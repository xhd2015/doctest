# Scenario

**Feature**: AND expression requires every listed label on the leaf

```
--label 'slow && ui-automation' -> runs both only
```

## Steps

1. Run AND filter against fixture mod.

```go
func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "slow && ui-automation"}
	return nil
}
```