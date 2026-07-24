## Preconditions
- `sh` is available to inspect environment variables.

## Steps
1. Run with `--model=gpt-4` and print `DOCTEST_SUBAGENT_MODEL`.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
    req.Args = append(req.Args, "--model=gpt-4", "sh", "-c", "echo $DOCTEST_SUBAGENT_MODEL")
    return nil
}
```
