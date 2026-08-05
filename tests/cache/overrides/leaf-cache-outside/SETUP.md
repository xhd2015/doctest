# Scenario

**Feature**: outside leaf-cache override is listed in clean dry-run plan

```
CacheHome/doctest/mapping-gen (seeded)
+ LeafCache outside CacheHome (seeded)
  -> cache --clean --dry-run
  -> [dry-run] would remove: <DoctestRoot> ...
  -> [dry-run] would remove: <LeafCache> ...
  -> neither path deleted
```

## Preconditions

- `LeafCache` absolute path is **not** under `DoctestRoot`.
- Product treats an explicit leaf-cache override (env `DOCTEST_LEAF_CACHE` in
  production; `Request.LeafCache` injection in tests) as an extra clean target
  when it lies outside the main doctest root.

## Steps

1. Seed one in-root bucket.
2. Create an outside leaf-cache directory with a payload; set `req.LeafCache`.
3. Set Args to `cache --clean --dry-run`.

```go
import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	_ = d
	ensureCacheHome(t, req)
	seedBucket(t, req, "mapping-gen", seedBytes(400))

	// Outside store: sibling temp dir, not under CacheHome/doctest.
	outside := t.TempDir()
	// Ensure outside is not under DoctestRoot (TempDir is independent).
	if strings.HasPrefix(outside, req.DoctestRoot) {
		t.Fatalf("outside path unexpectedly under DoctestRoot: %s", outside)
	}
	if err := os.MkdirAll(outside, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(outside, "pass.json"), []byte(`{"ok":true}`), 0o644); err != nil {
		t.Fatal(err)
	}
	req.LeafCache = outside
	req.Args = []string{"cache", "--clean", "--dry-run"}
	return nil
}
```
