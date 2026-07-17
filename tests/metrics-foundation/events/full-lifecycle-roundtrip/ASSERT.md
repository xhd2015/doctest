## Expected

- File has 5 non-empty JSONL lines, each ending with newline in raw file.
- Types in order: `run_start`, `leaf_start`, `leaf_end`, `leaf_end`, `run_end`.
- Each decoded object has `schema_version` == 1 (number).
- Pass leaf_end has `result=pass`; skip leaf_end has `result=skip`.
- run_end has `passed=1`, `total=2`, `skipped=1`, `exit_ok=true`.

```go
import (
	"bytes"
	"testing"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("unexpected error: %v (resp.Err=%s)", err, resp.Err)
	}
	if len(resp.FileData) == 0 {
		t.Fatal("file empty after close")
	}
	if !bytes.HasSuffix(resp.FileData, []byte("\n")) {
		t.Fatal("JSONL file must end with newline")
	}
	if len(resp.Lines) != 5 {
		t.Fatalf("line count = %d, want 5\nfile:\n%s", len(resp.Lines), resp.FileData)
	}
	if len(resp.Decoded) != 5 {
		t.Fatalf("decoded count = %d, want 5", len(resp.Decoded))
	}

	wantTypes := []string{"run_start", "leaf_start", "leaf_end", "leaf_end", "run_end"}
	for i, want := range wantTypes {
		got, _ := resp.Decoded[i]["type"].(string)
		if got != want {
			t.Fatalf("event[%d].type = %q, want %q", i, got, want)
		}
		sv := resp.Decoded[i]["schema_version"]
		switch v := sv.(type) {
		case float64:
			if int(v) != 1 {
				t.Fatalf("event[%d].schema_version = %v, want 1", i, sv)
			}
		case int:
			if v != 1 {
				t.Fatalf("event[%d].schema_version = %v, want 1", i, sv)
			}
		default:
			t.Fatalf("event[%d].schema_version type %T = %v", i, sv, sv)
		}
	}

	if resp.Decoded[2]["result"] != "pass" {
		t.Fatalf("leaf_end pass result = %v", resp.Decoded[2]["result"])
	}
	if resp.Decoded[3]["result"] != "skip" {
		t.Fatalf("leaf_end skip result = %v", resp.Decoded[3]["result"])
	}

	re := resp.Decoded[4]
	if re["passed"] != float64(1) && re["passed"] != 1 {
		t.Fatalf("run_end.passed = %v", re["passed"])
	}
	if re["total"] != float64(2) && re["total"] != 2 {
		t.Fatalf("run_end.total = %v", re["total"])
	}
	if re["skipped"] != float64(1) && re["skipped"] != 1 {
		t.Fatalf("run_end.skipped = %v", re["skipped"])
	}
	if re["exit_ok"] != true {
		t.Fatalf("run_end.exit_ok = %v", re["exit_ok"])
	}
}
```
