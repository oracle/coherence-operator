#!/usr/bin/env python3
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
"""Content regressions and supplemental checks of the actual Make loader."""

import copy
import gzip
import hashlib
import io
import json
import os
from pathlib import Path
import shlex
import shutil
import subprocess
import sys
import tarfile
import tempfile
import unittest

from image_archive import exit_status, validate


IMAGE = "example.invalid:5000/team/coherence:test"
PLATFORM = "linux/amd64"
REPO = Path(__file__).resolve().parents[2]


def encoded(document):
    return json.dumps(document, sort_keys=True, separators=(",", ":")).encode()


def digest(data):
    return "sha256:" + hashlib.sha256(data).hexdigest()


def descriptor(data, media):
    return {"digest": digest(data), "size": len(data), "mediaType": media}


def blob_name(data):
    return "blobs/sha256/" + hashlib.sha256(data).hexdigest()


def layer(payload):
    stream = io.BytesIO()
    with tarfile.open(fileobj=stream, mode="w") as tar:
        member = tarfile.TarInfo("file")
        member.size = len(payload)
        tar.addfile(member, io.BytesIO(payload))
    return stream.getvalue()


class Fixture:
    def __init__(self, image=IMAGE, architecture="amd64", variant=None):
        # The base and target share diffIDs but use different compressed blobs.
        # This is the shared-layer trigger; the failed save omits even the config.
        self.raw_layers = [layer(b"shared base content\n" * 1024), layer(b"shared runtime\n" * 1024),
                           layer(b"application\n")]
        self.base_layers = [gzip.compress(raw, compresslevel=1, mtime=0) for raw in self.raw_layers[:2]]
        self.layers = [gzip.compress(raw, compresslevel=9, mtime=0) for raw in self.raw_layers]
        self.config = {"os": "linux", "architecture": architecture,
                       "config": {"Env": ["SECRET_SENTINEL=must-not-appear"]},
                       "rootfs": {"type": "layers", "diff_ids": [digest(raw) for raw in self.raw_layers]}}
        self.platform = {"os": "linux", "architecture": architecture}
        if variant:
            self.platform["variant"] = variant
            self.config["variant"] = variant
        config_bytes = encoded(self.config)
        self.manifest = {"schemaVersion": 2, "mediaType": "application/vnd.oci.image.manifest.v1+json",
                         "config": descriptor(config_bytes, "application/vnd.oci.image.config.v1+json"),
                         "layers": [descriptor(data, "application/vnd.oci.image.layer.v1.tar+gzip")
                                    for data in self.layers]}
        self.members = {"oci-layout": encoded({"imageLayoutVersion": "1.0.0"}),
                        blob_name(config_bytes): config_bytes,
                        **{blob_name(data): data for data in self.layers}}
        self.entry = {"Config": blob_name(config_bytes), "Layers": [blob_name(data) for data in self.layers],
                      "RepoTags": [image]}
        self.image = image
        self.update_manifest()

    def update_manifest(self):
        data = encoded(self.manifest)
        self.members[blob_name(data)] = data
        desc = descriptor(data, self.manifest["mediaType"])
        desc.update(platform=self.platform, annotations={"io.containerd.image.name": self.image})
        self.index = {"schemaVersion": 2, "manifests": [desc]}
        self.metadata()

    def metadata(self):
        self.members["index.json"] = encoded(self.index)
        self.members["manifest.json"] = encoded([self.entry])

    def missing_content(self):
        for name in [self.entry["Config"], *self.entry["Layers"]]:
            self.members.pop(name)

    def anonymous(self):
        self.entry["RepoTags"] = None
        self.index["manifests"][0]["annotations"] = {}
        self.metadata()

    def nest(self, image):
        data = encoded(self.index)
        self.members[blob_name(data)] = data
        parent = descriptor(data, "application/vnd.oci.image.index.v1+json")
        parent["annotations"] = {"io.containerd.image.name": image} if image else {}
        self.index = {"schemaVersion": 2, "manifests": [parent]}
        self.metadata()

    def uncompressed(self):
        self.layers = self.raw_layers
        self.manifest["layers"] = [descriptor(data, "application/vnd.oci.image.layer.v1.tar")
                                   for data in self.layers]
        self.members.update({blob_name(data): data for data in self.layers})
        self.entry["Layers"] = [blob_name(data) for data in self.layers]
        self.update_manifest()

    def legacy(self):
        config = encoded(self.config)
        self.members = {hashlib.sha256(config).hexdigest() + ".json": config}
        names = []
        for number, raw in enumerate(self.raw_layers):
            name = str(number) + "/layer.tar"
            self.members[name] = raw
            names.append(name)
        self.entry.update(Config=next(iter(self.members)), Layers=names)
        self.members["manifest.json"] = encoded([self.entry])

    def write(self, path, duplicate=None):
        with tarfile.open(path, "w") as tar:
            for name, data in self.members.items():
                member = tarfile.TarInfo(name)
                member.size = len(data)
                tar.addfile(member, io.BytesIO(data))
            if duplicate:
                member = tarfile.TarInfo(duplicate)
                tar.addfile(member, io.BytesIO())
        return path


class ContentTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory(prefix="kind-archive-test-")
        self.addCleanup(self.temp.cleanup)
        self.path = Path(self.temp.name) / "image.tar"

    def check(self, fixture, image=IMAGE, platform=PLATFORM):
        return validate(fixture.write(self.path), image, platform)

    def test_valid_selected_content(self):
        result = self.check(Fixture())
        self.assertEqual(exit_status(result), 0, result)
        self.assertEqual(result["verified_objects"], 5)
        self.assertEqual(result["runtime_reference"], "not_checked")

    def test_shared_diffids_different_compressed_digests_successful_but_incomplete_save(self):
        fixture = Fixture()
        for base, target, raw in zip(fixture.base_layers, fixture.layers, fixture.raw_layers):
            self.assertEqual(digest(gzip.decompress(base)), digest(raw))
            self.assertEqual(digest(gzip.decompress(target)), digest(raw))
            self.assertNotEqual(digest(base), digest(target))
        fixture.missing_content()
        result = self.check(fixture)
        self.assertEqual(exit_status(result), 1)
        self.assertEqual(len(result["errors"]), 4, result)
        self.assertIn("missing required configuration", result["errors"][0])

    def test_missing_and_corrupt_manifest_configuration_and_each_layer(self):
        original = Fixture()
        required = [blob_name(encoded(original.manifest)), original.entry["Config"], *original.entry["Layers"]]
        for name in required:
            for mutation in ("missing", "corrupt"):
                with self.subTest(name=name, mutation=mutation):
                    fixture = copy.deepcopy(original)
                    if mutation == "missing":
                        fixture.members.pop(name)
                    else:
                        data = fixture.members[name]
                        fixture.members[name] = bytes([data[0] ^ 1]) + data[1:]
                    result = self.check(fixture)
                    self.assertEqual(exit_status(result), 1, result)
                    self.assertTrue(any(name in error for error in result["errors"]), result)

    def test_descriptor_size(self):
        for role in ("config", "layer"):
            with self.subTest(role=role):
                fixture = Fixture()
                desc = fixture.manifest["config"] if role == "config" else fixture.manifest["layers"][0]
                desc["size"] += 1
                fixture.update_manifest()
                result = self.check(fixture)
                self.assertEqual(exit_status(result), 1, result)
                self.assertTrue(any("invalid size" in error for error in result["errors"]), result)

    def test_truncated_tar(self):
        data = Fixture().write(self.path).read_bytes()
        for size in (0, 1025, len(data) // 2, len(data) - 1000):
            with self.subTest(size=size):
                self.path.write_bytes(data[:size])
                self.assertEqual(exit_status(validate(self.path, IMAGE, PLATFORM)), 1)

    def test_duplicate_and_unsafe_members(self):
        for name in ("manifest.json", "../outside", "/absolute"):
            with self.subTest(name=name):
                Fixture().write(self.path, duplicate=name)
                self.assertEqual(exit_status(validate(self.path, IMAGE, PLATFORM)), 1)

    def test_malformed_metadata_and_inconsistent_docker_paths(self):
        for member in ("index.json", "manifest.json", "oci-layout"):
            with self.subTest(member=member):
                fixture = Fixture()
                fixture.members[member] = b"{"
                self.assertEqual(exit_status(self.check(fixture)), 1)
        fixture = Fixture()
        fixture.entry["Layers"].reverse()
        fixture.metadata()
        self.assertEqual(exit_status(self.check(fixture)), 1)

    def test_oci_layout_requires_index_and_supported_version(self):
        for compressed in (True, False):
            original = Fixture()
            if not compressed:
                original.uncompressed()
            result = self.check(original)
            self.assertEqual(exit_status(result), 0, result)
            self.assertEqual(result["format"], "oci")
            for mutation in ("missing-index", "unsupported-layout", "unsupported-layout-missing-index",
                             "malformed-layout-missing-index", "missing-layout"):
                with self.subTest(compressed=compressed, mutation=mutation):
                    fixture = copy.deepcopy(original)
                    if "missing-index" in mutation:
                        fixture.members.pop("index.json")
                    if mutation.startswith("unsupported-layout"):
                        fixture.members["oci-layout"] = encoded({"imageLayoutVersion": "9.9.9"})
                    if mutation.startswith("malformed-layout"):
                        fixture.members["oci-layout"] = b"{"
                    if mutation == "missing-layout":
                        fixture.members.pop("oci-layout")
                    result = self.check(fixture)
                    self.assertEqual(exit_status(result), 1, result)
                    self.assertEqual(result["content"], "invalid", result)
                    self.assertTrue(result["errors"], result)

    def test_platform_and_variant(self):
        self.assertEqual(exit_status(self.check(Fixture(), platform="linux/arm64")), 1)
        for variant in (None, "v8"):
            result = self.check(Fixture(architecture="arm64", variant=variant), platform="linux/arm64/v8")
            self.assertEqual(exit_status(result), 0, result)
        self.assertEqual(exit_status(self.check(Fixture(architecture="arm64", variant="v9"),
                                                platform="linux/arm64/v8")), 1)
        fixture = Fixture(architecture="arm64")
        fixture.index["manifests"][0]["platform"] = {"os": "linux", "architecture": "amd64"}
        fixture.metadata()
        self.assertEqual(exit_status(self.check(fixture)), 1)

    def test_unselected_missing_sibling_is_not_required(self):
        fixture = Fixture()
        fixture.index["manifests"].append({"digest": "sha256:" + "a" * 64, "size": 100,
                                            "platform": {"os": "linux", "architecture": "arm64"}})
        fixture.metadata()
        self.assertEqual(exit_status(self.check(fixture)), 0)

    def test_nested_index_checks_parent_digest_and_preserves_name(self):
        fixture = Fixture()
        del fixture.index["manifests"][0]["annotations"]
        data = encoded(fixture.index)
        fixture.nest(IMAGE)
        self.assertEqual(exit_status(self.check(fixture)), 0)
        fixture.members[blob_name(data)] = b"corrupt parent index"
        self.assertEqual(exit_status(self.check(fixture)), 1)

    def test_nested_index_name_uses_top_level_descriptor(self):
        other = "example.invalid/different:tag"
        for outer in (IMAGE, other, None):
            for inner in (IMAGE, other, None):
                with self.subTest(outer=outer, inner=inner):
                    fixture = Fixture()
                    fixture.index["manifests"][0]["annotations"] = (
                        {"io.containerd.image.name": inner} if inner else {})
                    fixture.nest(outer)
                    result = self.check(fixture)
                    self.assertEqual(result["content"], "complete", result)
                    self.assertEqual(exit_status(result), 0 if outer == IMAGE else 1, result)
                    metadata = "present" if outer == IMAGE else "not_checked"
                    self.assertEqual(result["reference_metadata"], metadata)
        fixture = Fixture()
        fixture.nest(other)
        fixture.nest(IMAGE)
        self.assertEqual(exit_status(self.check(fixture)), 0)

    def test_invalid_descriptor_and_configuration_layer_count(self):
        fixture = Fixture()
        fixture.manifest["layers"][0]["digest"] = "sha256:invalid"
        fixture.update_manifest()
        self.assertEqual(exit_status(self.check(fixture)), 1)
        fixture = Fixture()
        fixture.config["rootfs"]["diff_ids"].pop()
        data = encoded(fixture.config)
        fixture.members[blob_name(data)] = data
        fixture.entry["Config"] = blob_name(data)
        fixture.manifest["config"] = descriptor(data, "application/vnd.oci.image.config.v1+json")
        fixture.update_manifest()
        self.assertEqual(exit_status(self.check(fixture)), 1)

    def test_tag_metadata_and_familiar_name(self):
        for location in ("docker", "oci"):
            with self.subTest(location=location):
                fixture = Fixture()
                if location == "docker":
                    fixture.entry["RepoTags"] = ["example.invalid/different:tag"]
                else:
                    fixture.index["manifests"][0]["annotations"] = {}
                fixture.metadata()
                result = self.check(fixture)
                self.assertEqual(result["content"], "complete")
                self.assertEqual(exit_status(result), 1)
        result = self.check(Fixture(image="docker.io/library/ubuntu:latest"), image="ubuntu")
        self.assertEqual(exit_status(result), 0, result)

    def test_digest_reference_content_complete_but_naming_unsupported(self):
        for reference in ("example.invalid/team/image@sha256:" + "a" * 64,
                          "example.invalid/team/image@sha256:" + "b" * 64,
                          IMAGE + "@sha256:" + "a" * 64):
            fixture = Fixture()
            fixture.anonymous()
            result = self.check(fixture, image=reference)
            self.assertEqual(result["content"], "complete", result)
            self.assertEqual(result["reference_metadata"], "unsupported", result)
            self.assertEqual(exit_status(result), 3)
            self.assertEqual(result["runtime_reference"], "not_checked")

    def test_legacy_docker_config_and_uncompressed_layers(self):
        original = Fixture()
        original.legacy()
        self.assertEqual(exit_status(self.check(original)), 0)
        for name in (original.entry["Config"], *original.entry["Layers"]):
            with self.subTest(name=name):
                fixture = copy.deepcopy(original)
                fixture.members[name] = b"corrupt content"
                self.assertEqual(exit_status(self.check(fixture)), 1)
        fixture = copy.deepcopy(original)
        fixture.entry["Layers"][1] = fixture.entry["Layers"][0]
        fixture.members["manifest.json"] = encoded([fixture.entry])
        self.assertEqual(exit_status(self.check(fixture)), 1, "a reused path must satisfy each diffID")


MOCK = r'''#!/usr/bin/env python3
import json, os, shutil, sys, tempfile
from pathlib import Path
command = Path(sys.argv[0]).name
args = sys.argv[1:]
with open(os.environ['TEST_CALLS'], 'a') as stream:
    stream.write(json.dumps({'command': command, 'args': args,
        'env': {key: os.environ.get(key) for key in ('DOCKER_HOST', 'DOCKER_CONFIG',
        'DOCKER_CONTEXT', 'DOCKER_TLS_VERIFY', 'DOCKER_CERT_PATH', 'OVERRIDE')}}) + '\n')
failure = os.environ.get('TEST_FAILURE', '')
if command == 'mktemp':
    if failure == 'mktemp': sys.exit(42)
    print(tempfile.mkdtemp(prefix='coherence-kind.', dir=os.environ['TEST_ARCHIVES']))
elif command == 'mock-docker':
    assert args[:2] == ['--context', 'test-context'], args
    args = args[2:]
    if args[0] == 'pull': sys.exit(41 if failure == 'pull' else 0)
    if args[:2] == ['image', 'save']:
        shutil.copyfile(os.environ['TEST_FIXTURE'], args[args.index('--output') + 1])
        sys.exit(43 if failure == 'save' else 0)
    if failure == 'diagnostics': sys.exit(45)
    if args[0] == 'version': print('Client=29.9.3 Server=29.9.3')
    elif args[0] == 'info': print('Driver=overlayfs DriverStatus=io.containerd.snapshotter.v1')
    elif args[:2] == ['image', 'inspect']: print('ID=sha256:selected OS=linux Architecture=amd64')
    else: raise AssertionError(args)
elif command == 'mock-kind':
    assert args[:4] == ['load', 'image-archive', '--name', 'test-cluster'], args
    assert Path(args[4]).is_file()
    sys.exit(44 if failure == 'import' else 0)
else: raise AssertionError(command)
'''


class MakeTests(unittest.TestCase):
    def test_loader_content_gate_arguments_cleanup_and_failure_diagnostics(self):
        make = os.environ.get("MAKE", "make")
        self.assertIsNotNone(shutil.which(make), "GNU Make is required")
        with tempfile.TemporaryDirectory(prefix="kind-loader-test-") as directory:
            root = Path(directory)
            bin_dir = root / "bin"
            bin_dir.mkdir()
            for command in ("mock-docker", "mock-kind", "mktemp"):
                path = bin_dir / command
                path.write_text(MOCK)
                path.chmod(0o755)
            ttl = root / "ttl-uuid"
            ttl.write_text("unused\n")
            for case in ("tag", "nested-tag", "nested-wrong-root", "nested-anonymous-root",
                         "missing-oci-index", "unsupported-oci-layout",
                         "incomplete", "corrupt", "truncated", "wrong-tag", "digest", "tag-digest",
                         "pull", "mktemp", "save", "import", "diagnostics", "arm64", "aarch64", "variant"):
                with self.subTest(case=case):
                    case_dir = root / case
                    case_dir.mkdir()
                    archives = case_dir / "temporary archives with spaces"
                    archives.mkdir()
                    fixture = Fixture()
                    reference = IMAGE
                    machine = "x86_64"
                    platform = PLATFORM
                    if case == "nested-tag":
                        fixture.index["manifests"][0]["annotations"] = (
                            {"io.containerd.image.name": "example.invalid/inner:tag"})
                        fixture.nest(IMAGE)
                    if case in ("nested-wrong-root", "nested-anonymous-root"):
                        fixture.nest("example.invalid/wrong:tag" if case == "nested-wrong-root" else None)
                    if case in ("missing-oci-index", "unsupported-oci-layout"):
                        fixture.uncompressed()
                        fixture.members.pop("index.json")
                        if case == "unsupported-oci-layout":
                            fixture.members["oci-layout"] = encoded({"imageLayoutVersion": "9.9.9"})
                    if case == "incomplete":
                        fixture.missing_content()
                    if case in ("corrupt", "diagnostics"):
                        fixture.members[fixture.entry["Layers"][0]] = b"corrupt"
                    if case == "wrong-tag":
                        fixture.entry["RepoTags"] = []
                        fixture.metadata()
                    if case in ("digest", "tag-digest"):
                        fixture.anonymous()
                        reference = (IMAGE if case == "tag-digest" else IMAGE.split(":test")[0]) + "@sha256:" + "a" * 64
                    if case in ("arm64", "aarch64", "variant"):
                        fixture = Fixture(architecture="arm64")
                        machine = "aarch64" if case == "aarch64" else "arm64"
                        platform = "linux/arm64/v8" if case == "variant" else "linux/arm64"
                    source = fixture.write(case_dir / "fixture.tar")
                    if case == "truncated":
                        source.write_bytes(source.read_bytes()[:1025])
                    calls_file = case_dir / "calls.jsonl"
                    env = os.environ.copy()
                    env.update(PATH=str(bin_dir) + os.pathsep + env["PATH"], FORCE_HOST_GO="true", GOTOOLCHAIN="local",
                               TEST_CALLS=str(calls_file), TEST_FIXTURE=str(source), TEST_ARCHIVES=str(archives),
                               TEST_FAILURE=case, TMPDIR=str(archives), DOCKER_HOST="unix:///test.sock",
                               DOCKER_CONFIG=str(root / "private Docker config"), DOCKER_CONTEXT="inherited-context",
                               DOCKER_TLS_VERIFY="1", DOCKER_CERT_PATH=str(root / "private certs"))
                    cmd = [make, "--no-print-directory", "-o", "kind-install", "kind-load-coherence",
                           "ISTIO_VERSION=test-unused", "TTL_UUID_FILE=" + str(ttl), "PATH=" + env["PATH"],
                           "GOTOOLCHAIN=local", "COHERENCE_IMAGE=" + reference, "UNAME_M=" + machine,
                           "DOCKER_CMD=env OVERRIDE=kept " + shlex.quote(str(bin_dir / "mock-docker")) + " --context test-context",
                           "KIND=" + shlex.quote(str(bin_dir / "mock-kind")), "KIND_CLUSTER=test-cluster"]
                    if case == "variant":
                        cmd.append("KIND_PLATFORM=" + platform)
                    process = subprocess.run(cmd, cwd=REPO, env=env, capture_output=True, text=True)
                    success = case in ("tag", "nested-tag", "arm64", "aarch64", "variant")
                    self.assertEqual(process.returncode, 0 if success else 2, process.stdout + process.stderr)
                    calls = [json.loads(line) for line in calls_file.read_text().splitlines()]
                    docker = [call["args"][2:] for call in calls if call["command"] == "mock-docker"]
                    self.assertEqual(docker[0], ["pull", "--platform", platform, reference])
                    saves = [args for args in docker if args[:2] == ["image", "save"]]
                    if case not in ("pull", "mktemp"):
                        self.assertEqual(len(saves), 1)
                        self.assertEqual(saves[0][:4], ["image", "save", "--platform", platform])
                        self.assertEqual(saves[0][-1], reference)
                    imports = [call for call in calls if call["command"] == "mock-kind"]
                    self.assertEqual(len(imports), int(success or case == "import"), process.stdout + process.stderr)
                    self.assertEqual(list(archives.iterdir()), [], "archive cleanup failed")
                    for call in calls:
                        for key in ("DOCKER_HOST", "DOCKER_CONFIG", "DOCKER_CONTEXT", "DOCKER_TLS_VERIFY", "DOCKER_CERT_PATH"):
                            self.assertEqual(call["env"][key], env[key])
                        if call["command"] == "mock-docker":
                            self.assertEqual(call["env"]["OVERRIDE"], "kept")
                    expected_error = {"pull": 41, "mktemp": 42, "save": 43, "import": 44,
                                      "digest": 3, "tag-digest": 3}.get(case, 1)
                    if not success:
                        self.assertIn("Error " + str(expected_error), process.stderr)
                    if case in ("save", "import", "incomplete", "corrupt", "truncated", "wrong-tag", "digest", "tag-digest",
                                "nested-wrong-root", "nested-anonymous-root", "missing-oci-index", "unsupported-oci-layout"):
                        self.assertIn("Client=29.9.3 Server=29.9.3", process.stderr)
                        self.assertIn("Driver=overlayfs", process.stderr)
                        self.assertIn("ID=sha256:selected", process.stderr)
                    self.assertNotIn("SECRET_SENTINEL", process.stdout + process.stderr)


if __name__ == "__main__":
    unittest.main(verbosity=2)
