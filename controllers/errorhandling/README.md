# Error Handling in Coherence Operator

This package provides standardized error handling patterns for the Coherence Operator. It includes utilities for error wrapping, context addition, and common error scenarios.

## Key Components

### Error Types

- **OperationError**: Represents an error that occurred during an operation. It includes:
  - Operation name
  - Resource name and namespace (optional)
  - Underlying error
  - Context map for additional information

### Error Creation

- **NewOperationError**: Creates a new operation error
- **NewResourceError**: Creates an error for a specific resource

### Common Error Scenarios

The package provides helper functions for common error scenarios:

- **NewCreateResourceError**: For resource creation failures
- **NewUpdateResourceError**: For resource update failures
- **NewDeleteResourceError**: For resource deletion failures
- **NewGetResourceError**: For resource retrieval failures
- **NewListResourceError**: For resource listing failures
- **NewPatchResourceError**: For resource patching failures
- **NewReconcileError**: For reconciliation failures
- **NewValidationError**: For validation failures
- **NewTimeoutError**: For timeout failures
- **NewConnectionError**: For connection failures
- **NewAuthenticationError**: For authentication failures
- **NewAuthorizationError**: For authorization failures

### Error Wrapping

- **WrapError**: Wraps an error with context information
- **WrapErrorf**: Wraps an error with formatted context information
- **WithStack**: Adds a stack trace to an error

### Error Handling

- **ErrorHandler**: Handles errors in the reconciliation loop
  - Categorizes errors (Transient, Permanent, Recoverable, Unknown)
  - Updates error tracking information
  - Updates resource status
  - Handles errors based on their category

## Usage Examples

### Creating Errors

```go
// Create a simple operation error
err := errorhandling.NewOperationError("update_config", originalErr)

// Add context to the error
err.WithContext("resource_type", "ConfigMap").WithContext("retry_count", "3")

// Create a resource-specific error
err := errorhandling.NewResourceError("update", "my-configmap", "default", originalErr)

// Use helper functions for common scenarios
err := errorhandling.NewCreateResourceError("my-pod", "default", originalErr)
```

### Handling Errors

```go
// Create an error handler
errorHandler := errorhandling.NewErrorHandler(client, logger, recorder)

// Handle an error
result, err := errorHandler.HandleError(ctx, originalErr, resource, "Failed to update resource")

// Handle a resource-specific error
result, err := errorHandler.HandleResourceError(ctx, originalErr, resource, "update", "Failed to update resource")

// Retry an operation with context
err := errorHandler.RetryWithContext(ctx, "update", "my-configmap", "default", func() error {
    // Operation that might fail
    return client.Update(ctx, configMap)
})
```

## Best Practices

1. **Always add context to errors**: Use WithContext to add relevant information to errors.
2. **Use the appropriate helper function**: Choose the most specific helper function for your error scenario.
3. **Handle errors based on their category**: Use the ErrorHandler to handle errors appropriately based on their category.
4. **Include stack traces**: Use WithStack to add stack traces to errors for better debugging.
5. **Log errors with context**: Use the LogAndWrapError functions to log errors with context.