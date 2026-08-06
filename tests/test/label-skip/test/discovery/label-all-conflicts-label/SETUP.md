# Scenario

**Feature**: --label-all rejects combination with --label (L2 parse)

```
runner.ParseTestOptions([--label-all, --label, e2e, .])
  -> mutually exclusive
```

## Steps

1. Parse with both flags (dir unused after parse fails).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{"--label-all", "--label", "e2e", "."}
	return nil
}
```
