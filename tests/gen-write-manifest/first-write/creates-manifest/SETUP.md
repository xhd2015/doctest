# Scenario

**Feature**: first WriteGoMod creates `doctest.gen-manifest` listing `go.mod`

```
WriteGoMod -> doctest.gen-manifest exists
  go.mod <content-hash>
  version field present
```

## Steps

1. Inherit cold gen root + single WriteGoMod Mode.
2. Confirm Mode and dirs are ready for the first-write create-manifest check.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-gomod" {
		t.Fatalf("creates-manifest expects Mode write-gomod, got %q", req.Mode)
	}
	if req.GenDir == "" || req.ModRoot == "" {
		t.Fatal("parent Setup must prepare GenDir and ModRoot")
	}
	return nil
}
```
