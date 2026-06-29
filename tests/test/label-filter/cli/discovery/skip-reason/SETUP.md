# Scenario

**Feature**: label-filter skips include `reason: label filter`

```
--label manual -> every skipped entry documents reason
```

## Steps

1. Reuse no-match command; assert reason line present on all skips.

```go
func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "manual"}
	return nil
}
```