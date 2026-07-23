## Expected

- Under gen-dir (or work tree), a nested `go.mod` contains a replace for
  `github.com/xhd2015/doctest/session` pointing at the session-mod cache path.
- Cache dir layout is valid.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	cacheDir := expectedSessionCacheDir(t)
	if req.GenDir == "" {
		t.Fatal("req.GenDir is empty; Setup must set request-local gen dir")
	}
	// Search req.GenDir and req.ModuleRoot for a go.mod with session replace.
	roots := []string{req.GenDir, req.ModuleRoot}
	var goModText string
	var goModPath string
	for _, root := range roots {
		_ = filepath.WalkDir(root, func(path string, d os.DirEntry, walkErr error) error {
			if walkErr != nil || d.IsDir() {
				return walkErr
			}
			if d.Name() != "go.mod" {
				return nil
			}
			b, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			s := string(b)
			if strings.Contains(s, sessionModPath) && strings.Contains(s, "replace") {
				goModText = s
				goModPath = path
				return filepath.SkipAll
			}
			return nil
		})
		if goModText != "" {
			break
		}
	}
	if goModText == "" {
		t.Fatalf("no go.mod with session replace found under req.GenDir/req.ModuleRoot\nstdout:\n%s\nstderr:\n%s\nrunErr=%v exit=%d",
			resp.Stdout, resp.Stderr, err, resp.ExitCode)
	}
	if !strings.Contains(goModText, "replace "+sessionModPath) &&
		!strings.Contains(goModText, "replace\t"+sessionModPath) &&
		!strings.Contains(goModText, sessionModPath+" =>") {
		t.Fatalf("go.mod %s missing session replace:\n%s", goModPath, goModText)
	}
	// Replace target should mention cache path segment session-mod.
	if !strings.Contains(goModText, "session-mod") && !strings.Contains(goModText, cacheDir) {
		t.Fatalf("replace should point at session-mod cache; go.mod:\n%s\ncache=%s", goModText, cacheDir)
	}
	assertSessionCacheLayout(t, cacheDir)
}
```
