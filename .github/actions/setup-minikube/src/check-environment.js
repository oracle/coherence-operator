'use strict';

const fs = require('fs');

const isLinux = () => process.platform.toLowerCase().indexOf('linux') === 0;
const isUbuntu = version => () => {
  const osRelease = '/etc/os-release';
  const osInfo = fs.existsSync(osRelease) && fs.readFileSync(osRelease);
  return (
    osInfo &&
    osInfo.indexOf('NAME="Ubuntu"') >= 0 &&
    osInfo.indexOf(`VERSION="${version}`) >= 0
  );
};
['18', '20', '22'].some(v => isUbuntu(v)())
const isValidLinux = () => isLinux() && ['18', '20', '22'].some(v => isUbuntu(v)());
const checkOperatingSystem = () => {
  if (!isValidLinux()) {
    throw Error('Unsupported OS, action only works in Ubuntu 18, 20, or 22');
  }
};

const checkEnvironment = () => {
  checkOperatingSystem();
};

module.exports = checkEnvironment;
