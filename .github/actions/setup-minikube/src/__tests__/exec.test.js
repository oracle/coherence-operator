describe('exec module test suite', () => {
  let exec;
  let child_process;
  beforeEach(() => {
    jest.resetModules();
    jest.mock('child_process');
    exec = require('../exec');
    child_process = require('child_process');
  });
  test('execSync, should spawn the provided command', () => {
    // Given
    child_process.execSync.mockImplementationOnce(() => {});
    // When
    exec.execSync('1337');
    // Then
    expect(child_process.execSync).toHaveBeenCalledTimes(1);
    expect(child_process.execSync).toHaveBeenCalledWith('1337');
  });
  test('logExecSync, should spawn the provided command with stdio redirect', () => {
    // Given
    child_process.execSync.mockImplementationOnce(() => {});
    // When
    exec.logExecSync('1337');
    // Then
    expect(child_process.execSync).toHaveBeenCalledTimes(1);
    expect(child_process.execSync).toHaveBeenCalledWith('1337', {
      stdio: 'inherit'
    });
  });
});
