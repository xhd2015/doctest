package build

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/testtree"
)

func TestMultiCaseRunIgnoresStaleNestedDoctestPackages(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "parent/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	nested := filepath.Join(root, "nested")
	writeTreeFile(t, nested, "README.md", "# nested")
	writeRootHarness(t, nested, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, nested, "bad/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, nested, "bad/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Fatal("nested doctest should not run when testing parent tree")
}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Test(nested, core.Options{GenDir: genDir}); err == nil {
		t.Fatal("expected nested doctest to fail")
	}
	if err := Test(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("parent doctest should ignore stale nested generated packages: %v", err)
	}
}

func TestChildCannotRedefineRun(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{ Source string }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Source: "root"}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Source: "leaf"}, nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected child Run redefine error")
	}
	if !strings.Contains(err.Error(), "cannot redefine Run") {
		t.Fatalf("expected cannot redefine Run, got %v", err)
	}
}

func TestExecutionOrderSetupRunAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{ Order []string }
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { req.Order = append(req.Order, "run"); return &Response{}, nil }
`, `
func Setup(t *testing.T, req *Request) error { req.Order = append(req.Order, "root setup"); return nil }
`)
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Order = append(req.Order, "parent setup"); return nil }
`))
	writeTreeFile(t, root, "parent/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Order = append(req.Order, "leaf setup"); return nil }
`))
	writeTreeFile(t, root, "parent/leaf/ASSERT.md", assertDoc(`
import "reflect"
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	req.Order = append(req.Order, "assert")
	want := []string{"root setup", "parent setup", "leaf setup", "run", "assert"}
	if !reflect.DeepEqual(req.Order, want) { t.Fatalf("order = %#v, want %#v", req.Order, want) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestSetupErrorFailsBeforeRunAndAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
import "fmt"
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { t.Fatal("run should not execute"); return nil, nil }
`, `
func Setup(t *testing.T, req *Request) error { return fmt.Errorf("setup failed") }
`)
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	t.Fatal("assert should not execute")
}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected setup failure")
	}
	if !strings.Contains(err.Error(), "setup failed") && !strings.Contains(err.Error(), "exit status 1") {
		t.Fatalf("expected setup failure in error, got %v", err)
	}
}

func TestRunErrorPassedToAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
import "fmt"
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return nil, fmt.Errorf("run failed") }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err == nil || err.Error() != "run failed" { t.Fatalf("expected run failed error, got %v", err) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestRequestMutatedThroughSetupChain(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{ Value int }
type Response struct{ Value int }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Value: req.Value}, nil }
`, `
func Setup(t *testing.T, req *Request) error { req.Value += 1; return nil }
`)
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Value += 2; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Value += 3; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if req.Value != 6 || resp.Value != 6 { t.Fatalf("req=%d resp=%d, want 6", req.Value, resp.Value) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestResponsePassedFromRunToAssert(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{ Message string }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Message: "ok"}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if resp == nil || resp.Message != "ok" { t.Fatalf("unexpected response: %#v", resp) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestDuplicateSetupHooksAcrossLevelsAllowed(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{ Count int }
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, `
func Setup(t *testing.T, req *Request) error { req.Count++; return nil }
`)
	writeTreeFile(t, root, "parent/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Count++; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { req.Count++; return nil }
`))
	writeTreeFile(t, root, "parent/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if req.Count != 3 { t.Fatalf("count = %d, want 3", req.Count) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestDuplicateNamesAcrossDifferentLeavesAllowed(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{ Name string }
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{Name: "root"}, nil }
`, "")
	writeTreeFile(t, root, "a/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "a/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	var duplicate = "a"
	if duplicate != "a" { t.Fatal(duplicate) }
}
`))
	writeTreeFile(t, root, "b/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "b/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	var duplicate = "b"
	if duplicate != "b" { t.Fatal(duplicate) }
}
`))

	if err := Test(root, core.Options{RemoveTemp: true}); err != nil {
		t.Fatalf("test: %v", err)
	}
}

func TestChildCannotRedefineRequestOrResponse(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
type Request struct{ Bad bool }
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected child Request redefinition error")
	}
	if !strings.Contains(err.Error(), "Request") {
		t.Fatalf("expected Request conflict error, got %v", err)
	}
}

