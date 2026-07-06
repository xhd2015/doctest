---
name: doc-style-test-code-specification
description: when user mentions doc-style test, with code written inside the docs
---

This specification builds on top of `doctest skill doc-spec show`.
It adds executable Go code to doc-style test documents, turning prose test cases
into runnable, verifiable tests.

The code block is just additional to the original doc, so what really matters is the doc itself, not the code.

`SETUP.md` and `ASSERT.md` must not only have code blocks, they must have meaningful description of steps to derive the code.

**NOTE: code is supplementary, the human-readable steps are precious.**

Never omit the description, and write merely code block.

## Go Code Block

Each SETUP.md or ASSERT.md may carry a **single** Go code block at the very end
of the file. Rules:

- The block is a fenced code block: `` ```go `` … `` ``` ``
- It must be the **final content** in the file — no prose, no other blocks after it
- At most **one** Go code block per file
- The block contains valid Go source (evaluated as `package testcase`)

## Request and Response Types

The root `DOCTEST.md` Go block defines a shared model inherited by all descendants:

```go
type Request struct {
    // Fields described in the ## Preconditions and ## Context sections
}

type Response struct {
    // Fields described in the ## Expected and ## Side Effects sections
}
```

Every `SETUP.md` and `ASSERT.md` in the tree refers to the same `Request` and
`Response` types and the same root-defined `Run` function. Child `SETUP.md` must
**not** redefine these types or `Run`. Root `SETUP.md` (when present) must also
**not** redefine them.

## Function Signatures

### func Setup

Declared in any `SETUP.md` along the ancestor chain. Called in order from root
to leaf before `Run`. Must implement the preconditions and steps described in
the markdown sections above the code block.

```go
func Setup(t *testing.T, req *Request) error
```

- Returns `error` — if non-nil, the test fails immediately (before `Run` and `Assert`)
- Each ancestor's `Setup` is called separately; errors report which level failed

### Helper functions

Besides `func Setup`, any `SETUP.md` Go block may declare additional functions
(helpers) for use within that directory's subtree:

- **Root `SETUP.md`** — helpers needed by all tests (e.g. building the binary under test, shared fixture writers).
- **Grouping or leaf `SETUP.md`** — helpers needed only by that node's descendants.

Helpers are inherited from ancestors and available to all descendant `Setup`,
`Run`, and `Assert` code. A child `SETUP.md` **must not** redefine a helper
with the same name as an ancestor.

### func Run

Defined **exclusively** in the root `DOCTEST.md` Go block. Executes the core
behavior under test. Must not be redefined by any descendant (Rule 9).

```go
func Run(t *testing.T, req *Request) (*Response, error)
```

- Returns `(*Response, error)` — the `error` is passed to `Assert`, not treated
  as a test failure
- Must implement exactly what the `## Steps` section describes
- If tests cannot share the same `Run(Request, Response)`, create a separate
  test tree rooted at its own `DOCTEST.md`

#### Run Is Root-Only

`func Run` is defined **exclusively in the root** `DOCTEST.md` Go block. No
descendant may redefine it. All leaves under a tree share the same
`Run(request) → response` contract.

If two test scenarios need a different `Run` function, they require separate
test trees — each rooted at its own `DOCTEST.md` with its own
`Request`/`Response`/`Run`.

#### Do Not Wrap Unit Tests

A doc-style test must execute the behavior it describes. Do **not** write a
doctest whose only behavior is running an existing unit test, such as:

```go
cmd := exec.Command("go", "test", "./pkg/foo", "-run", "TestSpecificCase")
```

That is an anti-pattern because the actual assertions live outside the
doc-style test tree, so the markdown no longer owns the specification.
Instead, migrate the relevant unit-test scenario into the doctest:

- call the target function, CLI, or API directly from `Run`
- store the observed result in `Response`
- assert the expected behavior in the leaf `ASSERT.md`
- keep any shell-out in `Run` for the actual program under test, not for
  delegating assertions to another test framework

