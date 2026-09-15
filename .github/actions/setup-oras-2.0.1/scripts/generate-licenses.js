// Copyright The ORAS Authors.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Reads the esbuild metafile to find all bundled packages, then collects
// their license texts from node_modules and writes dist/licenses.txt.

'use strict';

const fs = require('fs');
const path = require('path');

const metaPath = process.argv[2] || 'dist/meta.json';
const outPath = process.argv[3] || 'dist/licenses.txt';

const meta = JSON.parse(fs.readFileSync(metaPath, 'utf8'));

// Collect unique package paths from the metafile inputs.
const pkgPaths = {};
for (const inputFile of Object.keys(meta.inputs)) {
  // Match both flat and nested node_modules paths.
  const m = inputFile.match(/^((?:.*\/)?node_modules\/)(@[^/]+\/[^/]+|[^/@][^/]*)\//);
  if (!m) continue;
  const pkgName = m[2];
  const pkgDir = m[1] + pkgName;
  if (!pkgPaths[pkgName]) {
    pkgPaths[pkgName] = pkgDir;
  }
}

const LICENSE_FILES = ['LICENSE', 'LICENSE.md', 'LICENSE.txt', 'LICENCE', 'license', 'license.md'];

const entries = [];
for (const [pkgName, pkgDir] of Object.entries(pkgPaths).sort()) {
  const pkgJson = JSON.parse(fs.readFileSync(path.join(pkgDir, 'package.json'), 'utf8'));
  const licenseType = pkgJson.license || 'UNKNOWN';

  let licenseText = '';
  for (const candidate of LICENSE_FILES) {
    const candidate_path = path.join(pkgDir, candidate);
    if (fs.existsSync(candidate_path)) {
      licenseText = fs.readFileSync(candidate_path, 'utf8').trimEnd();
      break;
    }
  }

  entries.push(`${pkgName}\n${licenseType}\n${licenseText}`);
}

fs.writeFileSync(outPath, entries.join('\n\n') + '\n');
console.log(`wrote ${outPath} with ${entries.length} package licenses`);
