# Tests

This folder contains integration tests. Here is a brief explanation of how Go tests are organized in this project.

## Unit Tests

Unit tests reside, as it is _the Go way_, beside their implementation.

```
|- pkg
|  |- service
|  |  |- users_service.go       <-- implementation of the Users service
|  |  |- users_service_test.go  <-- unit tests of the Users service
```

## Integration Tests

Those tests require more setup and external dependencies (typically: the database).

### Utility functions

Some functions are in place to help with integration tests.

- `GetTestDB` ensures a single database connection is used and, for the duration of this tests
  (accepts the `t *testing.T` as a paramter), opens a temporary tranasaction that is rolled back
  automatically, ensuring the tests will not have modified the database in the end.
- `CreateTestAPIHandler` returns a HTTP handler that is using `GetTestDB` to ensure no call
  to this test API handler modifies the database.
