#!/usr/bin/env python3
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
"""Check a single-platform Docker image-save tar without extracting it."""

import argparse
import hashlib
import json
from pathlib import Path, PurePosixPath
import re
import sys
import tarfile


class InvalidArchive(ValueError):
    pass


def require(condition, message):
    if not condition:
        raise InvalidArchive(message)


def blob_path(digest):
    require(isinstance(digest, str) and re.fullmatch(r"sha256:[0-9a-f]{64}", digest),
            "invalid or unsupported content digest: " + str(digest))
    return "blobs/" + digest.replace(":", "/", 1)


def canonical_tag(reference):
    """Match Docker's familiar and fully qualified tag spellings."""
    parts = reference.split("/")
    if len(parts) == 1 or not ("." in parts[0] or ":" in parts[0] or parts[0] == "localhost"):
        parts.insert(0, "docker.io")
    if parts[0] == "index.docker.io":
        parts[0] = "docker.io"
    if parts[0] == "docker.io" and len(parts) == 2:
        parts.insert(1, "library")
    if ":" not in parts[-1]:
        parts[-1] += ":latest"
    return "/".join(parts)


def platform_matches(actual, wanted):
    variant = actual.get("variant") or ("v8" if actual.get("architecture") == "arm64" else "")
    return (actual.get("os") == wanted[0] and actual.get("architecture") == wanted[1]
            and (len(wanted) == 2 or variant == wanted[2]))


class Archive:
    def __init__(self, tar):
        self.tar = tar
        self.members = {}
        self.errors = []
        self.checked = {}
        for member in tar.getmembers():
            name = member.name.removeprefix("./")
            path = PurePosixPath(name)
            require(not path.is_absolute() and ".." not in path.parts,
                    "unsafe archive member: " + name)
            if member.isdir():
                continue
            require(name not in self.members, "duplicate archive member: " + name)
            require(member.isfile(), "non-regular archive member: " + name)
            self.members[name] = member

    def stream(self, name):
        require(name in self.members, "missing required archive member: " + name)
        return self.tar.extractfile(self.members[name])

    def document(self, name):
        with self.stream(name) as stream:
            require(self.members[name].size <= 4 * 1024 * 1024, "oversized image metadata: " + name)
            return json.load(stream)

    def check(self, name, digest, size=None, role="content"):
        # Reuse the descriptor size/SHA-256 checks for manifests, configs and layers.
        if name not in self.members:
            self.errors.append("missing required " + role + ": " + name)
            return False
        if size is not None and self.members[name].size != size:
            self.errors.append("invalid size for " + role + ": " + name
                               + " (expected " + str(size) + ", got " + str(self.members[name].size) + ")")
            return False
        if name not in self.checked:
            value = hashlib.sha256()
            with self.stream(name) as stream:
                for block in iter(lambda: stream.read(1024 * 1024), b""):
                    value.update(block)
            self.checked[name] = "sha256:" + value.hexdigest()
        if self.checked[name] != digest:
            self.errors.append("invalid SHA-256 for " + role + ": " + name)
            return False
        return True

    def descriptor(self, desc, role):
        name = blob_path(desc["digest"])
        require(type(desc["size"]) is int and desc["size"] >= 0, "invalid descriptor size: " + name)
        return name, self.check(name, desc["digest"], desc["size"], role)

    def manifests(self, descriptors, wanted, ancestors=(), image_descriptor=None):
        found = []
        for desc in descriptors:
            if desc.get("platform") and not platform_matches(desc["platform"], wanted):
                continue  # Unselected sibling platforms need not have been exported.
            require(desc["digest"] not in ancestors, "cyclic image index")
            # Containerd creates image names only from the outer index's descriptors.
            imported = image_descriptor if image_descriptor is not None else desc
            name, valid = self.descriptor(desc, "manifest/index")
            if not valid:
                continue
            document = self.document(name)
            require(document.get("schemaVersion") == 2, "unsupported manifest schema: " + name)
            if "manifests" in document:
                found.extend(self.manifests(document["manifests"], wanted, (*ancestors, desc["digest"]),
                                            imported))
            else:
                require(isinstance(document.get("config"), dict) and isinstance(document.get("layers"), list),
                        "invalid image manifest: " + name)
                found.append((desc, document, imported))
        return found


