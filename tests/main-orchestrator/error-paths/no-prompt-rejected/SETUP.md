# Scenario

**Feature**: `doctest agent implement` requires a prompt (**L2 short path**)

```
doctest agent implement --agent-runner fake-codex  # no prompt
  -> non-zero exit; stderr mentions "requires"
```

## Preconditions

- Missing-prompt is a fast-fail short path — in-process CLI, no binary/fake-codex.

## Steps

1. Clear parent e2e Env (Env requires UseCLI under Parallel rules).
2. Run `agent implement` with no prompt via in-process CLI.
3. Expect error.

```go
import (
	"testing"
)

func Setup(t *testing.T, req *Request) error {
	// L2 short path: no Env isolation, no product binary.
	req.UseCLI = false
	req.Bin = ""
	req.Env = nil
	req.Args = []string{"agent", "implement", "--agent-runner", "fake-codex"}
	return nil
}
```
