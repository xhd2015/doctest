# Scenario

**Feature**: `--changed` fails outside a git repository

```
# no .git directory
doctest <subcmd> --changed <dir> -> hard error (requires git repo)
```

## Preconditions

- A valid fixture test tree exists in a temp directory.
- The temp directory has **no** git repository initialized.

## Steps

1. Create a fixture test tree without `git init`.
2. Run the subcommand with `--changed` from that directory.

```go
import (
	"path/filepath"
	"testing"

	"github.com/xhd2015/doctest/libdoc/testtree"
)

func Setup(t *testing.T, req *Request) error {
	dir := t.TempDir()
	testtree.WritePassFailTree(t, filepath.Join(dir, "tests"), 2, 0)
	req.WorkDir = dir
	return nil
}
```