If the function that needs direct coverage is unexported, prefer extracting or
exporting a small pure helper with a production-meaningful name over wrapping a
unit test that can access package internals. The doctest should remain the
place where the scenario is expressed and verified.

### func Assert

Declared in every `ASSERT.md`. Validates the outcomes described in
`## Expected`, `## Side Effects`, `## Errors`, and `## Exit Code`.

```go
func Assert(t *testing.T, req *Request, resp *Response, err error)
```

- `resp` is the value returned by `Run`; `err` is the error returned by `Run`
- If `Run` should succeed, check `err != nil` and fail the test
- If `Run` should fail, check that `err` matches the expected error
- All assertions go through `t.Fatal`/`t.Fatalf` — no bare `return` for failures

### Output assertions (`github.com/xhd2015/doctest/assert`)

When validating **CLI or text output**, use the output assert template DSL instead
of `strings.Contains` loops or hand-rolled `strings.Index` / `strings.Count`
parsing. Full reference: `doctest skill output-assert show`.

```go
import "github.com/xhd2015/doctest/assert"

// v2 — strict full match (preferred for new tests):
assert.Output(t, resp.Stdout, `---
version: 2
__PORT__: type=number, example=8901, a port
---
Server listen on: __PORT__
...2 lines omitted...
<ansi-color bold gray>ready</ansi-color>`)
```

v2 is strict line-by-line full match. Use `...N lines omitted...` for variable
middle sections and regex lines (e.g. `^\.+$`) for flexible single lines.

Optionally document the same template in a `## Expected Output` fenced block
(recommended for readability; not required by vet).

**Prefer DSL constructs over:**

| Legacy | v2 DSL |
|--------|--------|
| `strings.Contains` loop | strict template + `...N lines omitted...` |
| Dot/summary parsing | `^\.+$` regex line + literals |
| Dual platform `Contains` | `(linux-msg\|darwin-msg)` regex alternation |
| `metricIsColored` helpers | `<ansi-color bold gray>…</ansi-color>` |

## Validation Rules

All files are validated by `doctest build`. Rules are checked
in order; the first violation for each file is reported. Multiple violations
across different files are collected and reported together.

### 1. SETUP.md must have a Go code block

