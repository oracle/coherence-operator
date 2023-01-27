describe('check-environment module test suite', () => {
  let checkEnvironment;
  let fs;
  beforeEach(() => {
    jest.resetModules();
    jest.mock('fs');
    checkEnvironment = require('../check-environment');
    fs = require('fs');
  });
  describe('checkEnvironment', () => {
    test('OS is windows, should throw Error', () => {
      // Given
      Object.defineProperty(process, 'platform', {value: 'win32'});
      process.platform = 'win32';
      // When - Then
      expect(checkEnvironment).toThrow(
        'Unsupported OS, action only works in Ubuntu 18, 20, or 22'
      );
    });
    test('OS is Linux but not Ubuntu, should throw Error', () => {
      // Given
      Object.defineProperty(process, 'platform', {value: 'linux'});
      fs.existsSync.mockImplementation(() => false);
      fs.readFileSync.mockImplementation(() => 'SOME DIFFERENT OS');
      // When - Then
      expect(checkEnvironment).toThrow(
        'Unsupported OS, action only works in Ubuntu 18, 20, or 22'
      );
      expect(fs.existsSync).toHaveBeenCalled();
      expect(fs.readFileSync).toHaveBeenCalledTimes(0);
    });
    test('OS is Linux but not Ubuntu 18, should throw Error', () => {
      // Given
      Object.defineProperty(process, 'platform', {value: 'linux'});
      fs.existsSync.mockImplementation(() => true);
      fs.readFileSync.mockImplementation(() => 'SOME DIFFERENT OS');
      // When - Then
      expect(checkEnvironment).toThrow(
        'Unsupported OS, action only works in Ubuntu 18, 20, or 22'
      );
      expect(fs.existsSync).toHaveBeenCalled();
      expect(fs.readFileSync).toHaveBeenCalled();
    });
    test('OS is Linux and Ubuntu 18, should not throw Error', () => {
      // Given
      Object.defineProperty(process, 'platform', {value: 'linux'});
      fs.existsSync.mockImplementation(() => true);
      fs.readFileSync.mockImplementation(
        () => `
        NAME="Ubuntu"
        VERSION="18.04.3 LTS (Bionic Beaver)"
        `
      );
      // When - Then
      expect(checkEnvironment).not.toThrow();
    });
    test('OS is Linux and Ubuntu 20, should not throw Error', () => {
      // Given
      Object.defineProperty(process, 'platform', {value: 'linux'});
      fs.existsSync.mockImplementation(() => true);
      fs.readFileSync.mockImplementation(
        () => `
        NAME="Ubuntu"
        VERSION="20.04.1 LTS (Focal Fossa)"
        `
      );
      // When - Then
      expect(checkEnvironment).not.toThrow();
    });
    test('OS is Linux and Ubuntu 22, should not throw Error', () => {
      // Given
      Object.defineProperty(process, 'platform', {value: 'linux'});
      fs.existsSync.mockImplementation(() => true);
      fs.readFileSync.mockImplementation(
        () => `
        NAME="Ubuntu"
        VERSION="22.04.1 LTS (Jammy Jellyfish)"
        `
      );
      // When - Then
      expect(checkEnvironment).not.toThrow();
    });
  });
});
