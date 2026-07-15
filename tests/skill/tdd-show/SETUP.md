# Scenario

**Feature**: the doc-style test based TDD specification exists under `agents/doctest/doc`

```
# expose embedded spec documents
doctest skill <name> --show -> embedded .md doc -> stdout

# list available skills
doctest skill --list -> skill names -> stdout
```

## Preconditions
- The doc-style test based TDD specification exists under `agents/doctest/doc`.

## Steps
1. Run `doctest skill tdd --show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "tdd", "--show"}
    return nil
}
```
