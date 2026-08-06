# DOCTEST_DEBUG gen-plan=1 — generate plan / result trees on stderr

## Version
0.0.2

**Mode: Classic TDD** — `gen-plan` is not implemented yet. Leaves that require
accepted `gen-plan=1` or plan/result stderr markers are expected **RED** until
the implementer lands parse + emit.

# DSN (Domain Specific Notion)

### Participants

- **DOCTEST_DEBUG** — GODEBUG-style process env (`comma key=value`). Fail-closed
  on unknown keys. Engine-internal; not a public CLI flag.
- **gen-plan key** — when `gen-plan=1`, the host prints generate **plan** and
  **result** hierarchy trees on **stderr** only (never pollutes test stdout / JSON).
- **bypass-go-test key** — existing key; skips host go test after generate. Pairs
  with gen-plan for plan-only / prepare-only runs.
- **Path arg** — each CLI path argument to `doctest test` (single tree root or
  multi-arg roots). Plan labels args as `arg[i/n]`.
- **Gen root** — isolated generate directory (`--gen-dir`). Bookkeeping files:
  `go.mod`, `go.sum`, `doctest.gen-manifest`, `doctest.tidy-done`.
- **Tree hierarchy** — relative paths under gen root (packages + dirs) printed
  as an indented file tree.
- **Merged view** — multi-arg / multi-tree plan ends with `gen-plan: merged`
  listing bookkeeping once, all trees, and `__workspace/…` when a hub is written.
- **Result coloring** — after generate + prune, `gen-plan: result` reprints the
  same hierarchy shape with status colors when color is on:
  gray=unchanged, green=modified/created, red=deleted; plus
  `summary: modified=N unchanged=M deleted=K`.

### Behaviors

1. **Parse** — `gen-plan=1` is a known key (bool). Unknown keys still error.
2. **Invocation header** (stderr) — short banner, e.g.
   `doctest: DOCTEST_DEBUG gen-plan=1` and `gen-plan: invocation` with args,
   labels, gen root, mode.
3. **Plan per arg** — `gen-plan: arg[i/n]  <arg>` then that arg’s hierarchy.
   - Single-arg: bookkeeping files appear under the arg[1/1] tree.
   - Multi-arg: each arg shows **only that tree’s package subtree** (no repeated
     go.mod per arg).
4. **Merged** (multi-arg or multi-tree only) — `gen-plan: merged` with
   bookkeeping + all trees + `__workspace` if hub written. Single-tree: no
   separate merged block.
5. **Result** — after generate (+ prune when it runs): `gen-plan: result` with
   the same structure as plan (arg shape for single; merged shape for multi),
   status colors when color on, legend optional, summary counts.
6. **Stdout clean** — gen-plan markers stay on stderr; test stdout/JSON free of
   `gen-plan:` noise.
7. **Combo** — `DOCTEST_DEBUG=gen-plan=1,bypass-go-test=1` accepts both; plan and
   result still print while go test is skipped.

### Pipeline sketch

```
DOCTEST_DEBUG=gen-plan=1[,bypass-go-test=1]
  -> Parse (fail-closed unknown keys)
  -> doctest test <args> --gen-dir <isolated>
       plan phase (stderr):
         doctest: DOCTEST_DEBUG gen-plan=1
         gen-plan: invocation …
         gen-plan: arg[i/n] <arg>
           <hierarchy>
         [gen-plan: merged …]          # multi only
       generate + optional prune
       result phase (stderr):
         gen-plan: result
           <same hierarchy; color statuses>
           summary: modified=N unchanged=M deleted=K
```

## Decision Tree

