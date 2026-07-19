# Doctest Structure Layout Tests

These tests specify the restructured root doctest tree layout: spec version in
`DOCTEST.md`, `Request`/`Response`/`Run` relocated from root `SETUP.md` to
`DOCTEST.md`, updated `doctest vet` rules, skill prompt version injection, and
build/test integration with the new layout.

## Version

0.0.2

## DSN (Domain Specific Notion)

### Participants

- **Doctest CLI** — the binary under test; exposes `vet`, `build`, `test`, and
  `skill` subcommands.
- **Test tree** — a directory of `DOCTEST.md`, optional root `SETUP.md`, and
  descendant `SETUP.md`/`ASSERT.md` leaves.
- **Vet validator** — walks a tree and reports structural violations (missing
  version section, misplaced `Request`/`Response`/`Run`, anti-patterns).
- **Skill embedder** — serves embedded markdown prompts via `doctest skill
  <name> show`, resolving `__DOCTEST_VERSION__` to the canonical version.
- **Tree assembler** — reads `Request`/`Response`/`Run` from root
  `DOCTEST.md` and compiles runnable leaves.

### Behaviors

- **`vet`** — checks `## Version` presence, `Request`/`Response`/`Run` in
  `DOCTEST.md` Go block, rejects them in root `SETUP.md`, optional root
  `SETUP.md`.
- **`skill --show`** — prints embedded docs with literal version string, no
  placeholder residue.
- **`build` / `test`** — discover and execute trees using the new root layout.

## Decision Tree

```
doctest-structure/                              [operation mode]
│
├── vet/                                        Target: structural validation
│   ├── version/                                Param: ## Version section
│   │   ├── missing                             → no ## Version → vet error
│   │   └── present                             → ## Version present → vet pass
│   ├── run-in-doctest/                         Param: Request/Response/Run in DOCTEST.md
│   │   ├── valid-with-root-setup               → types in DOCTEST.md, root SETUP has Setup only → pass
│   │   ├── valid-no-root-setup                 → types in DOCTEST.md, no root SETUP.md → pass
│   │   ├── missing-run                         → DOCTEST.md Go block lacks func Run → error
│   │   └── missing-types                       → DOCTEST.md Go block lacks Request/Response → error
│   └── run-in-setup/                           Param: legacy placement
│       └── rejected                            → Request/Response/Run in root SETUP.md → error
│
├── skill-show/                                 Target: version injection in prompts
│   ├── skill-tdd-show                          → doctest skill tdd --show
│   ├── skill-tdd-lite-show                     → doctest skill tdd-lite --show
│   ├── skill-designer-show                     → doctest skill designer --show
│   └── skill-implementer-show                  → doctest skill implementer --show
│
└── integration/                                Target: build/test with new layout
    ├── build-runs                              → minimal valid tree → doctest build succeeds
    └── test-runs                               → minimal valid tree → doctest test succeeds
```

## Test Index

| Leaf | Description |
|------|-------------|
| `vet/version/missing` | Vet fails when `DOCTEST.md` lacks `## Version` |
| `vet/version/present` | Vet passes when `## Version` is present (value not validated) |
| `vet/run-in-doctest/valid-with-root-setup` | Vet passes: types in `DOCTEST.md`, root `SETUP.md` has only `Setup` |
| `vet/run-in-doctest/valid-no-root-setup` | Vet passes: types in `DOCTEST.md`, no root `SETUP.md` |
| `vet/run-in-doctest/missing-run` | Vet fails when `DOCTEST.md` Go block lacks `func Run` |
| `vet/run-in-doctest/missing-types` | Vet fails when `DOCTEST.md` Go block lacks `Request` or `Response` |
| `vet/run-in-setup/rejected` | Vet fails when `Request`/`Response`/`Run` remain in root `SETUP.md` |
| `skill-show/skill-tdd-show` | `doctest skill tdd --show` contains `0.0.2`, no placeholder |
| `skill-show/skill-tdd-lite-show` | `doctest skill tdd-lite --show` contains `0.0.2`, no placeholder |
| `skill-show/skill-designer-show` | `doctest skill designer --show` contains `0.0.2`, no placeholder |
| `skill-show/skill-implementer-show` | `doctest skill implementer --show` contains `0.0.2`, no placeholder |
| `integration/build-runs` | `doctest build` succeeds on minimal new-layout tree |
| `integration/test-runs` | `doctest test` succeeds on minimal new-layout tree |

