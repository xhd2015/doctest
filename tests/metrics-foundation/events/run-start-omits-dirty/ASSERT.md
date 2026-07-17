## Expected

- Decoded run_start has `git_branch` and `git_commit`.
- None of the keys `dirty`, `git_dirty`, `dirty_worktree`, `is_dirty` are present.
- Raw JSON line does not contain `"dirty"` as a key substring pattern `"dirty"`.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if len(resp.Decoded) != 1 {
		t.Fatalf("want 1 event, got %d", len(resp.Decoded))
	}
	m := resp.Decoded[0]
	if m["type"] != "run_start" {
		t.Fatalf("type = %v", m["type"])
	}
	if m["git_branch"] != "feature/metrics" {
		t.Fatalf("git_branch = %v", m["git_branch"])
	}
	if m["git_commit"] != "deadbeefcafebabe" {
		t.Fatalf("git_commit = %v", m["git_commit"])
	}
	for _, k := range []string{"dirty", "git_dirty", "dirty_worktree", "is_dirty"} {
		if _, ok := m[k]; ok {
			t.Fatalf("run_start must not include field %q", k)
		}
	}
	// belt-and-suspenders on raw line keys
	raw := resp.Lines[0]
	for _, frag := range []string{`"dirty"`, `"git_dirty"`, `"dirty_worktree"`, `"is_dirty"`} {
		if strings.Contains(raw, frag) {
			t.Fatalf("raw run_start contains forbidden key fragment %s: %s", frag, raw)
		}
	}
}
```
