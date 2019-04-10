#
# Copyright (c) 2019, Oracle and/or its affiliates. All rights reserved.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
#
#!/usr/bin/env sh

#!/bin/sh -e -x -u

echo "Starting init script"
echo ""
echo "Lib directory is: ${LIB_DIR}"
echo "External Lib directory is: ${EXTERNAL_LIB_DIR}"
echo "Config directory is: ${CONF_DIR}"
echo "External config directory is: ${EXTERNAL_CONF_DIR}"

if [ -d ${LIB_DIR} ]; then
  echo "Copying files from ${LIB_DIR} to ${EXTERNAL_LIB_DIR}"
  cp -R ${LIB_DIR}/* ${EXTERNAL_LIB_DIR}
else
  echo "Lib directory ${LIB_DIR} does not exist - no files to copy"
fi

if [ -d ${CONF_DIR} ]; then
  echo "Copying files from ${CONF_DIR} to ${EXTERNAL_CONF_DIR}"
  cp -R ${CONF_DIR}/* ${EXTERNAL_CONF_DIR}
else
  echo "Config directory ${CONF_DIR} does not exist - no files to copy"
fi

echo ""
echo "---------------------------------------------------"
echo "Contents of ${EXTERNAL_LIB_DIR}"
echo "---------------------------------------------------"
ls -Ral ${EXTERNAL_LIB_DIR}
echo "---------------------------------------------------"
echo ""
echo "---------------------------------------------------"
echo "Contents of ${EXTERNAL_CONF_DIR}"
echo "---------------------------------------------------"
ls -Ral ${EXTERNAL_CONF_DIR}
echo "---------------------------------------------------"
echo ""
echo "Finished init script"