## How to Run

```sh
doctest vet ./tests/doctest-structure
doctest test ./tests/doctest-structure/...
doctest test ./tests/doctest-structure/vet/...
doctest test ./tests/doctest-structure/skill-show/...
doctest test ./tests/doctest-structure/integration/...
```

```go
import (
"github.com/xhd2015/doctest/session"
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/testbin"
)

const canonicalVersion = "0.0.2"

type Request struct {
	Args    []string
	Env     []string
	WorkDir string
	Timeout time.Duration
	Bin     string
}

type Response struct {
	ExitCode int
	Stdout   string
	Stderr   string
	Err      error
}

func Run(t *testing.T, req *Request) (*Response, error) {
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()

	bin := req.Bin
	if bin == "" {
		return nil, fmt.Errorf("req.Bin is not set")
	}
	cmd := exec.CommandContext(ctx, bin, req.Args...)
	cmd.Dir = req.WorkDir
	cmd.Env = append(os.Environ(), req.Env...)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	resp := &Response{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
		Err:    err,
	}
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

func buildDoctestBin(t *testing.T, d *session.Doctest) string {
	t.Helper()
	return testbin.Ensure(t, filepath.Join(d.DOCTEST_ROOT, "..", ".."))
}

var bt = string([]byte{96, 96, 96})

func goBlock(code string) string {
	return bt + "go\n" + code + bt
}

const minimalDSN = "## DSN (Domain Specific Notion)\n\n### Participants\n- **system** — under test.\n\n### Behaviors\n- **run** — executes the scenario.\n"

func validTypesGoBlock() string {
	return goBlock(
		"import \"testing\"\n\n" +
			"type Request struct{}\n" +
			"type Response struct{}\n\n" +
			"func Run(t *testing.T, req *Request) (*Response, error) {\n" +
			"\treturn &Response{}, nil\n" +
			"}",
	)
}

func typesWithoutRunGoBlock() string {
	return goBlock(
		"import \"testing\"\n\n" +
			"type Request struct{}\n" +
			"type Response struct{}\n",
	)
}

func typesWithoutRequestGoBlock() string {
	return goBlock(
		"import \"testing\"\n\n" +
			"type Response struct{}\n\n" +
			"func Run(t *testing.T, req *Request) (*Response, error) {\n" +
			"\treturn &Response{}, nil\n" +
			"}",
	)
}

type treeOpts struct {
	withVersion     bool
	version         string
	doctestGoBlock  string
	withRootSetup   bool
	rootSetupBody   string
	withLeaf        bool
}

func writeDoctestMD(dir string, opts treeOpts) error {
	var b strings.Builder
	b.WriteString("# Tests\n\n")
	if opts.withVersion {
		v := opts.version
		if v == "" {
			v = canonicalVersion
		}
		b.WriteString("## Version\n")
		b.WriteString(v)
		b.WriteString("\n\n")
	}
	b.WriteString(minimalDSN)
	b.WriteString("\n")
	block := opts.doctestGoBlock
	if block == "" {
		block = validTypesGoBlock()
	}
	b.WriteString(block)
	b.WriteString("\n")
	return os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte(b.String()), 0644)
}

func writeRootSetupOnly(dir string) error {
	content := "# Scenario\n\n**Feature**: root setup without types\n\n" +
		bt + "\n# root setup only\nroot -> Setup\n" + bt + "\n\n" +
		"## Steps\n1. noop\n\n" +
		goBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { return nil }") +
		"\n"
	return os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(content), 0644)
}

func writeRootSetupWithTypes(dir string) error {
	content := "# Scenario\n\n**Feature**: legacy root setup with types\n\n" +
		bt + "\n# legacy placement\nroot SETUP.md -> Request/Response/Run\n" + bt + "\n\n" +
		"## Steps\n1. noop\n\n" +
		goBlock(
			"import \"testing\"\n\n" +
				"type Request struct{}\n" +
				"type Response struct{}\n\n" +
				"func Run(t *testing.T, req *Request) (*Response, error) {\n" +
				"\treturn &Response{}, nil\n" +
				"}",
		) + "\n"
	return os.WriteFile(filepath.Join(dir, "SETUP.md"), []byte(content), 0644)
}

func writeMinimalLeaf(dir string) error {
	leafDir := filepath.Join(dir, "leaf")
	if err := os.MkdirAll(leafDir, 0755); err != nil {
		return err
	}
	leafSetup := "# Scenario\n\n**Feature**: minimal runnable leaf\n\n" +
		bt + "\n# leaf runs\nleaf -> Run -> pass\n" + bt + "\n\n" +
		"## Steps\n1. noop\n\n" +
		goBlock("import \"testing\"\n\nfunc Setup(t *testing.T, req *Request) error { return nil }") +
		"\n"
	leafAssert := "## Expected\n- Run succeeds without error.\n\n" +
		goBlock(
			"import \"testing\"\n\n" +
				"func Assert(t *testing.T, req *Request, resp *Response, err error) {\n" +
				"\tif err != nil {\n" +
				"\t\tt.Fatal(err)\n" +
				"\t}\n" +
				"}",
		) + "\n"
	if err := os.WriteFile(filepath.Join(leafDir, "SETUP.md"), []byte(leafSetup), 0644); err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(leafDir, "ASSERT.md"), []byte(leafAssert), 0644)
}

func writeTree(t *testing.T, opts treeOpts) string {
	t.Helper()
	dir := t.TempDir()
	if err := writeDoctestMD(dir, opts); err != nil {
		t.Fatal(err)
	}
	if opts.withRootSetup {
		body := opts.rootSetupBody
		if body == "types" {
			if err := writeRootSetupWithTypes(dir); err != nil {
				t.Fatal(err)
			}
		} else {
			if err := writeRootSetupOnly(dir); err != nil {
				t.Fatal(err)
			}
		}
	}
	if opts.withLeaf {
		if err := writeMinimalLeaf(dir); err != nil {
			t.Fatal(err)
		}
	}
	return dir
}

func setVetArgs(t *testing.T, req *Request, treeDir string) {
	t.Helper()
	req.Args = []string{"vet", treeDir}
}

func assertVetPass(t *testing.T, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected vet pass (exit 0), got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}

func assertVetFail(t *testing.T, resp *Response, err error, stderrSubstr string) {
	t.Helper()
	if err != nil {
		t.Fatalf("run failed unexpectedly: %v", err)
	}
	if resp.ExitCode == 0 {
		t.Fatalf("expected vet fail (nonzero exit), got 0\nstdout:\n%s", resp.Stdout)
	}
	combined := strings.ToLower(resp.Stdout + resp.Stderr)
	if !strings.Contains(combined, strings.ToLower(stderrSubstr)) {
		t.Fatalf("expected stderr/stdout to mention %q, got:\nstdout:\n%s\nstderr:\n%s", stderrSubstr, resp.Stdout, resp.Stderr)
	}
}

func assertSkillShowVersion(t *testing.T, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("exit code = %d, stderr:\n%s", resp.ExitCode, resp.Stderr)
	}
	if !strings.Contains(resp.Stdout, canonicalVersion) {
		t.Fatalf("stdout missing version %q:\n%s", canonicalVersion, resp.Stdout)
	}
	if strings.Contains(resp.Stdout, "__DOCTEST_VERSION__") {
		t.Fatalf("stdout still contains unresolved placeholder __DOCTEST_VERSION__")
	}
}

func assertCommandPass(t *testing.T, resp *Response, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("run failed: %v", err)
	}
	if resp.ExitCode != 0 {
		t.Fatalf("expected exit 0, got %d\nstdout:\n%s\nstderr:\n%s", resp.ExitCode, resp.Stdout, resp.Stderr)
	}
}
```