```
debug-gen-plan/
├── parse/                                 library debug.Parse (L2, unlabeled)
│   ├── accepts-gen-plan/                  gen-plan=1 → Settings.GenPlan
│   ├── accepts-with-bypass/               gen-plan=1,bypass-go-test=1 both on
│   ├── unknown-key-fail-closed/           still errors on unknown keys
│   └── invalid-bool/                      gen-plan=maybe → bool error
├── plan/                                  plan phase on stderr (product CLI)
│   ├── single-arg/
│   │   └── hierarchy-with-bookkeeping/    arg[1/1] + go.mod/manifest in tree
│   └── multi-arg/
│       └── per-arg-and-merged/            arg[1/2] arg[2/2] package-only + merged
└── result/                                result phase after generate
    ├── statuses/
    │   ├── cold/                          first run: gen-plan: result + summary
    │   └── warm/                          second run: unchanged increases
    └── color/
        ├── forced-on/                     --color → green/gray SGR on result
        └── forced-off/                    --no-color → no ANSI in gen-plan lines
```

## Test Index

| Leaf | Layer | Scenario |
|------|-------|----------|
| `parse/accepts-gen-plan` | L2 | `Parse("gen-plan=1")` succeeds; GenPlan true |
| `parse/accepts-with-bypass` | L2 | Both gen-plan and bypass-go-test on |
| `parse/unknown-key-fail-closed` | L2 | `bypass-go-test=1,not-a-key=1` still errors |
| `parse/invalid-bool` | L2 | `gen-plan=maybe` errors with bool message |
| `plan/single-arg/hierarchy-with-bookkeeping` | L3 e2e | stderr plan: arg[1/1], bookkeeping files in hierarchy |
| `plan/multi-arg/per-arg-and-merged` | L3 e2e | arg[1/2]+arg[2/2] without go.mod each; merged has bookkeeping |
| `result/statuses/cold` | L3 e2e | `gen-plan: result` + summary; modified≥1 on cold gen |
| `result/statuses/warm` | L3 e2e | second identical run → more unchanged (summary) |
| `result/color/forced-on` | L3 e2e | `--color`: ANSI green and/or gray on result lines |
| `result/color/forced-off` | L3 e2e | `--no-color`: no ESC sequences on gen-plan lines |

## How to Run

