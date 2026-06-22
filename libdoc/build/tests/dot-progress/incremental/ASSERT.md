## Expected
- `resp.Incremental` is `true` — the first dot appeared within 4s (before
  the slow leaf's 5s sleep completed), proving dots are printed as each
  test package finishes, not batched after all tests complete.
- `resp.DotCount` is `2` — exactly one dot per package (a_fast + z_slow).
- The two dots appear **before** the summary line `"(N Run, ...)"`.
- No dots appear **after** the summary line.
- `err` is nil (all tests pass).

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if !resp.Incremental {
		t.Fatal("dots are NOT printed incrementally — " +
			"the first dot appeared after the slow package finished, " +
			"indicating dots are batched")
	}

	if resp.DotCount != 2 {
		t.Fatalf("expected 2 dots (one per package), got %d. output:\n%s",
			resp.DotCount, resp.Output)
	}

	summaryIdx := strings.Index(resp.Output, "  (")
	if summaryIdx < 0 {
		t.Fatalf("expected summary line starting with '  (' in output:\n%s",
			resp.Output)
	}

	// Dots must be before the summary.
	dotsBefore := strings.Count(resp.Output[:summaryIdx], ".")
	if dotsBefore != 2 {
		t.Fatalf("expected 2 dots before summary, got %d. output:\n%s",
			dotsBefore, resp.Output)
	}

	// No progress-dot lines after the inline summary (duration decimals are OK).
	inlineEnd := strings.Index(resp.Output[summaryIdx:], "\n")
	rest := ""
	if inlineEnd >= 0 {
		rest = resp.Output[summaryIdx+inlineEnd+1:]
	}
	for _, line := range strings.Split(rest, "\n") {
		trimmed := strings.TrimSpace(line)
		if trimmed != "" && strings.Trim(trimmed, ".") == "" {
			t.Fatalf("unexpected progress dots after inline summary: %q", line)
		}
	}
}
```
