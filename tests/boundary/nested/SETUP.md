# Scenario

**Feature**: this is a nested root with its own DOCTEST.md boundary

```
# DOCTEST.md creates an inheritance firewall
parent SETUP.md -/-> nested DOCTEST.md (no cross-inheritance)

# each root has its own Run, Request, Response, setup chain
nested root -> self-contained test tree -> runs independently
```

## Preconditions
- This is a nested root with its own DOCTEST.md boundary.
- This root defines its own Request and Response types (different from parent root).

## Steps
1. The leaf sets `req.Name`.
2. Run returns a greeting using the Name field, or an error if Name is empty.

```go
import (
	"fmt"
	"testing"
)

```
