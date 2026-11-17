# Operator CI shared build

`Prepare Build Artifacts` runs `make all` once. Quality checks, Java/Go unit tests,
local E2E, remote E2E, and Helm+FIPS then run on separate runners.
Each test runner starts its own KinD cluster. Unit tests also install CRDs.
Helm and FIPS run in sequence on one runner, with deployment/CRD cleanup.

The final `build` check requires preparation and every test job to succeed.
It uploads the original success artifacts and runs publication on `main` only.
Failure/cancellation logs use `test-output-<job>` artifact names.

| Internal artifact | Contents | Consumers |
| --- | --- | --- |
| `operator-ci-build` | Binaries, generated source/config/assets/docs, manifests and chart | All consumers and final `build` |
| `operator-ci-java` | Maven module `target/` outputs; omit Jib caches and duplicate exploded Spring archives | Quality/unit tests and main-only publication |
| `operator-ci-images` | Operator plus 11 test image tags, including both Spring buildpacks | Three integration runners |

`artifacts.sh` uses tar to preserve modes, symlinks, hidden files and timestamps.
Images are saved together to share layers, then imported directly into KinD.
The loader checks every required tag on every node. Each integration cluster also
runs `make kind-load-coherence`. Docker 29.9.3 matches `coherence-matrix.yaml`
to avoid the incomplete image archives seen with an earlier Docker version.

`make-prebuilt.sh` checks required outputs and prevents completed Make targets
from running again. It applies the same rules to recursive Make calls.
Restored timestamps alone cannot prevent rebuilds: source files can be newer,
and some prerequisites always run. Tests, deployment setup and certificate
generation still run. Use outputs from the same checkout.
Local/remote E2E install JDK 25 after cleanup to provide `keytool`.

Local checks:

```sh
shellcheck hack/ci/artifacts.sh hack/ci/make-prebuilt.sh
bash -n hack/ci/artifacts.sh hack/ci/make-prebuilt.sh
actionlint -ignore 'invalid activity type "committed"' .github/workflows/build.yaml
git diff --check
```

The unsupported `pull_request.types: committed` and missing `push-snapshot-docs`
target are inherited baseline issues. Publication retains `|| true`.

Validate the split on GitHub: run every suite, check image imports and aggregate
failure/cancellation gating, and measure archive sizes, transfers and total time.
Confirm required-check settings and publication separately. No speedup has been measured.
