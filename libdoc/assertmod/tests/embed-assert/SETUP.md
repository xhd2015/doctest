# Scenario

**Feature**: build-time embedding concatenates assert package into a single deterministic file

```
# script/generate/embed-assert reads assert/*.go
go run script/generate/embed-assert -> libdoc/assertmod/assert.go

# libdoc/assertmod exposes ContentMD5
embed bytes -> ContentMD5() == md5(assert.go on disk)
```

## Preconditions

- Module root is four levels above this tree (`DOCTEST_ROOT/../../../..`).
- `assert/` directory contains production `.go` sources (no `*_test.go` in output).

## Steps

1. Descendant sets `runKind` and optional `SecondRun` for embed script tests.

## Context

- `runEmbedScript` invokes `go run <moduleRoot>/script/generate/embed-assert` with `-o` output path.

```go
import (
"github.com/xhd2015/doctest/session"
	"crypto/md5"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/libdoc/assertmod"
)

var runKind string

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.ModuleRoot = filepath.Join(d.DOCTEST_ROOT, "..", "..", "..", "..")
	req.AssertDir = filepath.Join(req.ModuleRoot, "assert")
	return nil
}
func runEmbedScript(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	outPath := req.OutputPath
	if outPath == "" {
		outPath = filepath.Join(t.TempDir(), "assert.go")
	}
	scriptPkg := filepath.Join(req.ModuleRoot, "script", "generate", "embed-assert")
	cmd := exec.Command("go", "run", scriptPkg, "-o", outPath, req.AssertDir)
	cmd.Dir = req.ModuleRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return &Response{Err: fmt.Errorf("embed script: %w\n%s", err, string(out))}, err
	}
	data, err := os.ReadFile(outPath)
	if err != nil {
		return nil, err
	}
	sum := md5.Sum(data)
	resp := &Response{
		OutputBytes:	data,
		OutputMD5:	fmt.Sprintf("%x", sum),
	}
	if req.SecondRun {
		outPath2 := filepath.Join(t.TempDir(), "assert2.go")
		cmd2 := exec.Command("go", "run", scriptPkg, "-o", outPath2, req.AssertDir)
		cmd2.Dir = req.ModuleRoot
		if out, err := cmd2.CombinedOutput(); err != nil {
			return &Response{Err: fmt.Errorf("second embed script: %w\n%s", err, string(out))}, err
		}
		data2, err := os.ReadFile(outPath2)
		if err != nil {
			return nil, err
		}
		sum2 := md5.Sum(data2)
		resp.SecondRunMD5 = fmt.Sprintf("%x", sum2)
	}
	return resp, nil
}
func runEmbedCacheKey(t *testing.T, req *Request) (*Response, error) {
	t.Helper()
	tmp := t.TempDir()
	outPath := filepath.Join(tmp, "assert.go")
	keyPath := filepath.Join(tmp, "cache_key.go")
	scriptPkg := filepath.Join(req.ModuleRoot, "script", "generate", "embed-assert")
	cmd := exec.Command("go", "run", scriptPkg, "-o", outPath, "-cache-key", keyPath, req.AssertDir)
	cmd.Dir = req.ModuleRoot
	if out, err := cmd.CombinedOutput(); err != nil {
		return &Response{Err: fmt.Errorf("embed cache key: %w\n%s", err, string(out))}, err
	}
	data, err := os.ReadFile(keyPath)
	if err != nil {
		return &Response{Err: err}, err
	}
	scriptKey, err := parseRawSourceCacheKey(string(data))
	if err != nil {
		return &Response{Err: err}, err
	}
	return &Response{
		ScriptCacheKey:	scriptKey,
		PackageCacheKey: assertmod.RawSourceCacheKeyMD5(),
	}, nil
}
func parseRawSourceCacheKey(src string) (string, error) {
	const prefix = `const rawSourceCacheKeyMD5 = "`
	i := strings.Index(src, prefix)
	if i < 0 {
		return "", fmt.Errorf("cache key constant not found")
	}
	rest := src[i+len(prefix):]
	j := strings.Index(rest, `"`)
	if j < 0 {
		return "", fmt.Errorf("cache key string not terminated")
	}
	return rest[:j], nil
}
```