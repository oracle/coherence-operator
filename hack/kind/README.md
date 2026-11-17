# KinD image archive checks

`make test-kind-image-archive` runs the content regressions and supplemental
Make/shell checks with Python 3.9+ (standard library only) and GNU Make. It needs
no Docker daemon, registry access or cluster. Coherence Certification also runs
these tests before creating its cluster.

The fixtures generate real manifests, image configurations and tiny tar layers.
The shared-layer fixture uses identical uncompressed layers with different gzip
digests, then models an export that exits successfully while omitting required
content. Tests reject missing, corrupt and truncated content, incorrect sizes,
wrong platforms, and inconsistent tag metadata. They also exercise the actual
loader's cleanup, command overrides, inherited Docker configuration and failure
status preservation. No image archives are stored in the repository.

`kind-load-coherence` validates the exported tar before importing it. For Docker's
OCI layout it checks the selected manifest/index, configuration and every layer
by descriptor size and SHA-256, along with the image configuration's platform
and the corresponding Docker manifest paths. For the classic Docker layout it
checks the configuration digest and each uncompressed layer against `diff_ids`.
Archives are read without extraction. The checker can also run directly:

```sh
python3 hack/kind/image_archive.py image.tar --image registry.example/team/image:tag --platform linux/amd64
```

The OCI layout marker requires a valid layout version and `index.json`; a missing
index cannot cause a fallback to classic Docker validation. An OCI index without
its layout marker is also rejected.

Content and reference metadata are separate results. A tagged export must retain
the requested tag in `RepoTags` and, for OCI exports, the image-name annotation
on the selected top-level descriptor. Names on nested descriptors do not replace
that import name.
Digest references (`repository@digest` and `tag@digest`) report complete content
independently, but return **unsupported** (exit 3) before import. Docker can omit
their requested name, leaving anonymous KinD imports that fail CRI lookup. This
known naming limitation has no alias or retag workaround here. Exit 1 denotes
invalid content or tag metadata; exit 0 means these archive checks passed.

Archive checks do not prove runtime reference lookup, unpacking or pod startup.
Native qualification must separately check the original requested reference via
CRI and complete/unpacked content on every node. Simulated arm64/variant cases
do not qualify an arm64 runtime. Authenticated registry pulls and final-source
hosted certification are separate gates.

On export, validation or import failure the loader prints the stage, requested
image/platform, actual Docker client/server versions, image-store driver/backend,
selected image ID and repository references. When available it checks the failed
archive as well. Diagnostics use selected fields rather than full Docker info,
image configuration or environment dumps, preserve the original failure status,
and remove the temporary archive directory.
