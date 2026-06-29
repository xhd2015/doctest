# Scenario

**Feature**: trailing operator in label expression is rejected

```
doctest test <mod> --label 'slow &&' -> stderr parse error
```

## Steps

1. Create fixture mod.
2. Run with invalid `--label`.

```go
func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "slow &&"}
	return nil
}
```