# Imports Process: Library-based import formatting

Tests that `WriteGeneratedCase` correctly uses `golang.org/x/tools/imports.Process()` instead of the `goimports` binary, removing the external binary dependency.

## Decision Tree

```
tests/imports-process/
├── DOCTEST.md                      # This file
├── SETUP.md                        # Root: builds doctest binary, run helpers
├── unused-import-removed/          # R1: SETUP.md has unused import → removed by imports.Process
│   ├── SETUP.md                    # Creates test tree with unused import, runs doctest test
│   └── ASSERT.md                   # Assert exit code 0, generated code compiles
└── syntax-error-reported/          # R2: SETUP.md has syntax error → imports.Process error reported
    ├── SETUP.md                    # Creates test tree with broken Go code
    └── ASSERT.md                   # Assert error message, no corrupted file written
```

## Test Index

| Leaf | Description |
|------|-------------|
| `unused-import-removed` | SETUP.md imports `"fmt"` but doesn't use it → `imports.Process` removes it → test compiles and passes |
| `syntax-error-reported` | SETUP.md contains invalid Go (unclosed string) → `imports.Process` fails → clean error reported |

## How to Run

```sh
doctest test -v ./tests/imports-process/
```
