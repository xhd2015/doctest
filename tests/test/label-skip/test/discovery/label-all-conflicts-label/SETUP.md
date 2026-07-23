# Scenario

**Feature**: --label-all rejects combination with --label (L2 parse)

```
runner.ParseTestOptions([--label-all, --label, heavy, .])
  -> mutually exclusive
```

## Steps

1. Parse with both flags (dir unused after parse fails).

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Args = []string{"--label-all", "--label", "heavy", "."}
	return nil
}
```
