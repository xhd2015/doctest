## Expected

- `doctest test` exits 0.
- Assert-mod cache entry list is identical before and after run.

## Exit Code

- Exit code 0.

```go
import "testing"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("doctest subprocess failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected no-import cache skip test to pass, got exit %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
	after := listAssertModCacheEntries(t)
	if len(after) != len(cacheEntriesBefore) {
		t.Fatalf("assert-mod cache entry count changed: before=%v after=%v", cacheEntriesBefore, after)
	}
	for i := range cacheEntriesBefore {
		if after[i] != cacheEntriesBefore[i] {
			t.Fatalf("assert-mod cache entries changed: before=%v after=%v", cacheEntriesBefore, after)
		}
	}
}
```