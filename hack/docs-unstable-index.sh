#!/usr/bin/env bash

# This shell script creates the index.html file
# in the docs-unstable directory.

DIRS=$(ls -ld docs-unstable/* | grep '^d' | awk '{print $9}' | sort -r)

echo "<html><body><h2>Unstable Release Documentation</h2><ul>" > docs-unstable/index.html;
for i in ${DIRS}; do
  IFS='/' read -ra NAME <<< "${i}"
  echo "<li><a href=https://oracle.github.io/coherence-operator/${i}index.html>${NAME[1]}</a></li>" >> docs-unstable/index.html
done
echo "</ul></body></html>" >> docs-unstable/index.html

cat docs-unstable/index.html