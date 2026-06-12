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

The root `SETUP.md` defines a shared model inherited by all descendants:

```go
type Request struct {
    // Fields described in the ## Preconditions and ## Context sections
}

type Response struct {
    // Fields described in the ## Expected and ## Side Effects sections
}
```

Every `SETUP.md` and `ASSERT.md` in the tree refers to the same `Request` and
`Response` types. Child `SETUP.md` must **not** redefine these types.

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

### func Run

Defined in at most one `SETUP.md` in the root-to-leaf chain. The **deepest**
`Run` in the ancestor path is used. Executes the core behavior under test.

```go
func Run(t *testing.T, req *Request) (*Response, error)
```

- Returns `(*Response, error)` — the `error` is passed to `Assert`, not treated
  as a test failure
- Must implement exactly what the `## Steps` section describes
- At least one `Run` must exist in the entire setup chain

#### Run Override: Nearest Wins

Only the **deepest** (nearest-to-leaf) `Run` in the ancestor chain is executed.
Any node along the path — including the leaf itself — can define its own `Run`
to override the inherited one.

```
Root SETUP:        func Run() → returns error "not implemented yet"
  mode-commit/:    (no Run — inherits root's Run)
    leaf/:         func Run() → implements commit for this specific scenario
```

The leaf uses its own `Run`, ignoring all ancestors. This allows a leaf to
specialize behavior while still inheriting Setup from the chain above it.

The convention is for the root to provide a default `Run` (e.g., a harness that
shells out to the tool under test). Children and leaves override it only when
they need different behavior.

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

### 3. Root SETUP.md must define `type Request` and `type Response`

The root `SETUP.md` (the one at the tree root, not in a subdirectory) must
declare both shared model types. All descendant files reference these types.

**Violation** — root defines Response but not Request:
```go
type Response struct {
    Output string
}

func Run(t *testing.T, req *Request) (*Response, error) { ... }
```
**Error**: `SETUP.md: must define type Request and type Response`

---

### 4. Every SETUP.md must have `func Setup` or `func Run`

Every `SETUP.md` must contribute at least one executable function. A Go block
with only type declarations and no functions is incomplete.

**Violation** — types only, no functions:
```go
type Request struct{ Value int }
type Response struct{ Result int }
```
**Error**: `SETUP.md: must have func Setup or func Run`

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

### 9. Child SETUP.md cannot redefine `type Request` or `type Response`

The `Request` and `Response` types are defined once at the root. Any child
`SETUP.md` that redefines them creates a conflicting model.

**Violation** — leaf redefines Request:
```go
type Request struct {
    Bad bool
}

func Setup(t *testing.T, req *Request) error { ... }
```
**Error**: `leaf/SETUP.md: child SETUP.md cannot redefine Request`

---

### 10. At least one SETUP.md in the chain must define `func Run`

Every runnable leaf needs a `Run` somewhere in its ancestor chain. Without
one, there is no executable behavior to test.

**Violation** — entire chain has only Setup, no Run anywhere:
- Root SETUP: `func Setup(...) error { ... }` (no Run)
- Leaf SETUP: `func Setup(...) error { ... }` (no Run)
**Error**: `leaf/ASSERT.md: no Run(t *testing.T, req *Request) (*Response, error) in setup chain`

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

Each fixture is a directory with its own `SETUP.md` (root) and at least one
leaf with `SETUP.md` + `ASSERT.md`. The fixture's documents contain Go code
that violates the rule being tested.

Example — `testdata/missing-request-type/`:

```
testdata/missing-request-type/
├── SETUP.md           # Defines Response and Run, but NO Request type
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
## Expected
- Compile fails
- Error message contains "<expected error text>"

```go
import "strings"

func Assert(t *testing.T, req *Request, resp *Response, err error) {
    if err != nil {
        t.Fatal(err)
    }
    if resp.Passed {
        t.Fatal("expected compile to fail")
    }
    if !strings.Contains(resp.Output, "<expected error text>") {
        t.Fatalf("expected '<expected error text>' in output, got:\n%s", resp.Output)
    }
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
