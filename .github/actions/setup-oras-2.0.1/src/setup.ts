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

import * as core from '@actions/core';
import * as tc from '@actions/tool-cache';
import { getReleaseInfo } from './lib/release';
import { hash } from './lib/checksum';

// setup sets up the ORAS CLI
async function setup(): Promise<void> {
  try {
    // inputs from user
    const version: string = core.getInput('version');
    const url: string = core.getInput('url');
    const checksum = core.getInput('checksum').toLowerCase();

    // download ORAS CLI and validate checksum
    const info = getReleaseInfo(version, url, checksum);
    const download_url = info.url;
    console.log(`downloading ORAS CLI from ${download_url}`);
    const pathToTarball: string = await tc.downloadTool(download_url);
    console.log("downloading ORAS CLI completed");
    const actual_checksum = await hash(pathToTarball);
    if (actual_checksum !== info.checksum) {
      throw new Error(`checksum of downloaded ORAS CLI ${actual_checksum} does not match expected checksum ${info.checksum}`);
    }
    console.log("successfully verified downloaded release checksum");

    // extract the tarball/zipball onto host runner
    const extract = download_url.endsWith('.zip') ? tc.extractZip : tc.extractTar;
    const pathToCLI: string = await extract(pathToTarball);

    // add `ORAS` to PATH
    core.addPath(pathToCLI);
  } catch (e) {
    if (e instanceof Error) {
      core.setFailed(e);
    } else {
      core.setFailed('unknown error during ORAS setup');
    }
  }
}

export = setup;

if (require.main === module) {
  setup();
}