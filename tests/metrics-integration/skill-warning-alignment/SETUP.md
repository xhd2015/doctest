# Scenario

**Feature**: WARNING banner phrases are present in the review-perf skill body

```
FormatDefaultSuiteSlowWarning()
  -> required phrases
skill review-perf --show
  -> same phrases appear (banner does not point at missing docs)
```

## Preconditions

- P2 warn formatter and P4 skill content exist.
- This is a thin cross-link; full skill content is covered under `tests/skill/review-perf-show`.

## Steps

1. Leaf sets `Op=align_skill_warn`.
2. Run loads warn message + skill show.
3. Assert each required warn phrase is a substring of skill stdout.

## Context

- Required phrases match `metrics-record/warning-message/required-phrases`:
  `WARNING:`, `skill:doctest-review-perf`, `review-perf --show`, `3 minutes`.
  Skill body may phrase WARNING as prose (**WARNING**) rather than the exact
  banner; we require the same product tokens the banner uses so agents can
  follow the recommendation.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	req.Op = "align_skill_warn"
	req.Args = []string{"skill", "review-perf", "--show"}
	return nil
}
```
