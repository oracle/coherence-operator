#!/bin/sh
#
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#

set -eu

VERSION=${1:?Usage: get-pack.sh VERSION DESTINATION}
DESTINATION=${2:?Usage: get-pack.sh VERSION DESTINATION}
SCRIPT_DIRECTORY=$(CDPATH='' cd -- "$(dirname "$0")" && pwd)

case "$(uname -s)" in
  Darwin) PLATFORM=macos ;;
  Linux) PLATFORM=linux ;;
  *) echo "Unsupported pack operating system: $(uname -s)" >&2; exit 1 ;;
esac

case "$(uname -m)" in
  x86_64|amd64) ARCH_SUFFIX= ;;
  arm64|aarch64) ARCH_SUFFIX=-arm64 ;;
  *) echo "Unsupported pack architecture: $(uname -m)" >&2; exit 1 ;;
esac

ARCHIVE=pack-v${VERSION}-${PLATFORM}${ARCH_SUFFIX}.tgz
CHECKSUM_FILE=$SCRIPT_DIRECTORY/pack-v${VERSION}.sha256
CHECKSUM=$(awk -v archive="$ARCHIVE" '$2 == archive { print $1 }' "$CHECKSUM_FILE")
if [ "${#CHECKSUM}" -ne 64 ]; then
  echo "No pinned checksum for $ARCHIVE in $CHECKSUM_FILE" >&2
  exit 1
fi

INSTALL_DIRECTORY=$(dirname "$DESTINATION")
mkdir -p "$INSTALL_DIRECTORY"
# Keep the verified binary on the destination filesystem for an atomic install.
DOWNLOAD_DIRECTORY=$(mktemp -d "$INSTALL_DIRECTORY/.pack-download.XXXXXX")
trap 'rm -rf "$DOWNLOAD_DIRECTORY"' EXIT HUP INT TERM

echo "Downloading pack $VERSION for $PLATFORM${ARCH_SUFFIX}"
curl -fsSL --retry 3 --connect-timeout 15 --max-time 120 \
  "https://github.com/buildpacks/pack/releases/download/v${VERSION}/${ARCHIVE}" \
  -o "$DOWNLOAD_DIRECTORY/$ARCHIVE"

if command -v sha256sum >/dev/null 2>&1; then
  printf '%s  %s\n' "$CHECKSUM" "$DOWNLOAD_DIRECTORY/$ARCHIVE" | sha256sum -c -
else
  printf '%s  %s\n' "$CHECKSUM" "$DOWNLOAD_DIRECTORY/$ARCHIVE" | shasum -a 256 -c -
fi

tar -xzf "$DOWNLOAD_DIRECTORY/$ARCHIVE" -C "$DOWNLOAD_DIRECTORY" pack
chmod 755 "$DOWNLOAD_DIRECTORY/pack"
# The archive timestamp predates the installer; let Make cache this installation.
touch "$DOWNLOAD_DIRECTORY/pack"
mv -f "$DOWNLOAD_DIRECTORY/pack" "$DESTINATION"
