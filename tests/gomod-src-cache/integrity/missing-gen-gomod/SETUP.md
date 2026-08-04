# Scenario

**Feature**: missing gen go.mod forces rebuild despite matching fingerprint file

```
first write -> leave gomod-src intact
delete genDir/go.mod
second write
  -> rebuild succeeds; go.mod restored
```

## Steps

1. Fresh dirs; first write; keep fingerprint on disk.
2. Set DeleteGenGoMod so Run removes gen go.mod before second call.

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
	firstWrite(t, req)
	if !fileExists(filepath.Join(req.GenDir, gomodSrcName)) {
		t.Fatal("expected gomod-src after first write")
	}
	req.DeleteGenGoMod = true
	return nil
}
```
