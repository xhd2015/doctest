# Scenario

**Feature**: suite organization / shared setup

```
suite organization
```

## Preconditions
- The path_resolve package is importable from this test tree.
- The root Run dispatches to the function under test based on runType.

## Steps
1. Each leaf sets input fields and runType via `Setup`.
2. The root Run calls the corresponding function.
3. Each leaf asserts results via `Assert`.

## Context
- These are direct unit tests, not CLI integration tests.
- The root Run dispatches to the appropriate test function.

```go
import (
	"testing"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
)

```