def validate(path, reference, platform):
    result = {"image": reference, "platform": platform, "content": "invalid",
              "reference_metadata": "not_checked", "runtime_reference": "not_checked", "errors": []}
    try:
        wanted = platform.split("/")
        require(len(wanted) in (2, 3) and all(wanted), "expected platform os/architecture[/variant]")
        path = Path(path)
        result["archive_bytes"] = path.stat().st_size
        with path.open("rb") as stream:
            require(path.stat().st_size >= 1024 and path.stat().st_size % 512 == 0, "truncated image tar")
            stream.seek(-1024, 2)
            require(stream.read() == b"\0" * 1024, "missing image tar end markers")
        with tarfile.open(path, "r:") as tar:
            archive = Archive(tar)
            result["errors"] = archive.errors
            entries = archive.document("manifest.json")
            require(isinstance(entries, list) and len(entries) == 1, "expected one Docker manifest entry")
            entry = entries[0]
            require(isinstance(entry.get("Layers"), list), "invalid Docker layer list")
            tags = entry.get("RepoTags") or []
            require(isinstance(tags, list) and all(isinstance(tag, str) for tag in tags), "invalid RepoTags")
            image_descriptor = None
            if "oci-layout" in archive.members or "index.json" in archive.members:
                result["format"] = "oci"
                require(archive.document("oci-layout").get("imageLayoutVersion") == "1.0.0", "unsupported OCI layout")
                index = archive.document("index.json")
                require(index.get("schemaVersion") == 2, "invalid OCI index schema")
                selected = archive.manifests(index["manifests"], wanted)
                require(len(selected) == 1, "expected one selected image manifest for " + platform)
                desc, manifest, image_descriptor = selected[0]
                result["selected_manifest"] = desc["digest"]
                config_name, config_valid = archive.descriptor(manifest["config"], "configuration")
                layer_names = [archive.descriptor(layer, "layer")[0] for layer in manifest["layers"]]
                require(entry["Config"] == config_name and entry["Layers"] == layer_names,
                        "Docker manifest paths disagree with the selected OCI manifest")
                result["configuration"] = manifest["config"]["digest"]
            else:
                result["format"] = "docker"
                config_name = entry["Config"]
                config_digest = "sha256:" + PurePosixPath(config_name).stem
                blob_path(config_digest)
                config_valid = archive.check(config_name, config_digest, role="configuration")
                layer_names = entry["Layers"]
                result["configuration"] = config_digest
            result["layers"] = len(layer_names)
            if config_valid:
                config = archive.document(config_name)
                require(platform_matches(config, wanted), "image configuration does not match " + platform)
                require(config["rootfs"]["type"] == "layers", "invalid image rootfs type")
                diff_ids = config["rootfs"]["diff_ids"]
                require(len(diff_ids) == len(layer_names), "configuration diff_ids disagree with layer count")
                for diff_id in diff_ids:
                    blob_path(diff_id)
                if image_descriptor is None:
                    for name, digest in zip(layer_names, diff_ids):
                        archive.check(name, digest, role="uncompressed layer")
            if archive.errors:
                return result
            result["content"] = "complete"
            result["verified_objects"] = len(archive.checked)
            if "@" in reference:
                result["reference_metadata"] = "unsupported"
                result["errors"].append(
                    "Digest-pinned reference naming is not qualified: Docker can export complete content "
                    "without the requested name, which KinD imports anonymously. CRI lookup remains a known failed gate.")
            else:
                expected = canonical_tag(reference)
                require(expected in [canonical_tag(tag) for tag in tags], "requested tag missing from RepoTags")
                if image_descriptor is not None:
                    name = image_descriptor.get("annotations", {}).get("io.containerd.image.name", "")
                    require(name and canonical_tag(name) == expected,
                            "requested tag missing or different in OCI image-name annotation")
                result["reference_metadata"] = "present"
    except (InvalidArchive, OSError, tarfile.TarError, ValueError, KeyError, TypeError, AttributeError) as error:
        result["errors"].append(str(error))
    return result


def exit_status(result):
    if result["content"] != "complete":
        return 1
    if result["reference_metadata"] == "unsupported":
        return 3
    return 0 if result["reference_metadata"] == "present" and not result["errors"] else 1


def main():
    parser = argparse.ArgumentParser(description=__doc__)
    parser.add_argument("archive", type=Path)
    parser.add_argument("--image", required=True)
    parser.add_argument("--platform", required=True)
    args = parser.parse_args()
    result = validate(args.archive, args.image, args.platform)
    print(json.dumps(result, indent=2))
    return exit_status(result)


if __name__ == "__main__":
    sys.exit(main())
