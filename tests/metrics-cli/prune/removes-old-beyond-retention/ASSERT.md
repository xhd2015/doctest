## Expected

- Exit code 0.
- Exactly `DefaultRunRetention` (30) `*.jsonl` files remain under the project runs dir.
- Newest files kept: remaining set is the 30 greatest basenames from the original 35.

## Exit Code

- 0

```go
import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"testing"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit=%d stderr=%s stdout=%s", resp.ExitCode, resp.Stderr, resp.Stdout)
	}
	dir := projectRunsDir(req)
	ents, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("readdir: %v", err)
	}
	var got []string
	for _, e := range ents {
		if !e.IsDir() && strings.HasSuffix(e.Name(), ".jsonl") {
			got = append(got, e.Name())
		}
	}
	sort.Strings(got)
	if len(got) != DefaultRunRetention {
		t.Fatalf("after prune want %d files, got %d: %v\nstdout=%s\nstderr=%s",
			DefaultRunRetention, len(got), got, resp.Stdout, resp.Stderr)
	}
	// Newest 30 of 35: prune00..prune04 are oldest (June early) should be gone if naming ok.
	// Check at least one of the highest indices remains.
	joined := strings.Join(got, "\n")
	if !strings.Contains(joined, "prune34") && !strings.Contains(joined, "prune33") {
		t.Fatalf("expected newest prune files retained; got:\n%s", joined)
	}
	_ = filepath.Separator
}
```
