## Expected

- Warn message is non-empty and contains the standard banner tokens.
- Skill show exits 0.
- Skill stdout contains:
  - `WARNING` (section or prose about the banner)
  - `skill:doctest-review-perf`
  - `review-perf --show`
  - `3 minutes`

## Exit Code

- Skill exit 0.

```go
import (
	"strings"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("harness error: %v", err)
	}
	if resp.WarnMessage == "" {
		t.Fatal("FormatDefaultSuiteSlowWarning returned empty")
	}
	for _, n := range []string{"WARNING:", "skill:doctest-review-perf", "review-perf --show", "3 minutes"} {
		if !strings.Contains(resp.WarnMessage, n) {
			t.Fatalf("warn message missing %q:\n%s", n, resp.WarnMessage)
		}
	}
	if resp.SkillExit != 0 {
		t.Fatalf("skill show exit=%d stderr=%s", resp.SkillExit, resp.SkillStderr)
	}
	skill := resp.SkillStdout
	// Align banner recommendation tokens with skill body.
	needles := []string{
		"WARNING",
		"skill:doctest-review-perf",
		"review-perf --show",
		"3 minutes",
	}
	for _, n := range needles {
		if !strings.Contains(skill, n) {
			t.Fatalf("skill show missing align phrase %q:\n%s", n, skill)
		}
	}
}
```
