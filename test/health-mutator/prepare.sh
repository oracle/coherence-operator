#!/usr/bin/env bash
# Copyright (c) 2026, Oracle and/or its affiliates.
# Licensed under the Universal Permissive License v 1.0 as shown at
# http://oss.oracle.com/licenses/upl.
set -euo pipefail
umask 077
output=${1:?output directory required}
dns=${2:?certificate DNS name required}
mkdir "$output"
openssl rand -hex 24 > "$output/store-password"
# Credential files intentionally have no trailing newline.
printf 'health-operator' > "$output/username"
printf '%s' "$(openssl rand -hex 32)" > "$output/password"
keytool -genkeypair -alias health -keyalg RSA -keysize 2048 -validity 30 \
  -dname "CN=$dns" -ext "SAN=dns:$dns" -storetype JKS \
  -keystore "$output/server.jks" -storepass:file "$output/store-password" \
  -keypass:file "$output/store-password" >/dev/null 2>&1
keytool -exportcert -rfc -alias health -keystore "$output/server.jks" \
  -storepass:file "$output/store-password" -file "$output/ca.pem" >/dev/null 2>&1