func TestGeneratedCodeHasDoctestRootConst(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	// Default hierarchical unified: leaf is non-test RunTestLeaf package.
	leafData, err := os.ReadFile(filepath.Join(genDir, "leaf", "leaf.go"))
	if err != nil {
		t.Fatalf("read leaf.go: %v", err)
	}
	code := string(leafData)
	absRoot, _ := filepath.Abs(root)
	assertGeneratedMatchesFixture(t, code, absRoot, "generated_leaf.go.fixture")
	if !strings.Contains(code, "func init()") {
		t.Fatal("expected func init() registration in unified leaf package")
	}
	if !strings.Contains(code, "session.Doctest") {
		t.Fatal("expected session.Doctest inject in unified leaf")
	}
	if !strings.Contains(code, "DOCTEST_ROOT") {
		t.Fatal("expected DOCTEST_ROOT in unified leaf")
	}
}

func TestGeneratedRootCaseChdirsToDoctestRoot(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	// Root-level ASSERT leaf: gen path is genRoot itself under unified layout.
	// Prefer suite package as the runnable entry; leaf package is non-test.
	suitePath := filepath.Join(genDir, "suite", "suite_test.go")
	if _, err := os.Stat(suitePath); err != nil {
		t.Fatalf("expected suite/suite_test.go: %v", err)
	}
	// Root leaf package file (non-test) next to gen root modules or named by leaf.
	var leafCode string
	entries, _ := os.ReadDir(genDir)
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		if strings.HasSuffix(e.Name(), ".go") && !strings.HasSuffix(e.Name(), "_test.go") {
			data, err := os.ReadFile(filepath.Join(genDir, e.Name()))
			if err != nil {
				t.Fatalf("read %s: %v", e.Name(), err)
			}
			leafCode = string(data)
			break
		}
	}
	if leafCode == "" {
		// Tree-scoped leaf may sit under a dedicated leaf package dir at ".".
		// Walk for RunTestLeaf.
		_ = filepath.Walk(genDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info == nil || info.IsDir() {
				return nil
			}
			if !strings.HasSuffix(path, ".go") || strings.HasSuffix(path, "_test.go") {
				return nil
			}
			data, readErr := os.ReadFile(path)
			if readErr != nil {
				return nil
			}
			if strings.Contains(string(data), "func RunTestLeaf") {
				leafCode = string(data)
				return filepath.SkipAll
			}
			return nil
		})
	}
	if leafCode == "" {
		t.Fatal("expected unified root leaf RunTestLeaf package under gen")
	}
	absRoot, _ := filepath.Abs(root)
	assertGeneratedMatchesFixture(t, leafCode, absRoot, "generated_root_leaf.go.fixture")
	if !strings.Contains(leafCode, "session.Doctest") {
		t.Fatal("expected session.Doctest inject in root leaf")
	}
}

