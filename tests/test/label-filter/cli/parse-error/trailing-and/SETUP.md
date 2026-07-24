# Scenario

**Feature**: trailing operator in label expression is rejected (L2 parse)

```
runner.ParseTestOptions([., --label, 'slow &&']) -> parse/syntax error
```

## Steps

1. Parse with invalid `--label` (dir unused after parse fails).

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.Args = []string{".", "--label", "slow &&"}
	return nil
}
```
