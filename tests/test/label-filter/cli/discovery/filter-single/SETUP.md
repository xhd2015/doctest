# Scenario

**Feature**: single-label filter runs all leaves with that label

```
--label slow -> runs slow + both
```

## Steps

1. Run `doctest test <mod> --label slow`.

```go
func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "slow"}
	return nil
}
```