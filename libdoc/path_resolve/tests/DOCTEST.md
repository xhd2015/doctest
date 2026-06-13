# path_resolve Package Tests

These doc-style tests specify the contract for the `path_resolve` package.
Each function (`IsDotDotDotPattern`, `ExtractBasePath`, `ResolveRoot`,
`FindDotDotDotDirs`) is tested through direct function calls, not CLI
integration.

Tests are organized by function, with leaves covering all inputs and edge cases.
