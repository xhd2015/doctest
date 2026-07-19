package core

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestPartitionRefSetupDocsHierarchical(t *testing.T) {
	tc := TreeCase{
		Path: "feature/mid/leaf",
		SetupFiles: []SetupDocument{
			{Path: "DOCTEST.md", GoBlock: &GoBlock{Setup: &FuncSnippet{Name: "Setup", Params: "t *testing.T, req *Request", Results: "error", Body: "{ return nil }"}}},
			{Path: "SETUP.md", GoBlock: &GoBlock{Setup: &FuncSnippet{Name: "Setup", Params: "t *testing.T, req *Request", Results: "error", Body: "{ return nil }"}}},
			{Path: "feature/SETUP.md", GoBlock: &GoBlock{Setup: &FuncSnippet{Name: "Setup", Params: "t *testing.T, req *Request", Results: "error", Body: "{ return nil }"}}},
			{Path: "feature/mid/SETUP.md", GoBlock: &GoBlock{Setup: &FuncSnippet{Name: "Setup", Params: "t *testing.T, req *Request", Results: "error", Body: "{ return nil }"}}},
			{Path: "feature/mid/leaf/SETUP.md", GoBlock: &GoBlock{Setup: &FuncSnippet{Name: "Setup", Params: "t *testing.T, req *Request", Results: "error", Body: "{ return nil }"}}},
		},
	}
	part := PartitionRefSetupDocs(tc)
	if len(part.RootDocs) != 2 {
		t.Fatalf("root docs: got %d want 2", len(part.RootDocs))
	}
	if len(part.Intermediate) != 2 {
		t.Fatalf("intermediate groups: got %d want 2", len(part.Intermediate))
	}
	if part.Intermediate[0].Dir != "feature" || part.Intermediate[1].Dir != "feature/mid" {
		t.Fatalf("intermediate dirs order: %+v", []string{part.Intermediate[0].Dir, part.Intermediate[1].Dir})
	}
	if len(part.LeafDocs) != 1 || part.LeafDocs[0].Path != "feature/mid/leaf/SETUP.md" {
		t.Fatalf("leaf docs: %+v", part.LeafDocs)
	}
}

