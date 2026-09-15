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

import * as os from 'os';
import releaseJson from './data/releases.json';

// release is the type of official ORAS CLI release
interface releases {
    [version: string]: {
        [platform: string]: {
            [arch: string]: {
                checksum: string,
                url: string
            }
        }
    }
}

// Get release info of a certain verion of ORAS CLI
export function getReleaseInfo(version: string, url: string, checksum: string) {
    if (url && checksum) {
        // if customized ORAS CLI link and checksum are provided, version is ignored
        return {
            checksum: checksum,
            url: url
        }
    }

    // sanity checks
    if (url && !checksum) {
        throw new Error("user provided url of customized ORAS CLI release but without SHA256 checksum");
    }
    if (!url && checksum) {
        throw new Error("user provided SHA256 checksum but without url");
    }

    // get the official release
    const releases = releaseJson as releases;
    if (!(version in releases)) {
        console.log(`official ORAS CLI releases does not contain version ${version}`)
        throw new Error(`official ORAS CLI releases does not contain version ${version}`);
    }

    const platform = mapPlatform();
    const arch = mapArch();
    const download = releases[version][platform][arch];
    if (!download) {
        throw new Error(`official ORAS CLI releases does not contain version ${version}, platform ${platform}, arch ${arch} is not supported`);
    }
    return download;
}


// getPlatform maps os.platform() to ORAS supported platforms.
export function mapPlatform(): string {
    const platform: string = os.platform();
    switch (platform) {
        case 'linux':
            return 'linux';
        case 'darwin':
            return 'darwin';
        case 'win32':
            return 'windows';
        case 'freebsd':
            return 'freebsd';
        default:
            throw new Error(`unsupported platform: ${platform}`);
    }
}

// mapArch maps os.arch() to ORAS supported architectures.
export function mapArch(): string {
    const architecture: string = os.arch();
    switch (architecture) {
        case 'x64':
            return 'amd64';
        case 'arm64':
            return 'arm64';
        case 'arm64':
            return 'arm64';
        case 'ppc64':
            return 'ppc64le';
        case 'riscv64':
            return 'riscv64';
        case 's390x':
            return 's390x';
        case 'arm':
            return 'armv7';
        default:
            throw new Error(`unsupported architecture: ${architecture}`);
    }
}

export function getBinaryExtension(): string {
    const platform = mapPlatform();
    return platform === 'windows' ? '.exe' : '';
}
