# Scenario

**Feature**: repeatable `--label` flags combine with OR

```
--label slow --label heavy  ≡  --label 'slow || heavy'
```

## Steps

1. Run with two `--label` flags.

```go
func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "slow", "--label", "heavy"}
	return nil
}
```