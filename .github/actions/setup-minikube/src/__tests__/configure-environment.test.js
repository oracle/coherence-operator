describe('configure-docker module test suite', () => {
  let configureEnvironment;
  let logExecSync;
  beforeEach(() => {
    jest.resetModules();
    jest.mock('../exec');
    jest.mock('../download');
    configureEnvironment = require('../configure-environment');
    logExecSync = require('../exec').logExecSync;
  });
  test('configureEnvironment, should run all configuration commands', () => {
    // Given
    logExecSync.mockImplementation(() => {});
    // When
    configureEnvironment();
    // Then
    expect(logExecSync).toHaveBeenCalledTimes(2);
  });
  test('configureEnvironment with docker driver, should run all configuration commands', () => {
    // Given
    logExecSync.mockImplementation(() => {});
    // When
    configureEnvironment({driver: 'docker'});
    // Then
    expect(logExecSync).toHaveBeenCalledTimes(3);
  });
});
