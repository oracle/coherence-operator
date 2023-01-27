'use strict';

const core = require('@actions/core');

const errorHandler = error => {
  console.error(error);
  core.error(error.message);
  core.setFailed(error.message);
};

module.exports = errorHandler;
