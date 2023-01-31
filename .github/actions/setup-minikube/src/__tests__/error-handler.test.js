describe('error-handler module test suite', () => {
  let errorHandler;
  let core;
  beforeEach(() => {
    jest.resetModules();
    jest.mock('@actions/core');
    errorHandler = require('../error-handler');
    core = require('@actions/core');
  });
  test('errorHandler, should set action failed', () => {
    // Given
    console.error = jest.fn(() => {});
    core.setFailed.mockImplementationOnce(() => {});
    // When
    errorHandler(Error('Something bad happened'));
    // Then
    expect(console.error).toHaveBeenCalledWith(Error('Something bad happened'));
    expect(core.setFailed).toHaveBeenCalledTimes(1);
    expect(core.setFailed).toHaveBeenCalledWith('Something bad happened');
  });
});
