# Scenario

**Feature**: manifest hash hit skips target rewrite (mtime stable)

```
write rel once -> force old mtime -> write same content again
  -> target mtime unchanged
  -> manifest still lists rel
```

## Steps

1. Inherit write-if-changed Setup (gen root + RelPath).
2. First write FileContent in Setup; snapshot target mtime.
3. Run Mode `write-file-second` with same content.

```go
import (
	"testing"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Mode = "write-file-second"
	if err := writeGenRelFile(t, req.GenDir, req.RelPath, req.FileContent); err != nil {
		t.Fatalf("first writeGenRelFile: %v", err)
	}
	snapshotTargetMtime(t, req)
	// Second write uses same FileContent (SecondFileContent empty).
	return nil
}
```
