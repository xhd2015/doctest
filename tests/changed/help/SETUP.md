# Scenario

**Feature**: `--changed` in subcommand help via in-process CLI (no product binary)

```
cli.RunWithWriter -> doctest <subcmd> --help -> stdout lists --changed
```

## Preconditions

- Help is covered in-process via `cli.RunWithWriter` (same usage strings as the product binary).
- Unlabeled (fast); no `testbin`, no `label: heavy`.
- Policy selection stays in-process under `git-context/in-git-repo/`.

## Steps

1. Root help Setup is a no-op (no binary).
2. Leaf sets `Args` to `<subcmd> --help`.
3. `Run` calls `cli.RunWithWriter` when Args are set without TreeDir.

## Context

- No fixture tree required for help.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// In-process CLI: no testbin, no UseCLI binary path.
	_ = d
	_ = req
	return nil
}
```
