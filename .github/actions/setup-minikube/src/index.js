#!/usr/bin/env node
'use strict';

const errorHandler = require('./error-handler');
const checkEnvironment = require('./check-environment');
const configureEnvironment = require('./configure-environment');
const loadInputs = require('./load-inputs');
const download = require('./download');
const install = require('./install');

const run = async () => {
  checkEnvironment();
  const inputs = loadInputs();
  await configureEnvironment(inputs);
  const downloadedFile = await download.downloadMinikube(inputs);
  await install(downloadedFile, inputs);
};

process.on('unhandledRejection', errorHandler);
run().catch(errorHandler);
