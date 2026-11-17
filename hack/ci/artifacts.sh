#!/usr/bin/env bash
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
set -euo pipefail
cd "$(dirname "$0")/../.."
command=${1:?Expected pack or load-images}
artifact_dir=${2:?Expected absolute artifact directory}

case "$command" in
  pack)
    mkdir -p "$artifact_dir"
    # Tar preserves modes, symlinks, hidden generated files and mtimes. Avoid
    # Maven image caches and exploded duplicate Spring jars in the Java bundle.
    tar -czf "$artifact_dir/build.tar.gz" --exclude=bin/.gimme \
      bin build/_output config pkg/data/assets \
      api/v1/zz_generated.deepcopy.go docs/about/04_coherence_spec.adoc
    find java -mindepth 2 -maxdepth 2 -type d -name target -print > "$artifact_dir/java-files.txt"
    tar -czf "$artifact_dir/java.tar.gz" --exclude='*/target/jib-cache' --exclude='java/operator-test-spring*/target/spring' \
      -T "$artifact_dir/java-files.txt"
    rm "$artifact_dir/java-files.txt"
    image_list=$(make --no-print-directory -s ci-image-list)
    images=()
    while IFS= read -r image; do images+=("$image"); done <<< "$image_list"
    # One save deduplicates shared layers. Stream compression avoids keeping a
    # second uncompressed copy on the build runner. No registry publication.
    docker image save --platform linux/amd64 "${images[@]}" | gzip -1 > "$artifact_dir/images.tar.gz"
    du -h "$artifact_dir/"*.tar.gz
    ;;
  load-images)
    cluster=${KIND_CLUSTER:?Expected the same KIND_CLUSTER used by make kind}
    archive="$artifact_dir/images.tar"
    trap 'rm -f "$archive"' EXIT
    gzip -dc "$artifact_dir/images.tar.gz" > "$archive"
    # Import straight into each node, avoiding another copy in the host daemon.
    build/tools/bin/kind load image-archive --name "$cluster" "$archive"
    image_list=$(make --no-print-directory -s ci-image-list)
    nodes=$(build/tools/bin/kind get nodes --name "$cluster")
    test -n "$nodes"
    for node in $nodes; do
      while IFS= read -r image; do
        docker exec "$node" crictl inspecti "$image" > /dev/null
      done <<< "$image_list"
    done
    ;;
  *) echo "Unknown artifact command: $command" >&2; exit 1 ;;
esac
