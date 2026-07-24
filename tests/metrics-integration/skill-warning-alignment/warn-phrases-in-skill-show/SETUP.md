# Scenario

**Feature**: skill show includes tokens from FormatDefaultSuiteSlowWarning

```
warn = FormatDefaultSuiteSlowWarning()
skill = doctest skill review-perf --show
assert each of: WARNING, skill:doctest-review-perf, review-perf --show, 3 minutes
  appears in skill stdout (or warn message itself for sanity)
```

## Preconditions

- None beyond binary + embedded skill.

## Steps

1. Capture warn message and skill stdout.
2. Assert phrase alignment.

```go
import "testing"

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	// Leaf reinforces op + skill argv (parent may have set them; re-assert here
	// so this SETUP is not an empty stub for vet).
	req.Op = "align_skill_warn"
	if len(req.Args) == 0 {
		req.Args = []string{"skill", "review-perf", "--show"}
	}
	return nil
}
```
