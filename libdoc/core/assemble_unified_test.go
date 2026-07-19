package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestAssembleUnifiedLeafHierarchicalNoInline(t *testing.T) {
	rootBlock := &GoBlock{
		Types: map[string]bool{"Request": true, "Response": true},
		TypeDecls: []string{
			"type Request struct{ N int }",
			"type Response struct{ OK bool }",
		},
		Run: &FuncSnippet{
			Name:    "Run",
			Params:  "t *testing.T, req *Request",
			Results: "(*Response, error)",
			Body:    "{ return &Response{OK: true}, nil }",
		},
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ req.N = 1; return nil }",
		},
	}
	featureBlock := &GoBlock{
		Helpers: []FuncSnippet{{
			Name:    "featureMarkerHelper",
			Params:  "",
			Results: "string",
			Body:    `{ return "feature-marker" }`,
		}},
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ _ = featureMarkerHelper(); req.N += 10; return nil }",
		},
	}
	midBlock := &GoBlock{
		Helpers: []FuncSnippet{{
			Name:    "midMarkerHelper",
			Params:  "",
			Results: "string",
			Body:    `{ return "mid-marker" }`,
		}},
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ _ = midMarkerHelper(); req.N += 100; return nil }",
		},
	}
	leafBlock := &GoBlock{
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ req.N += 1000; return nil }",
		},
	}
	assertBlock := GoBlock{
		Assert: &FuncSnippet{
			Name:   "Assert",
			Params: "t *testing.T, req *Request, resp *Response, err error",
			Body:   "{ if err != nil { t.Fatal(err) }; if req.N != 1111 { t.Fatalf(\"N=%d\", req.N) } }",
		},
	}

	tc := TreeCase{
		Name: "leaf",
		Path: "feature/mid/leaf",
		SetupFiles: []SetupDocument{
			{Path: "DOCTEST.md", GoBlock: rootBlock},
			{Path: "feature/SETUP.md", GoBlock: featureBlock},
			{Path: "feature/mid/SETUP.md", GoBlock: midBlock},
			{Path: "feature/mid/leaf/SETUP.md", GoBlock: leafBlock},
		},
		AssertFile: AssertDocument{Path: "feature/mid/leaf/ASSERT.md", GoBlock: assertBlock},
	}

	leafSrc, err := AssembleUnifiedLeafSource(tc, false, "testcase", "/tmp/doctest-root", RefRootImportPath, RefRootPkgName, "testcase/__registry")
	if err != nil {
		t.Fatalf("AssembleUnifiedLeafSource: %v", err)
	}

	// Registry + RunTestLeaf contract preserved.
	if !strings.Contains(leafSrc, "func "+RunTestLeafName) {
		t.Fatalf("unified leaf must define %s:\n%s", RunTestLeafName, leafSrc)
	}
	if !strings.Contains(leafSrc, "registry.Register") {
		t.Fatalf("unified leaf must register with registry:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, "testcase/__registry") {
		t.Fatalf("unified leaf should import registry:\n%s", leafSrc)
	}

	// Must import intermediate packages (not inline their bodies).
	if !strings.Contains(leafSrc, "testcase/feature") {
		t.Fatalf("unified leaf should import feature intermediate:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, "testcase/feature/mid") {
		t.Fatalf("unified leaf should import mid intermediate:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, ".Setup(t, d, req)") {
		t.Fatalf("unified leaf should call intermediate Setup:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "feature-marker") {
		t.Fatalf("unified leaf must not inline feature marker body:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "mid-marker") {
		t.Fatalf("unified leaf must not inline mid marker body:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "func featureMarkerHelper") || strings.Contains(leafSrc, "func FeatureMarkerHelper") {
		t.Fatalf("unified leaf must not define feature helper:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "func midMarkerHelper") || strings.Contains(leafSrc, "func MidMarkerHelper") {
		t.Fatalf("unified leaf must not define mid helper:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "type Request") {
		t.Fatalf("unified leaf must not redefine type Request:\n%s", leafSrc)
	}
	// Leaf-local setup still present as closure.
	if !strings.Contains(leafSrc, "setup0") {
		t.Fatalf("leaf-local setup closure missing:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, "req.N += 1000") {
		t.Fatalf("leaf-local setup body should remain inlined:\n%s", leafSrc)
	}
}

func TestWriteUnifiedTreeWritesIntermediatePackages(t *testing.T) {
	rootBlock := &GoBlock{
		Types: map[string]bool{"Request": true, "Response": true},
		TypeDecls: []string{
			"type Request struct{}",
			"type Response struct{}",
		},
		Run: &FuncSnippet{
			Name:    "Run",
			Params:  "t *testing.T, req *Request",
			Results: "(*Response, error)",
			Body:    "{ return &Response{}, nil }",
		},
	}
	featureBlock := &GoBlock{
		Helpers: []FuncSnippet{{
			Name:    "parentOnlyHelper",
			Params:  "",
			Results: "int",
			Body:    "{ return 42 }",
		}},
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ _ = parentOnlyHelper(); return nil }",
		},
	}
	leafABlock := &GoBlock{
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ return nil }",
		},
	}
	leafBBlock := &GoBlock{
		Setup: &FuncSnippet{
			Name:    "Setup",
			Params:  "t *testing.T, req *Request",
			Results: "error",
			Body:    "{ return nil }",
		},
	}
	assertBlock := GoBlock{
		Assert: &FuncSnippet{
			Name:   "Assert",
			Params: "t *testing.T, req *Request, resp *Response, err error",
			Body:   "{ _ = req; _ = resp; _ = err }",
		},
	}

	// Two sibling leaves under the same intermediate parent — package written once.
	cases := []TreeCase{
		{
			Name: "leaf-a",
			Path: "feature/leaf-a",
			SetupFiles: []SetupDocument{
				{Path: "DOCTEST.md", GoBlock: rootBlock},
				{Path: "feature/SETUP.md", GoBlock: featureBlock},
				{Path: "feature/leaf-a/SETUP.md", GoBlock: leafABlock},
			},
			AssertFile: AssertDocument{Path: "feature/leaf-a/ASSERT.md", GoBlock: assertBlock},
		},
		{
			Name: "leaf-b",
			Path: "feature/leaf-b",
			SetupFiles: []SetupDocument{
				{Path: "DOCTEST.md", GoBlock: rootBlock},
				{Path: "feature/SETUP.md", GoBlock: featureBlock},
				{Path: "feature/leaf-b/SETUP.md", GoBlock: leafBBlock},
			},
			AssertFile: AssertDocument{Path: "feature/leaf-b/ASSERT.md", GoBlock: assertBlock},
		},
	}

	genRoot := t.TempDir()
	if err := WriteUnifiedTree(genRoot, cases, "/doctest", false, "testcase"); err != nil {
		t.Fatalf("WriteUnifiedTree: %v", err)
	}

	// Intermediate package on disk once at parent path.
	interPath := filepath.Join(genRoot, "feature", RefIntermediateFileName)
	interBytes, err := os.ReadFile(interPath)
	if err != nil {
		t.Fatalf("expected intermediate package at %s: %v", interPath, err)
	}
	interSrc := string(interBytes)
	if !strings.Contains(interSrc, "package feature") {
		t.Fatalf("intermediate package name:\n%s", interSrc)
	}
	if !strings.Contains(interSrc, "ParentOnlyHelper") {
		t.Fatalf("intermediate should export helper:\n%s", interSrc)
	}

	// Only one setup.go under feature/ (shared), not duplicated under each leaf.
	if _, err := os.Stat(filepath.Join(genRoot, "feature", "leaf-a", RefIntermediateFileName)); err == nil {
		t.Fatalf("intermediate must not be duplicated under leaf-a")
	}
	if _, err := os.Stat(filepath.Join(genRoot, "feature", "leaf-b", RefIntermediateFileName)); err == nil {
		t.Fatalf("intermediate must not be duplicated under leaf-b")
	}

	for _, tc := range cases {
		leafPath := filepath.Join(genRoot, tc.Path, UnifiedLeafFileName(tc))
		leafBytes, err := os.ReadFile(leafPath)
		if err != nil {
			t.Fatalf("read leaf %s: %v", leafPath, err)
		}
		leafSrc := string(leafBytes)
		if strings.Contains(leafSrc, "return 42") {
			t.Fatalf("leaf %s must not contain intermediate helper body:\n%s", tc.Path, leafSrc)
		}
		if !strings.Contains(leafSrc, "testcase/feature") {
			t.Fatalf("leaf %s should import intermediate path:\n%s", tc.Path, leafSrc)
		}
		if !strings.Contains(leafSrc, ".Setup(") {
			t.Fatalf("leaf %s should call intermediate Setup:\n%s", tc.Path, leafSrc)
		}
		if !strings.Contains(leafSrc, "func "+RunTestLeafName) {
			t.Fatalf("leaf %s missing %s:\n%s", tc.Path, RunTestLeafName, leafSrc)
		}
		if !strings.Contains(leafSrc, "registry.Register") {
			t.Fatalf("leaf %s missing registry.Register:\n%s", tc.Path, leafSrc)
		}
	}

	// Root + suite extras still written.
	if _, err := os.Stat(filepath.Join(genRoot, RefRootDirName, "droot.go")); err != nil {
		t.Fatalf("root package missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(genRoot, UnifiedRegistryDirName, "registry.go")); err != nil {
		t.Fatalf("registry missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(genRoot, UnifiedAllLeavesDirName, "all.go")); err != nil {
		t.Fatalf("allleaves missing: %v", err)
	}
	if _, err := os.Stat(filepath.Join(genRoot, UnifiedSuiteDirName, "suite_test.go")); err != nil {
		t.Fatalf("suite missing: %v", err)
	}

	// Full compile of generated suite against the real session package.
	modRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("module root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(modRoot, "go.mod")); err != nil {
		t.Logf("skip go test -c: module root not found at %s", modRoot)
		return
	}
	goMod := "module testcase\n\ngo 1.22\n\nrequire github.com/xhd2015/doctest v0.0.0\n\nreplace github.com/xhd2015/doctest => " + filepath.ToSlash(modRoot) + "\n"
	if err := os.WriteFile(filepath.Join(genRoot, "go.mod"), []byte(goMod), 0644); err != nil {
		t.Fatalf("write go.mod: %v", err)
	}
	tidy := exec.Command("go", "mod", "tidy")
	tidy.Dir = genRoot
	tidy.Env = append(os.Environ(), "GO111MODULE=on")
	if out, err := tidy.CombinedOutput(); err != nil {
		t.Fatalf("go mod tidy: %v\n%s", err, out)
	}
	cmd := exec.Command("go", "test", "./suite/", "-count=1", "-c", "-o", os.DevNull)
	cmd.Dir = genRoot
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	out, err := cmd.CombinedOutput()
	if err != nil {
		interBytes, _ = os.ReadFile(interPath)
		leafA, _ := os.ReadFile(filepath.Join(genRoot, "feature", "leaf-a", "leaf.go"))
		t.Fatalf("go test -c unified suite: %v\n%s\n--- feature/setup.go ---\n%s\n--- leaf-a ---\n%s",
			err, out, interBytes, leafA)
	}
}