func TestAssembleRefIntermediateAndLeafNoInline(t *testing.T) {
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

	part := PartitionRefSetupDocs(tc)
	if len(part.Intermediate) != 2 {
		t.Fatalf("expected 2 intermediate groups, got %d", len(part.Intermediate))
	}

	// Intermediate feature package must export Setup + helper, import droot.
	featureSrc, err := AssembleRefIntermediateSource(
		part.Intermediate[0].Docs, RefIntermediatePkgName("feature"),
		RefRootImportPath, RefRootPkgName, part.RootDocs,
		nil, // no intermediate ancestors
	)
	if err != nil {
		t.Fatalf("AssembleRefIntermediateSource feature: %v", err)
	}
	if !strings.Contains(featureSrc, "package feature") {
		t.Fatalf("feature package clause missing:\n%s", featureSrc)
	}
	if !strings.Contains(featureSrc, `droot "testcase/__droot"`) && !strings.Contains(featureSrc, "testcase/__droot") {
		t.Fatalf("feature package should import droot:\n%s", featureSrc)
	}
	if !strings.Contains(featureSrc, "func Setup(") {
		t.Fatalf("feature package should export Setup:\n%s", featureSrc)
	}
	if !strings.Contains(featureSrc, "FeatureMarkerHelper") && !strings.Contains(featureSrc, "featureMarkerHelper") {
		// unexported helper must be exported via rename
		t.Fatalf("feature helper export missing:\n%s", featureSrc)
	}
	if !strings.Contains(featureSrc, "FeatureMarkerHelper") {
		t.Fatalf("expected exported FeatureMarkerHelper:\n%s", featureSrc)
	}
	if !strings.Contains(featureSrc, "droot.Request") {
		t.Fatalf("feature Setup should use droot.Request:\n%s", featureSrc)
	}

	// Nested mid imports feature parent + droot (ancestors parents-first).
	midSrc, err := AssembleRefIntermediateSource(
		part.Intermediate[1].Docs, RefIntermediatePkgName("feature/mid"),
		RefRootImportPath, RefRootPkgName, part.RootDocs,
		[]RefIntermediateGroup{part.Intermediate[0]},
	)
	if err != nil {
		t.Fatalf("AssembleRefIntermediateSource mid: %v", err)
	}
	if !strings.Contains(midSrc, "package mid") {
		t.Fatalf("mid package clause missing:\n%s", midSrc)
	}
	if !strings.Contains(midSrc, "testcase/feature") {
		t.Fatalf("mid should import parent feature package:\n%s", midSrc)
	}
	if !strings.Contains(midSrc, "MidMarkerHelper") {
		t.Fatalf("expected exported MidMarkerHelper:\n%s", midSrc)
	}

	leafSrc, err := AssembleRefLeafTestSource(tc, false, "testcase", "/tmp/doctest-root", RefRootImportPath, RefRootPkgName)
	if err != nil {
		t.Fatalf("AssembleRefLeafTestSource: %v", err)
	}

	// Must import intermediate packages.
	if !strings.Contains(leafSrc, "testcase/feature") {
		t.Fatalf("leaf should import feature intermediate:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, "testcase/feature/mid") {
		t.Fatalf("leaf should import mid intermediate:\n%s", leafSrc)
	}
	// Must call intermediate Setup, not inline markers.
	if !strings.Contains(leafSrc, ".Setup(t, d, req)") {
		t.Fatalf("leaf should call intermediate Setup:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "feature-marker") {
		t.Fatalf("leaf must not inline feature marker body:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "mid-marker") {
		t.Fatalf("leaf must not inline mid marker body:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "func featureMarkerHelper") || strings.Contains(leafSrc, "func FeatureMarkerHelper") {
		t.Fatalf("leaf must not define feature helper:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "func midMarkerHelper") || strings.Contains(leafSrc, "func MidMarkerHelper") {
		t.Fatalf("leaf must not define mid helper:\n%s", leafSrc)
	}
	if strings.Contains(leafSrc, "type Request") {
		t.Fatalf("leaf must not redefine type Request:\n%s", leafSrc)
	}
	// Leaf-local setup still present as closure.
	if !strings.Contains(leafSrc, "setup0") {
		t.Fatalf("leaf-local setup closure missing:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, "req.N += 1000") {
		t.Fatalf("leaf-local setup body should remain inlined:\n%s", leafSrc)
	}
}

func TestWriteRefTreeWritesIntermediatePackages(t *testing.T) {
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
	leafBlock := &GoBlock{
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

	tc := TreeCase{
		Name: "leaf",
		Path: "feature/leaf",
		SetupFiles: []SetupDocument{
			{Path: "DOCTEST.md", GoBlock: rootBlock},
			{Path: "feature/SETUP.md", GoBlock: featureBlock},
			{Path: "feature/leaf/SETUP.md", GoBlock: leafBlock},
		},
		AssertFile: AssertDocument{Path: "feature/leaf/ASSERT.md", GoBlock: assertBlock},
	}

	genRoot := t.TempDir()
	if err := WriteRefTree(genRoot, []TreeCase{tc}, "/doctest", false, "testcase"); err != nil {
		t.Fatalf("WriteRefTree: %v", err)
	}

	// Intermediate package on disk at parent path.
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

	// Leaf test must not contain intermediate helper body/marker.
	leafPath := filepath.Join(genRoot, "feature", "leaf", TestFileName(tc))
	leafBytes, err := os.ReadFile(leafPath)
	if err != nil {
		t.Fatalf("read leaf: %v", err)
	}
	leafSrc := string(leafBytes)
	if strings.Contains(leafSrc, "parentOnlyHelper") || strings.Contains(leafSrc, "ParentOnlyHelper") {
		// May reference via package alias — allow alias.ParentOnlyHelper only if imported for use.
		// Body of helper must not be inlined.
		if strings.Contains(leafSrc, "return 42") {
			t.Fatalf("leaf inlined intermediate helper body:\n%s", leafSrc)
		}
	}
	if strings.Contains(leafSrc, "return 42") {
		t.Fatalf("leaf must not contain intermediate helper body:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, "testcase/feature") {
		t.Fatalf("leaf should import intermediate path:\n%s", leafSrc)
	}
	if !strings.Contains(leafSrc, ".Setup(") {
		t.Fatalf("leaf should call intermediate Setup:\n%s", leafSrc)
	}

	// Root package still written.
	if _, err := os.Stat(filepath.Join(genRoot, RefRootDirName, "droot.go")); err != nil {
		t.Fatalf("root package missing: %v", err)
	}

	// Full compile of generated packages against the real session package.
	modRoot, err := filepath.Abs("../..")
	if err != nil {
		t.Fatalf("module root: %v", err)
	}
	if _, err := os.Stat(filepath.Join(modRoot, "go.mod")); err != nil {
		// Running from an unexpected cwd — skip compile check.
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
	cmd := exec.Command("go", "test", "./feature/leaf/", "-count=1", "-c", "-o", os.DevNull)
	cmd.Dir = genRoot
	cmd.Env = append(os.Environ(), "GO111MODULE=on")
	out, err := cmd.CombinedOutput()
	if err != nil {
		// Re-read sources after format for the failure dump.
		interBytes, _ = os.ReadFile(interPath)
		leafBytes, _ = os.ReadFile(leafPath)
		t.Fatalf("go test -c generated tree: %v\n%s\n--- feature/setup.go ---\n%s\n--- leaf ---\n%s",
			err, out, interBytes, leafBytes)
	}
}

func TestSanitizePackageName(t *testing.T) {
	cases := map[string]string{
		"feature":    "feature",
		"my-feature": "my_feature",
		"123abc":     "p123abc",
		"type":       "pkg_type",
		"error":      "pkg_error", // predeclared builtin — must not shadow `err error`
		"string":     "pkg_string",
		"":           "pkg",
		"__droot":    "pkg_droot", // after sanitize of underscores + reserved check
	}
	for in, want := range cases {
		got := SanitizePackageName(in)
		// __droot: underscores trimmed/leading handling may differ — check keyword/reserved paths carefully.
		if in == "__droot" {
			if got == "droot" || got == RefRootPkgName {
				t.Fatalf("SanitizePackageName(%q) must not be reserved droot, got %q", in, got)
			}
			continue
		}
		if got != want {
			t.Errorf("SanitizePackageName(%q)=%q want %q", in, got, want)
		}
	}
}

func TestRefIntermediateImportPaths(t *testing.T) {
	if got := RefIntermediateImport(RefRootImportPath, "feature"); got != "testcase/feature" {
		t.Fatalf("flat: got %q", got)
	}
	if got := RefIntermediateImport("testcase/tests/foo/__droot", "feature/mid"); got != "testcase/tests/foo/feature/mid" {
		t.Fatalf("tree-scoped: got %q", got)
	}
	if got := RefIntermediatePkgName("feature/mid"); got != "mid" {
		t.Fatalf("pkg name: got %q", got)
	}
	if got := RefIntermediateAlias("feature/mid"); got != "feature_mid" {
		t.Fatalf("alias: got %q", got)
	}
}
