# Scenario

**Feature**: parent go.sum create/change invalidates gomod-src cache

```
first write without parent go.sum -> seed tidy-done
create/change parent go.sum
second write
  -> miss; gen go.sum written; tidy-done dropped
```

## Steps

1. Parent go.mod only (no go.sum) for first write.
2. Set ChangeSourceGoSum so Run writes parent go.sum before second call.

```go
import (
	"path/filepath"
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-second"
	req.ModPath = defaultModPath
	req.HasMod = true
	prepareFreshGen(t, req, "module example.com/app\n\ngo 1.21\n")
	// Ensure no parent go.sum on first write.
	_ = filepath.Join(req.ModRoot, "go.sum")
	firstWrite(t, req)
	seedTidyDone(t, req.GenDir)
	if fileExists(filepath.Join(req.GenDir, "go.sum")) {
		t.Fatal("precondition: first write should not create gen go.sum without parent go.sum")
	}
	req.ChangeSourceGoSum = "example.com/dep v1.0.0 h1:abc=\nexample.com/dep v1.0.0/go.mod h1:def=\n"
	return nil
}
```