func TestCompileGoModGeneratedWithModule(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/test\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "tests/README.md", "# tree")
	writeRootHarness(t, filepath.Join(srcRoot, "tests"), `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, srcRoot, "tests/leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, srcRoot, "tests/leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(filepath.Join(srcRoot, "tests"), core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod in gen dir: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got %q", goMod)
	}
	if !strings.Contains(goMod, "example.com/test") {
		t.Fatalf("expected require example.com/test, got %q", goMod)
	}
	if !strings.Contains(goMod, "replace example.com/test => "+srcRoot) {
		t.Fatalf("expected replace with abs path %q, got %q", srcRoot, goMod)
	}
}

func TestCompileNoGoModWhenNoSourceModule(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod in gen dir: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "module testcase") {
		t.Fatalf("expected module testcase, got %q", goMod)
	}
	// No source-module replace when the tree has no go.mod. Harness inject still
	// adds assert-mod + session-mod replaces (always-on for shared gen-root stability).
	if !strings.Contains(goMod, "replace github.com/xhd2015/doctest/session =>") {
		t.Fatalf("expected session replace, got %q", goMod)
	}
	if !strings.Contains(goMod, "replace github.com/xhd2015/doctest/assert =>") {
		t.Fatalf("expected assert replace, got %q", goMod)
	}
	for _, line := range strings.Split(goMod, "\n") {
		line = strings.TrimSpace(line)
		if !strings.HasPrefix(line, "replace ") {
			continue
		}
		if strings.Contains(line, "github.com/xhd2015/doctest/session") ||
			strings.Contains(line, "github.com/xhd2015/doctest/assert") {
			continue
		}
		t.Fatalf("expected only assert/session replaces when no source module, got %q", goMod)
	}
}

func TestCompileWithExternalPackageImport(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/mylib\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "pkg/helper.go", "package pkg\n\nfunc Name() string { return \"hello\" }\n")

	testsDir := filepath.Join(srcRoot, "tests")
	writeTreeFile(t, testsDir, "README.md", "# tree")
	writeTreeFile(t, testsDir, "DOCTEST.md", doctestDoc(`
import "example.com/mylib/pkg"

type Request struct{}
type Response struct{ Val string }
func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Val: pkg.Name()}, nil
}
`))
	writeTreeFile(t, testsDir, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, testsDir, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if resp.Val != "hello" { t.Fatalf("expected hello, got %q", resp.Val) }
}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(testsDir, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	goModData, err := os.ReadFile(filepath.Join(genDir, "go.mod"))
	if err != nil {
		t.Fatalf("expected go.mod in gen dir: %v", err)
	}
	goMod := string(goModData)
	if !strings.Contains(goMod, "replace example.com/mylib => "+srcRoot) {
		t.Fatalf("expected replace with path %q, got %q", srcRoot, goMod)
	}

	cmd := exec.Command("go", "test", "-count=1", "./...")
	cmd.Dir = genDir
	cmd.Env = append(os.Environ(), core.DoctestSessionIDEnv+"=test-session")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("go test in gen dir failed: %v\n%s", err, out)
	}
}

func TestSourceFilesCopiedWithPackageRename(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/mylib\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "pkg/helper.go", "package pkg\n\nfunc Exported() string { return \"exported\" }\nfunc private() string { return \"private\" }\n")
	writeTreeFile(t, srcRoot, "pkg/types.go", "package pkg\n\ntype MyType struct{ Val string }\n")

	testsDir := filepath.Join(srcRoot, "tests")
	writeTreeFile(t, testsDir, "README.md", "# tree")
	writeTreeFile(t, testsDir, "DOCTEST.md", doctestDoc(`import "example.com/mylib/pkg"

type Request struct{}
type Response struct{ Msg string }
func Run(t *testing.T, req *Request) (*Response, error) {
	return &Response{Msg: Exported() + " " + private()}, nil
}`))
	writeTreeFile(t, testsDir, "SETUP.md", `# Setup
- Go module: example.com/mylib
- Package under test: pkg