```sh
doctest vet ./tests/debug-gen-plan/
doctest test ./tests/debug-gen-plan/                 # L2 parse only (default labels)
doctest test --label e2e ./tests/debug-gen-plan/   # include plan/result product leaves
# Classic TDD: expect RED on gen-plan acceptance + plan/result until implementer lands feature
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/cli"
	"github.com/xhd2015/doctest/libdoc/debug"
)

// Request drives parse-only or product generate+plan/result scenarios.
type Request struct {
	// Mode: "parse" | "cli" | "cli-twice"
	//   parse     — debug.Parse(DebugEnv) only
	//   cli       — one doctest test invocation with DebugEnv
	//   cli-twice — two invocations sharing GenDir (warm result)
	Mode string

	// DebugEnv is the DOCTEST_DEBUG value (e.g. "gen-plan=1,bypass-go-test=1").
	// Empty → no DOCTEST_DEBUG for subprocess; parse of "" is zero settings.
	DebugEnv string

	// Args for doctest CLI (without binary). Typical: test --gen-dir … --no-color tree
	Args []string
	// Args2 for Mode=cli-twice second run (defaults to Args when empty).
	Args2 []string

	// Fixture / isolation paths filled by Setup helpers.
	WorkDir    string
	GenDir     string
	FixtureDir string // single-tree root
	TreeRoot   string // same as FixtureDir for single; tree-a for multi
	TreeRootB  string // multi-arg second root
	ModuleRoot string

	// Product binary path (testbin.Ensure). Required when DebugEnv is set
	// (subprocess injects env; FromEnv is process-initial).
	Bin     string
	Env     []string // extra child env (DOCTEST_CACHE_HOME, GOWORK=off, …)
	Timeout time.Duration

	// ColorMode: "" | "on" | "off" — leaf may also put --color/--no-color in Args.
	ColorMode string
}

// Response captures parse and/or CLI outcomes.
type Response struct {
	// Parse path
	ParseErr   string
	GenPlan    bool
	BypassGoTest bool

	// CLI path
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error

	// Warm second run (Mode=cli-twice)
	SecondExitCode int
	SecondStdout   string
	SecondStderr   string
	SecondErr      error
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if req.Timeout <= 0 {
		req.Timeout = 2 * time.Minute
	}
	resp := &Response{}

	switch req.Mode {
	case "parse", "":
		s, err := debug.Parse(req.DebugEnv)
		if err != nil {
			resp.ParseErr = err.Error()
			return resp, nil
		}
		// GenPlan field lands with the feature; reflect so this tree compiles
		// pre-implement (Classic TDD RED on assert, not harness compile).
		resp.GenPlan = settingsBool(s, "GenPlan")
		resp.BypassGoTest = s.BypassGoTest
		return resp, nil

	case "cli", "cli-twice":
		r1, err := runCLI(t, req, req.Args)
		if err != nil && r1 == nil {
			return nil, err
		}
		if r1 != nil {
			resp.ExitCode = r1.ExitCode
			resp.Stdout = r1.Stdout
			resp.Stderr = r1.Stderr
			resp.Err = r1.Err
		}
		if req.Mode == "cli-twice" {
			args2 := req.Args2
			if len(args2) == 0 {
				args2 = req.Args
			}
			r2, err2 := runCLI(t, req, args2)
			if err2 != nil && r2 == nil {
				return resp, err2
			}
			if r2 != nil {
				resp.SecondExitCode = r2.ExitCode
				resp.SecondStdout = r2.Stdout
				resp.SecondStderr = r2.Stderr
				resp.SecondErr = r2.Err
			}
		}
		return resp, nil

	default:
		t.Fatalf("unknown Mode %q", req.Mode)
		return nil, nil
	}
}

// runCLI prefers in-process cli.RunWithWriters when DebugEnv is empty.
// When DebugEnv is set, FromEnv is process-initial and Parallel-safe injection
// requires a product subprocess with cmd.Env (never process os.Setenv).
func runCLI(t *testing.T, req *Request, args []string) (*Response, error) {
	t.Helper()
	resp := &Response{}
	if req.DebugEnv == "" && req.WorkDir == "" && len(req.Env) == 0 {
		var stdout, stderr bytes.Buffer
		err := cli.RunWithWriters(&stdout, &stderr, args)
		resp.Stdout = stdout.String()
		resp.Stderr = stderr.String()
		resp.Err = err
		if err != nil {
			resp.ExitCode = 1
			if resp.Stderr == "" {
				resp.Stderr = err.Error()
			}
		}
		return resp, nil
	}

	if req.Bin == "" {
		return nil, fmt.Errorf("DebugEnv/WorkDir/Env require req.Bin (root Setup Ensure)")
	}
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, req.Bin, args...)
	cmd.Dir = req.WorkDir
	base := filterEnv(os.Environ(), "DOCTEST_DEBUG")
	env := append(base, req.Env...)
	if req.DebugEnv != "" {
		env = append(env, "DOCTEST_DEBUG="+req.DebugEnv)
	}
	cmd.Env = env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	resp.Stdout = stdout.String()
	resp.Stderr = stderr.String()
	resp.Err = err
	if err == nil {
		return resp, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		resp.ExitCode = exitErr.ExitCode()
		return resp, nil
	}
	if ctx.Err() != nil {
		return resp, ctx.Err()
	}
	return resp, err
}

func filterEnv(environ []string, dropKeys ...string) []string {
	drop := map[string]bool{}
	for _, k := range dropKeys {
		drop[k] = true
	}
	out := make([]string, 0, len(environ))
	for _, e := range environ {
		k, _, _ := strings.Cut(e, "=")
		if drop[k] {
			continue
		}
		out = append(out, e)
	}
	return out
}

// settingsBool reads an optional bool field on debug.Settings by name.
// Returns false when the field is missing (pre-feature) or not a bool.
func settingsBool(s debug.Settings, field string) bool {
	v := reflect.ValueOf(s)
	f := v.FieldByName(field)
	if !f.IsValid() || f.Kind() != reflect.Bool {
		return false
	}
	return f.Bool()
}
```
