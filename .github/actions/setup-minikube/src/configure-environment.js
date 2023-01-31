'use strict';

const core = require('@actions/core');
const logExecSync = require('./exec').logExecSync;
const download = require('./download');
const MAX_ATTEMPTS = 10;
const RETRY_WAIT_MILLIS = 5_000;

const waitForDocker = async (attempt = 0) => {
  try {
    logExecSync(`docker version -f '{{.Server.Version}} - {{.Client.Version}}'`);
    core.info('Docker daemon is ready');
    return true;
  } catch (e) {
    if (attempt++ < MAX_ATTEMPTS) {
      await new Promise(r => setTimeout(r, RETRY_WAIT_MILLIS));
      return await waitForDocker(attempt);
    }
  }
  core.warning('Docker daemon appears to be unready, hoping for the best');
  return false;
};

const configureEnvironment = async (inputs = {}) => {
  core.info('Updating Environment configuration to support Minikube');
  logExecSync('sudo apt update -y');
  logExecSync('sudo apt-get install -y conntrack');
  if (inputs.driver === 'docker') {
    core.info('Waiting for Docker to be ready');
    await waitForDocker();
  }
  await download.installCriCtl(inputs);
  await download.installCriDockerd(inputs);
};

module.exports = configureEnvironment;
