# Scenario

**Feature**: OR expression runs leaves with any matching label

```
--label 'slow || heavy' -> runs slow, both, heavy
```

## Steps

1. Run OR filter.

```go
func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "slow || heavy"}
	return nil
}
```