`+setupDoc(`func Setup(t *testing.T, req *Request) error { _ = req; return nil }`))

	writeTreeFile(t, testsDir, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, testsDir, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {
	if err != nil { t.Fatal(err) }
	if resp.Msg != "exported private" { t.Fatalf("got %q", resp.Msg) }
}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(testsDir, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	helperPath := filepath.Join(genDir, "tests", "leaf", "helper.go")
	helperData, err := os.ReadFile(helperPath)
	if err != nil {
		t.Fatalf("expected helper.go in gen dir: %v", err)
	}
	if !strings.Contains(string(helperData), "package pkg_tc") {
		t.Fatalf("expected package pkg_tc in helper.go, got:\n%s", helperData)
	}
	if strings.Contains(string(helperData), "package pkg\n") {
		t.Fatal("expected original package name to be replaced")
	}
	if !strings.Contains(string(helperData), "func Exported()") {
		t.Fatal("expected Exported func in copied helper.go")
	}
	if !strings.Contains(string(helperData), "func private()") {
		t.Fatal("expected private func in copied helper.go")
	}

	// Unified leaf is non-test; package under test is also copied into __droot
	// so Run can call unexported symbols.
	leafPath := filepath.Join(genDir, "tests", "leaf", "leaf.go")
	testData, err := os.ReadFile(leafPath)
	if err != nil {
		t.Fatalf("expected leaf.go: %v", err)
	}
	if !strings.Contains(string(testData), "package pkg_tc") {
		t.Fatalf("expected package pkg_tc in leaf package, got:\n%s", testData)
	}
	drootHelper := filepath.Join(genDir, "tests", "__droot", "helper.go")
	if _, err := os.Stat(drootHelper); err != nil {
		t.Fatalf("expected package-under-test sources in __droot: %v", err)
	}
}

func TestNoTestGoFilesCopied(t *testing.T) {
	srcRoot := t.TempDir()
	writeTreeFile(t, srcRoot, "go.mod", "module example.com/mylib\n\ngo 1.21\n")
	writeTreeFile(t, srcRoot, "pkg/helper.go", "package pkg\n\nfunc Ok() string { return \"ok\" }\n")
	writeTreeFile(t, srcRoot, "pkg/helper_test.go", "package pkg\n\nimport \"testing\"\n\nfunc TestFoo(t *testing.T) {}\n")

	testsDir := filepath.Join(srcRoot, "tests")
	writeTreeFile(t, testsDir, "README.md", "# tree")
	writeTreeFile(t, testsDir, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`))
	writeTreeFile(t, testsDir, "SETUP.md", `# Setup
- Go module: example.com/mylib
- Package under test: pkg

`+setupDoc(`func Setup(t *testing.T, req *Request) error { _ = req; return nil }`))

	writeTreeFile(t, testsDir, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, testsDir, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(testsDir, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	if _, err := os.Stat(filepath.Join(genDir, "tests", "leaf", "helper.go")); os.IsNotExist(err) {
		t.Fatal("expected helper.go to be copied")
	}
	if _, err := os.Stat(filepath.Join(genDir, "tests", "leaf", "helper_test.go")); !os.IsNotExist(err) {
		t.Fatal("expected helper_test.go NOT to be copied")
	}
}

func TestBackwardCompatNoPackageUnderTest(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir}); err != nil {
		t.Fatalf("build: %v", err)
	}

	leafData, err := os.ReadFile(filepath.Join(genDir, "leaf", "leaf.go"))
	if err != nil {
		t.Fatalf("read leaf.go: %v", err)
	}
	if !strings.Contains(string(leafData), "package testcase") {
		t.Fatalf("expected package testcase when no metadata, got:\n%s", leafData)
	}
}

func TestCompileDefaultKeepsTempDir(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	var stderr bytes.Buffer
	if err := Build(root, core.Options{Stderr: &stderr}); err != nil {
		t.Fatalf("build: %v", err)
	}

	out := stderr.String()
	i := strings.Index(out, "\u2192 ")
	if i < 0 {
		t.Fatalf("expected gen dir in output, got %q", out)
	}
	rest := out[i+len("\u2192 "):]
	genDir := strings.TrimSpace(strings.SplitN(rest, "\n", 2)[0])

	if _, err := os.Stat(genDir); os.IsNotExist(err) {
		t.Fatalf("expected temp dir %s to be kept, but it was removed", genDir)
	}
	defer os.RemoveAll(genDir)
}

func TestCompileRmRemovesTempDir(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	var stderr bytes.Buffer
	if err := Build(root, core.Options{RemoveTemp: true, Stderr: &stderr}); err != nil {
		t.Fatalf("build: %v", err)
	}

	out := stderr.String()
	i := strings.Index(out, "\u2192 ")
	if i < 0 {
		t.Fatalf("expected gen dir in output, got %q", out)
	}
	rest := out[i+len("\u2192 "):]
	genDir := strings.TrimSpace(strings.SplitN(rest, "\n", 2)[0])

	if _, err := os.Stat(genDir); !os.IsNotExist(err) {
		t.Fatalf("expected temp dir %s to be removed, but it still exists", genDir)
	}
}

func TestCompileRmDoesNotRemoveGenDir(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeRootHarness(t, root, `
type Request struct{}
type Response struct{}
func Run(t *testing.T, req *Request) (*Response, error) { return &Response{}, nil }
`, "")
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	genDir := filepath.Join(t.TempDir(), "generated")
	if err := Build(root, core.Options{GenDir: genDir, RemoveTemp: true}); err != nil {
		t.Fatalf("build: %v", err)
	}

	if _, err := os.Stat(genDir); os.IsNotExist(err) {
		t.Fatalf("expected gen dir %s to exist with --rm, but it was removed", genDir)
	}
}