Every `SETUP.md` must contain at least one fenced Go code block (`` ```go ``).

**Violation** — a prose-only SETUP.md:
````markdown
## Preconditions
- A git repo exists

## Steps
1. Run commit
````
**Error**: `SETUP.md: must have a Go code block`

---

### 2. Go block must be the final content in the file

The Go code block must be the very last thing in the file. No markdown
prose, no text, no other blocks after the closing `` ``` ``.

**Violation** — trailing content after the block:
````markdown
```go
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { ... }
```

This sentence comes after the block and makes it non-final.
````
**Error**: `SETUP.md: go block must be final content`

---

### 3. Root DOCTEST.md must define `type Request`, `type Response`, and `func Run`

The root `DOCTEST.md` Go block (final content in the file) must declare both
shared model types **and** the `Run` function. `Run` is the tree's core behavior
contract — all leaves under this root share the same `Run`. Descendant files
reference these types and inherit `Run`; they must not redefine any of them.
Root `SETUP.md` (when present) must not define `Request`, `Response`, or `Run`.

**Violation** — DOCTEST.md defines Response and Run but not Request:
```go
type Response struct {
    Output string
}

func Run(t *testing.T, req *Request) (*Response, error) { ... }
```
**Error**: `DOCTEST.md: must define type Request and type Response`

**Violation** — DOCTEST.md defines Request and Response but not Run:
```go
type Request struct {
    Input string
}

type Response struct {
    Output string
}
```
**Error**: `DOCTEST.md: must have func Run`

**Violation** — root SETUP.md redefines shared types or Run:
```go
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { ... }
```
**Error**: `SETUP.md: Request and Response must be defined in DOCTEST.md, not SETUP.md`

---

### 4. Root DOCTEST.md must have `func Run`; every SETUP.md must have `func Setup`

The root `DOCTEST.md` Go block must define `func Run`. Every `SETUP.md` along
the ancestor chain (including optional root `SETUP.md`) must contribute a
`func Setup`. A Go block with only type declarations and no functions is
incomplete.

**Violation** — DOCTEST.md has only types, no Run:
```go
type Request struct{ Value int }
type Response struct{ Result int }
```
**Error**: `DOCTEST.md: must have func Run`

**Violation** — non-root SETUP.md has only types, no Setup:
```go
// (empty or only type/var declarations)
```
**Error**: `<dir>/SETUP.md: must have func Setup`

---

### 5. `func Setup` must have the correct signature

```go
func Setup(t *testing.T, req *Request) error
```

**Violation** — wrong parameter type:
```go
func Setup(t *testing.T, req string) error {
    return nil
}
```
**Error**: `SETUP.md: Setup must be func Setup(t *testing.T, req *Request) error`

---

### 6. `func Run` must have the correct signature

```go
func Run(t *testing.T, req *Request) (*Response, error)
```

**Violation** — missing `*Response` in return:
```go
func Run(t *testing.T, req *Request) error {
    return nil
}
```
**Error**: `SETUP.md: Run must be func Run(t *testing.T, req *Request) (*Response, error)`

---

### 7. ASSERT.md must have a Go code block with `func Assert`

Every `ASSERT.md` must contain a Go code block that defines a function
named exactly `Assert`. A block with a differently-named function does not
satisfy this rule.

**Violation** — wrong function name:
```go
func Check(t *testing.T, req *Request, resp *Response, err error) {
    t.Log("not an Assert function")
}
```
**Error**: `leaf/ASSERT.md: missing func Assert(t *testing.T, req *Request, resp *Response, err error)`

---

### 8. `func Assert` must have the correct signature

```go
func Assert(t *testing.T, req *Request, resp *Response, err error)
```

**Violation** — missing parameters:
```go
func Assert(t *testing.T) {
}
```
**Error**: `leaf/ASSERT.md: Assert must be func Assert(t *testing.T, req *Request, resp *Response, err error)`

---

### 9. Child SETUP.md cannot redefine `type Request`, `type Response`, or `func Run`

The `Request` and `Response` types and the `Run` function are defined once in
the root `DOCTEST.md` Go block. Any `SETUP.md` that redefines them creates a
conflicting model.

**Violation** — leaf redefines Request:
```go
type Request struct {
    Bad bool
}

func Setup(t *testing.T, req *Request) error { ... }
```
**Error**: `leaf/SETUP.md: child SETUP.md cannot redefine Request`

**Violation** — leaf redefines Run:
```go
func Run(t *testing.T, req *Request) (*Response, error) { ... }

func Setup(t *testing.T, req *Request) error { ... }
```
**Error**: `leaf/SETUP.md: child SETUP.md cannot redefine Run`

---

### 10. Root DOCTEST.md must define `func Run`

Every runnable leaf needs a `Run` at the root. Since `Run` can only be defined
in `DOCTEST.md` (Rule 9), the root `DOCTEST.md` Go block must provide one.

**Violation** — DOCTEST.md has only types, no Run:
- DOCTEST.md Go block: `type Request` + `type Response` (no Run)
- Leaf SETUP: `func Setup(...) error { ... }`
**Error**: `DOCTEST.md: must have func Run`

---

### 11. `func Setup` body must not be a stub

A Setup function whose body is solely `return nil` contributes nothing. Each
`func Setup` must contain real logic that implements the behavior described
in the markdown sections above it.

**Violation** — bare stub:
```go
func Setup(t *testing.T, req *Request) error {
    return nil
}
```
**Error**: `leaf/SETUP.md: func Setup body must not be a stub (return nil) — implement the behavior described in this document`

## Fixture Directory: `testdata/`

Fixtures are minimal test case trees, each designed to trigger exactly one
validation rule. They live under `testdata/` so they are skipped by the test
runner's tree discovery.

Each fixture is a directory with its own `DOCTEST.md` (root) and at least one
leaf with `SETUP.md` + `ASSERT.md`. The fixture's documents contain Go code
that violates the rule being tested.

Example — `testdata/missing-request-type/`:

```
testdata/missing-request-type/
├── DOCTEST.md         # Defines Response and Run, but NO Request type
└── leaf/
    ├── SETUP.md       # Valid leaf setup
    └── ASSERT.md      # Valid assertion
```

## Test Case Directory: `verify-*/`

Each test case is a runnable leaf under the validation tree's root. It points
`req.InputDir` at a fixture and asserts the expected compile result.

### verify-*/SETUP.md

Sets `req.InputDir` to the fixture path using a `func Setup`:

````markdown
## Steps
- Point InputDir to the <fixture-name> fixture

```go
import "path/filepath"

func Setup(t *testing.T, req *Request) error {
    req.InputDir = filepath.Join(DOCTEST_ROOT, "testdata", "<fixture-name>")
    return nil
}
```
````

### verify-*/ASSERT.md

Asserts the expected compile outcome using `func Assert`:

````markdown
## Expected Output

```
---
version: 2
---
(expected compile error text)
```

## Expected
- Compile fails
- Error message matches output template

```go
import (
    "testing"
    "github.com/xhd2015/doctest/assert"
)

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    if resp.Passed {
        t.Fatal("expected compile to fail")
    }
    assert.Output(t, resp.Output, `` +
`---
version: 2
---
(expected compile error text)`)
}
```
````

## Root SETUP.md: Test Harness

The validation tree's root `SETUP.md` defines the shared harness — the
`Request` and `Response` types plus a `Run` function that shells out to
the tool under test:

```go
import "os/exec"

type Request struct {
    InputDir string
}

type Response struct {
    Output string
    Passed bool
}

func Run(t *testing.T, req *Request) (*Response, error) {
    cmd := exec.Command("doctest", "build", req.InputDir)
    out, _ := cmd.CombinedOutput()
    return &Response{
        Output: string(out),
        Passed: cmd.ProcessState.ExitCode() == 0,
    }, nil
}
```

Each leaf overrides only `req.InputDir` in its `Setup`; the root `Run` handles
the actual invocation.

## Working Directory

Each generated test function runs with its working directory set to its own
case directory — the directory containing the leaf's `SETUP.md` and `ASSERT.md`.

The test can access `DOCTEST_ROOT` constant defined as the root of all tests, use it to refer to testdata placed at the tree root:

```go
req.InputDir = filepath.Join(DOCTEST_ROOT, "testdata", "child-redefines-request")
```

Each generated test function also defines **`DOCTEST_SESSION_ID`** — a package-level
string variable unique per `doctest test` invocation. Doctest injects it into
every generated test (via `syscall.Getenv` in generated boilerplate only). **Harness
code in `SETUP.md` / `ASSERT.md` must reference `DOCTEST_SESSION_ID` directly** —
do not call `os.Getenv("DOCTEST_SESSION_ID")` or `os.LookupEnv("DOCTEST_SESSION_ID")`.
Reading the session id through `os.Getenv` is recorded in Go's test cache key and
can pin or bust caching; `doctest vet` rejects it.

Use `DOCTEST_SESSION_ID` for session-scoped cache directories, file locks, or other
cross-package coordination within one `doctest test` run:

```go
cacheDir := filepath.Join(os.TempDir(), "my-harness-"+DOCTEST_SESSION_ID)
```

See **Session-scoped shared setup** in the design spec for the file-lock + cache
pattern used to amortize heavy setup (build binaries once, seed archives once)
across parallel leaf packages.

To reference testdata placed alongside a test case (under its own directory),
use a relative path — it resolves against the case directory:

```go
srcTestData := "./testdata"
```

## Validation

```sh
# works like go build
doctest build tests/validation
```

## Execution

Run all validation tests from the tree root:

```sh
# works like go test, run all tests
doctest test tests/validation

# run individual group of tests in a sub dir
doctest test tests/validation/group-a
doctest test tests/validation/group-a/nested-sub-group
```

This discovers every leaf (each `verify-*/` directory with an `ASSERT.md`),
generates Go test files, and runs each leaf's `Setup` → `Run` → `Assert` chain.
The tool reports which leaves passed or failed.
