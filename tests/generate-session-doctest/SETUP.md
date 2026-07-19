# Scenario

**Feature**: generated tests inject `d *session.Doctest` and stop using Chdir + free DOCTEST_* vars

```
# classic / ref / unified assemble
TreeCase -> Assemble* -> generated source with d inject

# signature rules
author Setup/Run/Assert with optional d after t -> parse/rules accept
```

## Preconditions

- Module import path is `github.com/xhd2015/doctest`.
- Package `libdoc/core` exposes `AssembleTestSource`, `AssembleRefRootSource`,
  `AssembleRefLeafTestSource`, `AssembleUnifiedLeafSource`, and parse helpers.
- Package `libdoc/rules` exposes `CheckSetupSignature` / `CheckRunSignature` /
  `CheckAssertSignature`.
- P1 type `session.Doctest` already exists at `session/doctest.go`.
- Current product still emits Chdir + free vars → these leaves are RED until P2 implementer.

## Steps

1. Leaf Setup sets `req.Op` and any author-mode / path fields.
2. Root `Run` calls the selected assemble or parse path.
3. Leaf Assert inspects `resp.Source` (and `RootSrc` for ref) or `ParseErr`.

## Context

- Fixtures are in-memory `core.TreeCase` values (no temp markdown trees for assemble leaves).
- Signature-rules leaves feed markdown snippets through `Parse*Document` and `rules.Check*`.
- Shared helpers for the inject contract live in root `DOCTEST.md` Go block.

```go
import "testing"

func Setup(t *testing.T, req *Request) error {
	// Defaults: nested leaf path, author omits d, temp abs root filled by Run.
	if req.CasePath == "" {
		req.CasePath = "nested/leaf"
	}
	if req.AuthorDMode == "" {
		req.AuthorDMode = "omit"
	}
	return nil
}

// assertInjectContract checks the NEW generated-source contract shared by all assemble paths.
func assertInjectContract(t *testing.T, label, src string) {
	t.Helper()
	if src == "" {
		t.Fatalf("%s: empty generated source", label)
	}
	if !hasSessionDoctestType(src) {
		t.Fatalf("%s: expected session.Doctest in generated source\n%s", label, src)
	}
	if !hasSessionImport(src) {
		t.Fatalf("%s: expected import github.com/xhd2015/doctest/session\n%s", label, src)
	}
	if !hasDConstruct(src) {
		t.Fatalf("%s: expected d := &session.Doctest{...} construct\n%s", label, src)
	}
	if hasLeafChdirBoilerplate(src) {
		t.Fatalf("%s: generated source must not contain leaf os.Chdir / Getwd restore boilerplate\n%s", label, src)
	}
	if hasPackageFreeDoctestVars(src) {
		t.Fatalf("%s: generated source must not declare/assign package free DOCTEST_ROOT / DOCTEST_SESSION_ID\n%s", label, src)
	}
	if !passesDToSetupRunAssert(src) {
		t.Fatalf("%s: Setup/Run/Assert call sites must pass d\n%s", label, src)
	}
}
```
