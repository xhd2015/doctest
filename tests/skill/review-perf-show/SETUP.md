# Scenario

**Feature**: the doctest-review-perf skill exists and is registered as `review-perf`

```
# expose embedded spec documents
doctest skill <name> --show -> embedded .md doc -> stdout

# list available skills
doctest skill --list -> skill names -> stdout
```

## Preconditions
- The review-perf skill document is registered and embeddable (CLI name `review-perf`, skill name `doctest-review-perf`).

## Steps
1. Run `doctest skill review-perf --show`.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
    req.Args = []string{"skill", "review-perf", "--show"}
    return nil
}
```
