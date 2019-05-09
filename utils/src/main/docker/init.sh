#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#!/usr/bin/env sh

#!/bin/sh -e -x -u

echo "Starting init script"

if [ -z "${UTIL_DIR}" ]; then
  UTIL_DIR="/utils"
fi

echo "Creating target directories under ${UTIL_DIR}"

mkdir ${UTIL_DIR}/scripts ${UTIL_DIR}/lib ${UTIL_DIR}/conf

echo "Copying files to ${UTIL_DIR}"

cp files/*.sh ${UTIL_DIR}/scripts/
cp files/*.jar ${UTIL_DIR}/lib/

if [ ! -d /snapshot ]; then
    mkdir /snapshot
fi
chmod 0777 /snapshot

if [ ! -d /persistence ]; then
    mkdir /persistence
fi
chmod 0777 /persistence

if [ ! -d /persistence/active ]; then
    mkdir /persistence/active
fi
if [ ! -d /persistence/trash ]; then
    mkdir /persistence/trash
fi
chmod 0777 /persistence/active
chmod 0777 /persistence/trash

echo "Finished init script"
