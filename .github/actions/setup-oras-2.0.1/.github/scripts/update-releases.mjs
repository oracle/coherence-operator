#!/usr/bin/env node
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

// update-releases.mjs — add any missing ORAS releases to src/lib/data/releases.json.
//
// Usage:
//   node .github/scripts/update-releases.mjs              # detect new versions via GitHub API
//   node .github/scripts/update-releases.mjs 1.3.2 1.3.3  # add specific versions
//
// Outputs the list of added versions to stdout, one per line, and writes
// ADDED_VERSIONS=<comma-list> to $GITHUB_OUTPUT when run in Actions.

import { readFile, writeFile, appendFile } from 'node:fs/promises';
import { fileURLToPath } from 'node:url';
import { dirname, resolve } from 'node:path';

const ORAS_REPO = 'oras-project/oras';
const RELEASES_JSON = resolve(
    dirname(fileURLToPath(import.meta.url)),
    '..', '..', 'src', 'lib', 'data', 'releases.json'
);

async function ghFetch(url) {
    const headers = { 'Accept': 'application/vnd.github+json', 'User-Agent': 'setup-oras-update-releases' };
    if (process.env.GITHUB_TOKEN) headers['Authorization'] = `Bearer ${process.env.GITHUB_TOKEN}`;
    const res = await fetch(url, { headers });
    if (!res.ok) throw new Error(`GET ${url} -> ${res.status} ${res.statusText}`);
    return res.json();
}

async function fetchText(url) {
    const res = await fetch(url, { redirect: 'follow' });
    if (!res.ok) throw new Error(`GET ${url} -> ${res.status} ${res.statusText}`);
    return res.text();
}

function parseChecksums(text, version) {
    // Each line: "<sha256>  oras_<version>_<platform>_<arch>.tar.gz" (or .zip)
    const entry = {};
    const filenameRe = new RegExp(`^oras_${version.replace(/\./g, '\\.')}_([^_]+)_([^.]+)\\.(tar\\.gz|zip)$`);
    for (const rawLine of text.split('\n')) {
        const line = rawLine.trim();
        if (!line) continue;
        const parts = line.split(/\s+/);
        if (parts.length < 2) continue;
        const [checksum, filename] = parts;
        const m = filename.match(filenameRe);
        if (!m) continue; // skip e.g. checksum lines for source tarballs or unrelated assets
        const [, platform, arch, ext] = m;
        entry[platform] ??= {};
        entry[platform][arch] = {
            checksum,
            url: `https://github.com/${ORAS_REPO}/releases/download/v${version}/${filename}`,
        };
    }
    if (Object.keys(entry).length === 0) {
        throw new Error(`No release artifacts parsed from checksums file for v${version}`);
    }
    return entry;
}

async function fetchVersionEntry(version) {
    const url = `https://github.com/${ORAS_REPO}/releases/download/v${version}/oras_${version}_checksums.txt`;
    const text = await fetchText(url);
    return parseChecksums(text, version);
}

async function listUpstreamVersions() {
    // List published, non-prerelease, non-draft releases tagged like vX.Y.Z.
    const releases = await ghFetch(`https://api.github.com/repos/${ORAS_REPO}/releases?per_page=100`);
    return releases
        .filter(r => !r.draft && !r.prerelease)
        .map(r => r.tag_name)
        .filter(t => /^v\d+\.\d+\.\d+$/.test(t))
        .map(t => t.replace(/^v/, ''));
}

function compareSemver(a, b) {
    const pa = a.split('.').map(Number);
    const pb = b.split('.').map(Number);
    for (let i = 0; i < 3; i++) {
        if (pa[i] !== pb[i]) return pa[i] - pb[i];
    }
    return 0;
}

async function main() {
    const cliVersions = process.argv.slice(2);
    const raw = await readFile(RELEASES_JSON, 'utf8');
    const releases = JSON.parse(raw);

    let candidates;
    let missing;
    if (cliVersions.length > 0) {
        // Backfill mode: caller explicitly asked for these versions.
        candidates = cliVersions.map(v => v.replace(/^v/, ''));
        missing = candidates.filter(v => !(v in releases));
    } else {
        // Auto mode: only pull versions newer than the current max, so we don't
        // backfill 0.x releases that the project intentionally never tracked.
        const currentMax = Object.keys(releases).sort(compareSemver).pop();
        candidates = await listUpstreamVersions();
        missing = candidates
            .filter(v => !(v in releases))
            .filter(v => compareSemver(v, currentMax) > 0);
    }
    if (missing.length === 0) {
        console.log('No new versions to add.');
        if (process.env.GITHUB_OUTPUT) {
            await appendFile(process.env.GITHUB_OUTPUT, 'ADDED_VERSIONS=\n');
        }
        return;
    }

    const added = [];
    for (const version of missing) {
        console.log(`Adding v${version}...`);
        releases[version] = await fetchVersionEntry(version);
        added.push(version);
    }

    // Re-sort by semver so newer versions append to the end deterministically.
    const sorted = Object.fromEntries(
        Object.keys(releases).sort(compareSemver).map(k => [k, releases[k]])
    );

    // Match existing formatting: 4-space indent, trailing newline.
    await writeFile(RELEASES_JSON, JSON.stringify(sorted, null, 4) + '\n');

    console.log(`Added ${added.length} version(s): ${added.join(', ')}`);
    if (process.env.GITHUB_OUTPUT) {
        await appendFile(process.env.GITHUB_OUTPUT, `ADDED_VERSIONS=${added.join(',')}\n`);
    }
}

main().catch(err => {
    console.error(err.stack || err.message);
    process.exit(1);
});
