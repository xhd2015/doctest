# Scenario

**Feature**: unified manifest replaces standalone `doctest.gomod-fp`

```
WriteGoMod -> doctest.gomod-fp absent
```

## Steps

1. Inherit cold gen root + single WriteGoMod.
2. Confirm Mode ready; Assert will require gomod-fp absence.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	if req.Mode != "write-gomod" {
		t.Fatalf("no-gomod-fp expects Mode write-gomod, got %q", req.Mode)
	}
	if req.GenDir == "" {
		t.Fatal("parent Setup must prepare GenDir")
	}
	return nil
}
```
