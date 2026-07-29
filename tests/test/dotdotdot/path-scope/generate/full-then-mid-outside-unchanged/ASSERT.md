---
label: heavy
---

## Expected

After full `./tree/...` then mid `./tree/mid/...` (both bypass go test), every
non-bookkeeping file outside `tree/mid/` is byte-identical to the post-full
snapshot.

```go
import (
	"crypto/sha256"
	"encoding/hex"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func pathScopeHashGen(t *testing.T, root string) map[string]string {
	t.Helper()
	out := make(map[string]string)
	err := filepath.Walk(root, func(p string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		rel, err := filepath.Rel(root, p)
		if err != nil {
			return err
		}
		b, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		sum := sha256.Sum256(b)
		out[filepath.ToSlash(rel)] = hex.EncodeToString(sum[:])
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	return out
}

func pathScopeIsBookkeeping(rel string) bool {
	switch rel {
	case "go.mod", "go.sum", "doctest.gen-manifest", "doctest.tidy-done":
		return true
	default:
		return false
	}
}

func pathScopeUnderMid(rel string) bool {
	return rel == "tree/mid" || strings.HasPrefix(rel, "tree/mid/")
}

func pathScopeRunBin(t *testing.T, req *Request, args []string) (exit int, out string) {
	t.Helper()
	cmd := exec.Command(req.Bin, args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), req.Env...)
	b, err := cmd.CombinedOutput()
	out = string(b)
	if err == nil {
		return 0, out
	}
	if ee, ok := err.(*exec.ExitError); ok {
		return ee.ExitCode(), out
	}
	t.Fatalf("run: %v\n%s", err, out)
	return 1, out
}

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatal(err)
	}
	// Ignore default Run; drive full then mid ourselves for a clean two-phase.
	genDir := filepath.Join(req.WorkDir, ".gen")
	if err := os.RemoveAll(genDir); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(genDir, 0755); err != nil {
		t.Fatal(err)
	}
	exit, out := pathScopeRunBin(t, req, []string{"test", "-v", "--gen-dir", genDir, "-count=1", "./tree/..."})
	if exit != 0 {
		t.Fatalf("full gen exit=%d\n%s", exit, out)
	}
	before := pathScopeHashGen(t, genDir)

	exit, out = pathScopeRunBin(t, req, []string{"test", "-v", "--gen-dir", genDir, "-count=1", "./tree/mid/..."})
	if exit != 0 {
		t.Fatalf("mid gen exit=%d\n%s", exit, out)
	}
	after := pathScopeHashGen(t, genDir)

	var rewritten []string
	for rel, sum := range before {
		if pathScopeIsBookkeeping(rel) || pathScopeUnderMid(rel) {
			continue
		}
		if after[rel] != sum {
			rewritten = append(rewritten, rel)
		}
	}
	for rel := range after {
		if pathScopeIsBookkeeping(rel) || pathScopeUnderMid(rel) {
			continue
		}
		if _, ok := before[rel]; !ok {
			rewritten = append(rewritten, rel+" (new outside mid)")
		}
	}
	if len(rewritten) > 0 {
		t.Fatalf("mid gen rewrote/created content outside tree/mid:\n  %s", strings.Join(rewritten, "\n  "))
	}
}
```