func TestAtLeastOneRunRequiredInSetupChain(t *testing.T) {
	root := t.TempDir()
	writeTreeFile(t, root, "README.md", "# tree")
	writeTreeFile(t, root, "DOCTEST.md", doctestDoc(`
type Request struct{}
type Response struct{}
`))
	writeTreeFile(t, root, "SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/SETUP.md", setupDoc(`
func Setup(t *testing.T, req *Request) error { _ = req; return nil }
`))
	writeTreeFile(t, root, "leaf/ASSERT.md", assertDoc(`
func Assert(t *testing.T, req *Request, resp *Response, err error) {}
`))

	err := Test(root, core.Options{RemoveTemp: true})
	if err == nil {
		t.Fatal("expected missing Run error")
	}
	if !strings.Contains(err.Error(), "Run") {
		t.Fatalf("expected missing Run error, got %v", err)
	}
}

func TestDotProgressIncremental(t *testing.T) {
	if testing.Short() {
		t.Skip("slow dot progress timing test")
	}

	subRoot := t.TempDir()
	testtree.WriteMinimalRunnableTree(t, subRoot, []testtree.LeafSpec{
		{Name: "a_fast", Steps: "No setup needed.", Expected: "Always passes."},
		{Name: "z_slow", Steps: "Sleep 2 seconds to simulate a long-running test.", Expected: "Always passes.",
			SetupGo: "import (\"testing\"; \"time\")\n\nfunc Setup(t *testing.T, req *Request) error { time.Sleep(2 * time.Second); return nil }"},
	})

	oldStdout := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	os.Stdout = w

	type dotInfo struct {
		firstDot time.Duration
		output   string
	}
	ch := make(chan dotInfo, 1)
	start := time.Now()
	go func() {
		var buf bytes.Buffer
		firstDot := time.Duration(-1)
		tmp := make([]byte, 1)
		for {
			n, readErr := r.Read(tmp)
			if n > 0 {
				buf.WriteByte(tmp[0])
				if tmp[0] == '.' && firstDot < 0 {
					firstDot = time.Since(start)
				}
			}
			if readErr != nil {
				break
			}
		}
		ch <- dotInfo{firstDot, buf.String()}
	}()

	if err := Test(subRoot, core.Options{RemoveTemp: true, Count: 1}); err != nil {
		t.Fatalf("build.Test: %v", err)
	}
	w.Close()
	info := <-ch
	os.Stdout = oldStdout

	totalElapsed := time.Since(start)
	incremental := info.firstDot >= 0 && (totalElapsed-info.firstDot) > 800*time.Millisecond
	if !incremental {
		t.Fatal("dots are NOT printed incrementally — the first dot appeared after the slow package finished")
	}

	inlineIdx := strings.Index(info.output, "  (")
	if inlineIdx < 0 {
		t.Fatalf("expected summary line in output:\n%s", info.output)
	}
	if dots := strings.Count(info.output[:inlineIdx], "."); dots != 2 {
		t.Fatalf("expected 2 dots before summary, got %d. output:\n%s", dots, info.output)
	}
}

// assertGeneratedMatchesFixture compares generated Go source to a golden
// fixture under testdata/. Absolute doctest roots are normalized to
// {{DOCTEST_ROOT}} so the checked-in file is path-stable and reviewable.
//
// Set UPDATE_FIXTURES=1 to rewrite the fixture from the current generator.
func assertGeneratedMatchesFixture(t *testing.T, got, absRoot, fixtureName string) {
	t.Helper()
	normalized := strings.ReplaceAll(got, absRoot, "{{DOCTEST_ROOT}}")
	// Stable trailing newline for diffs.
	if !strings.HasSuffix(normalized, "\n") {
		normalized += "\n"
	}

	fixturePath := filepath.Join("testdata", fixtureName)
	if os.Getenv("UPDATE_FIXTURES") == "1" {
		if err := os.MkdirAll(filepath.Dir(fixturePath), 0o755); err != nil {
			t.Fatalf("mkdir testdata: %v", err)
		}
		if err := os.WriteFile(fixturePath, []byte(normalized), 0o644); err != nil {
			t.Fatalf("write fixture %s: %v", fixturePath, err)
		}
		t.Logf("updated fixture %s", fixturePath)
		return
	}

	wantBytes, err := os.ReadFile(fixturePath)
	if err != nil {
		t.Fatalf("read fixture %s: %v (run with UPDATE_FIXTURES=1 to create)", fixturePath, err)
	}
	want := string(wantBytes)
	if normalized != want {
		t.Fatalf("generated source does not match fixture %s\n--- got ---\n%s\n--- want ---\n%s", fixtureName, normalized, want)
	}
}

