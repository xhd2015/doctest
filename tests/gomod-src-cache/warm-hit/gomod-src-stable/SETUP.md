# Scenario

**Feature**: warm hit leaves doctest.gomod-src content unchanged

```
snapshot gomod-src -> second identical write -> content equal
```

## Steps

1. Inherit warm-hit Setup (snapshots gomod-src content).
2. Confirm non-empty snapshot before measured call.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.Mode != "write-second" {
		t.Fatalf("gomod-src-stable expects Mode write-second, got %q", req.Mode)
	}
	if req.SnapGomodSrcBefore == "" {
		t.Fatal("parent warm Setup must snapshot gomod-src content")
	}
	return nil
}
```
