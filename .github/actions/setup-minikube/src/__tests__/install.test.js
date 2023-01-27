describe('install module test suite', () => {
  let core;
  let io;
  let path;
  let exec;
  let install;
  beforeEach(() => {
    jest.resetModules();
    jest.mock('@actions/core');
    jest.mock('@actions/io', () => ({
      mv: jest.fn(() => {})
    }));
    jest.mock('path');
    jest.mock('../exec');
    core = require('@actions/core');
    io = require('@actions/io');
    path = require('path');
    exec = require('../exec');
    install = require('../install');
  });
  test('install, should perform necessary steps', async () => {
    // Given
    const inputs = {minikubeVersion: 'v1.33.7'};
    exec.logExecSync.mockImplementation();
    exec.execSync.mockImplementation(() => '');
    // When
    await install('minikubeFileLocation', inputs);
    // Then
    expect(exec.logExecSync).toHaveBeenCalledTimes(5);
    expect(exec.execSync).toHaveBeenCalledTimes(1);
  });
});
