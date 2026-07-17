# Scenario

**Feature**: metrics help documents subcommands; unknown subcommand fails

```
# help
doctest metrics --help -> list path last top summary show prune

# unknown
doctest metrics nosuch -> non-zero exit
```

## Preconditions

- Help does not require MetricsRoot fixtures.

## Steps

1. Leaf sets Args for help or unknown subcommand.
2. Run CLI from an arbitrary WorkDir.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// Help leaves do not need metrics fixtures; still give a cwd.
	if req.WorkDir == "" {
		req.WorkDir = t.TempDir()
	}
	return nil
}
```
