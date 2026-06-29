# Scenario

**Feature**: expression matching no leaves exits 0 with full skip summary

```
--label manual -> all leaves skipped, no PASS line
```

## Steps

1. Run filter that matches nothing.

```go
func Setup(t *testing.T, req *Request) error {
	mod := writeLabelFilterMod(t)
	req.Args = []string{"test", mod, "--label", "manual"}
	return nil
}
```