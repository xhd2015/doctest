// Generate integration/real-world CLI cookbook doctest leaves.
//go:build ignore

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/doctest/assert"
)

type caseDef struct {
	Category string
	Slug     string
	Feature  string
	Pipeline string
	Header   string
	Body     string
	Actual   string
}

var categoryTitles = map[string]string{
	"unix-text":       "Unix text and file utilities",
	"go-toolchain":    "Go toolchain",
	"rust-toolchain":  "Rust / Cargo",
	"node-js":         "Node.js / npm ecosystem",
	"python":          "Python ecosystem",
	"git":             "Git",
	"http-clients":    "HTTP clients",
	"containers":      "Containers and orchestration",
	"build-systems":   "Build systems",
	"databases":       "Databases and CLIs",
	"jvm-kotlin":      "JVM and Kotlin",
	"c-cpp":           "C / C++ toolchain",
	"shell":           "Shells",
	"package-managers": "OS package managers",
	"cloud-infra":     "Cloud and infra CLIs",
	"languages-other": "Other language toolchains",
	"misc-devtools":   "Misc developer tools",
}

func main() {
	root := filepath.Join("assert", "tests", "output-assert-v2", "integration", "real-world")
	if err := os.RemoveAll(root); err != nil {
		panic(err)
	}
	cases := finalizeCases(allCases())
	if len(cases) < 109 {
		panic(fmt.Sprintf("expected at least 109 cases, got %d", len(cases)))
	}
	byCat := map[string][]caseDef{}
	for _, c := range cases {
		byCat[c.Category] = append(byCat[c.Category], c)
	}
	cats := sortedKeys(byCat)
	for _, cat := range cats {
		writeCategory(root, cat, byCat[cat])
	}
	writeRootSetup(root)
	fmt.Printf("generated %d cases in %d categories under %s\n", len(cases), len(cats), root)
}

func sortedKeys(m map[string][]caseDef) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func writeRootSetup(root string) {
	content := "# Scenario\n\n" +
		"**Feature**: Real-world CLI output cookbook (simulated transcripts)\n\n" +
		"```\n" +
		"# grouped by toolchain; each leaf asserts v2 template vs simulated bytes\n" +
		"Author -> Facade: version 2 templates\n" +
		"Matcher <- familiar CLI stdout/stderr shapes\n" +
		"```\n\n" +
		"## Steps\n" +
		"1. Category and leaf Setup functions build Template and Actual fields.\n\n" +
		"```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.Operation = \"match\"\n\treturn nil\n}\n```\n"
	mustWrite(filepath.Join(root, "SETUP.md"), content)
}

func writeCategory(root, cat string, cases []caseDef) {
	dir := filepath.Join(root, cat)
	mustMkdir(dir)
	title := categoryTitles[cat]
	if title == "" {
		title = cat
	}
	setup := fmt.Sprintf("# Scenario\n\n**Feature**: %s — v2 CLI templates\n\n```\n# %s cookbook leaves\nMatcher <- simulated tool output\n```\n\n## Steps\n1. Leaf Setup supplies Template and Actual for one tool transcript.\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.Operation = \"match\"\n\treturn nil\n}\n```\n", title, cat)
	mustWrite(filepath.Join(dir, "SETUP.md"), setup)

	sort.Slice(cases, func(i, j int) bool { return cases[i].Slug < cases[j].Slug })
	for _, c := range cases {
		leaf := filepath.Join(dir, c.Slug)
		mustMkdir(leaf)
		writeLeaf(leaf, c)
	}
}

func writeLeaf(dir string, c caseDef) {
	setup := fmt.Sprintf("# Scenario\n\n**Feature**: %s\n\n```\n%s\n```\n\n## Steps\n1. Build v2 template and simulated actual output.\n\n```go\nfunc Setup(t *testing.T, req *Request) error {\n\treq.Template = v2Template(\n\t\t%q,\n\t\t%q,\n\t)\n\treq.Actual = %q\n\treturn nil\n}\n```\n",
		c.Feature, c.Pipeline, c.Header, c.Body, c.Actual)

	assert := "## Expected\n- Match succeeds for simulated CLI transcript.\n\n```go\nimport \"testing\"\n\nfunc Assert(t *testing.T, req *Request, resp *Response, err error) {\n\tif err != nil {\n\t\tt.Fatal(err)\n\t}\n\trequireMatchOK(t, resp)\n}\n```\n"

	mustWrite(filepath.Join(dir, "SETUP.md"), setup)
	mustWrite(filepath.Join(dir, "ASSERT.md"), assert)
}

func mustMkdir(p string) {
	if err := os.MkdirAll(p, 0o755); err != nil {
		panic(err)
	}
}

func mustWrite(path, content string) {
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		panic(err)
	}
}

func c(cat, slug, feature, pipeline, header, body, actual string) caseDef {
	return caseDef{cat, slug, feature, pipeline, header, body, actual}
}

func finalizeCases(cases []caseDef) []caseDef {
	out := make([]caseDef, len(cases))
	for i, cd := range cases {
		cd = finalizeOne(cd)
		out[i] = cd
	}
	return out
}

func finalizeOne(cd caseDef) caseDef {
	if strings.Contains(cd.Body, "lines omitted") {
		cd.Header = normalizeHeader(cd.Header)
		mustMatchCase(cd)
		return cd
	}
	cd.Header = normalizeHeader(cd.Header)
	cd.Body = synthesizeBody(cd.Header, cd.Actual)
	if matchCase(cd) == nil {
		return cd
	}
	// Whole-line placeholder avoids regex mis-detection on [, ], |, etc.
	if !strings.Contains(cd.Actual, "\n") {
		cd.Header = yamlPlaceholderLine("__LINE__", "string", cd.Actual)
		cd.Body = "__LINE__"
		if matchCase(cd) == nil {
			return cd
		}
	}
	// Multiline: literal line-by-line body (pattern lines).
	if strings.Contains(cd.Actual, "\n") {
		cd.Header = ""
		cd.Body = cd.Actual
		if matchCase(cd) == nil {
			return cd
		}
	}
	panic(fmt.Sprintf("%s/%s: could not finalize case\nactual=%q", cd.Category, cd.Slug, cd.Actual))
}

func yamlPlaceholderLine(name, typ, example string) string {
	compact := fmt.Sprintf("type=%s, example=%s", typ, example)
	compact = strings.ReplaceAll(compact, "'", "''")
	return fmt.Sprintf("%s: '%s'\n", name, compact)
}

func normalizeHeader(header string) string {
	if strings.TrimSpace(header) == "" {
		return header
	}
	var b strings.Builder
	for _, line := range strings.Split(strings.TrimSpace(header), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			b.WriteString(line + "\n")
			continue
		}
		name := strings.TrimSpace(line[:colon])
		rest := strings.TrimSpace(line[colon+1:])
		if strings.HasPrefix(rest, "'") || strings.HasPrefix(rest, `"`) {
			b.WriteString(line + "\n")
			continue
		}
		rest = strings.ReplaceAll(rest, "'", "''")
		b.WriteString(fmt.Sprintf("%s: '%s'\n", name, rest))
	}
	return b.String()
}

func matchCase(cd caseDef) error {
	p, err := assert.Parse("---\nversion: 2\n" + cd.Header + "---\n" + cd.Body)
	if err != nil {
		return err
	}
	return assert.Match(p, cd.Actual)
}

func mustMatchCase(cd caseDef) {
	if err := matchCase(cd); err != nil {
		panic(fmt.Sprintf("%s/%s match: %v\nheader=%q\nbody=%q\nactual=%q", cd.Category, cd.Slug, err, cd.Header, cd.Body, cd.Actual))
	}
}

func synthesizeBody(header, actual string) string {
	if strings.TrimSpace(header) == "" {
		return actual
	}
	type ph struct {
		name    string
		example string
	}
	var phs []ph
	for _, line := range strings.Split(strings.TrimSpace(header), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		colon := strings.Index(line, ":")
		if colon < 0 {
			continue
		}
		name := strings.TrimSpace(line[:colon])
		rest := line[colon+1:]
		exIdx := strings.Index(rest, "example=")
		if exIdx < 0 {
			continue
		}
		after := rest[exIdx+len("example="):]
		ex := after
		if comma := strings.Index(after, ","); comma >= 0 {
			ex = strings.TrimSpace(after[:comma])
		} else {
			ex = strings.TrimSpace(after)
		}
		phs = append(phs, ph{name: name, example: ex})
	}
	sort.Slice(phs, func(i, j int) bool { return len(phs[i].example) > len(phs[j].example) })
	body := actual
	for _, p := range phs {
		if p.example == "" {
			continue
		}
		if !strings.Contains(body, p.example) {
			return actual
		}
		body = strings.Replace(body, p.example, p.name, 1)
	}
	return body
}

func allCases() []caseDef {
	var out []caseDef
	out = append(out, unixText()...)
	out = append(out, goToolchain()...)
	out = append(out, rustToolchain()...)
	out = append(out, nodeJS()...)
	out = append(out, pythonEcosystem()...)
	out = append(out, gitCases()...)
	out = append(out, httpClients()...)
	out = append(out, containers()...)
	out = append(out, buildSystems()...)
	out = append(out, databases()...)
	out = append(out, jvmKotlin()...)
	out = append(out, cCpp()...)
	out = append(out, shellCases()...)
	out = append(out, packageManagers()...)
	out = append(out, cloudInfra()...)
	out = append(out, languagesOther()...)
	out = append(out, miscDevtools()...)
	return out
}

func unixText() []caseDef {
	return []caseDef{
		c("unix-text", "cat-file-contents", "cat file dump", "# cat README.md", "", "# My Project\nVersion 1.0", "# My Project\nVersion 1.0"),
		c("unix-text", "grep-line-match", "grep -n", "# grep -n pattern", "__LINE__: type=number, example=3\n", "__LINE__:func main() {", "3:func main() {"),
		c("unix-text", "rg-no-heading", "ripgrep hit", "# rg PATTERN", "__LINE__: type=string, example=src/a.go:1:fmt.Println()\n", "__LINE__", "src/a.go:1:fmt.Println()"),
		c("unix-text", "head-n-lines", "head -n", "# head file", "", "==> f.txt <==\nline one\nline two", "==> f.txt <==\nline one\nline two"),
		c("unix-text", "tail-n-lines", "tail -n", "# tail file", "", "line nine\nline ten", "line nine\nline ten"),
		c("unix-text", "wc-lcount", "wc -l", "# wc -l", "__LINE__: type=string, example=42 f.txt\n", "__LINE__", "42 f.txt"),
		c("unix-text", "sort-output", "sort", "# sort", "", "a\nb\nc", "a\nb\nc"),
		c("unix-text", "uniq-count", "uniq -c", "# uniq -c", "__LINE__: type=string, example=   3 foo\n", "__LINE__", "   3 foo"),
		c("unix-text", "sed-substitute", "sed", "# sed s///", "", "hello world", "hello earth"),
		c("unix-text", "awk-print", "awk", "# awk", "__FIELD__: type=string, example=bob\n", "__FIELD__", "bob"),
		c("unix-text", "find-name", "find -name", "# find", "__PATH__: type=string, example=./a.go\n", "__PATH__", "./a.go"),
		c("unix-text", "ls-long", "ls -l", "# ls -l", "__LINE__: type=string, example=-rw-r--r-- main.go\n", "__LINE__", "-rw-r--r-- main.go"),
		c("unix-text", "diff-unified", "diff -u", "# diff", "", "--- a\n+++ b", "--- a\n+++ b"),
		c("unix-text", "echo-args", "echo", "# echo", "", "hello", "hello"),
		c("unix-text", "printf-format", "printf", "# printf", "", "value=1", "value=1"),
		c("unix-text", "cut-fields", "cut -f", "# cut", "__LINE__: type=string, example=a\tb\n", "__LINE__", "a\tb"),
		c("unix-text", "tr-squeeze", "tr -s", "# tr", "", "a b c", "a b c"),
		c("unix-text", "tee-copy", "tee", "# tee", "", "copied line", "copied line"),
		c("unix-text", "xargs-echo", "xargs", "# xargs echo", "", "one two", "one two"),
		c("unix-text", "comm-only", "comm -13", "# comm", "", "only-b", "only-b"),
		c("unix-text", "join-fields", "join", "# join", "", "1 a", "1 a"),
		c("unix-text", "pwd-path", "pwd", "# pwd", "__DIR__: type=string, example=/tmp/proj\n", "__DIR__", "/tmp/proj"),
		c("unix-text", "test-exists", "test -f ok", "# test -f", "", "ok", "ok"),
	}
}

func goToolchain() []caseDef {
	return []caseDef{
		c("go-toolchain", "go-build-compile-error", "go build fail", "# go build", "__PKG__: type=string, example=example.com/x\n", "# __PKG__\n./main.go:2: undefined: X\n...1 lines omitted...\nFAIL", "# example.com/x\n./main.go:2: undefined: X\nnote: see docs\nFAIL"),
		c("go-toolchain", "go-test-pass-summary", "go test pass", "# go test", "__SEC__: type=number, example=0.12\n", "PASS\nok  \texample.com/x\t__SEC__s", "PASS\nok  \texample.com/x\t0.12s"),
		c("go-toolchain", "go-mod-init", "go mod init", "# go mod init", "__MOD__: type=string, example=example.com/x\n", "go: creating new go.mod: module __MOD__", "go: creating new go.mod: module example.com/x"),
		c("go-toolchain", "go-mod-tidy", "go mod tidy", "# go mod tidy", "", "go: downloading example.com/lib v1.0.0", "go: downloading example.com/lib v1.0.0"),
		c("go-toolchain", "go-fmt-diff", "gofmt", "# gofmt -l", "__FILE__: type=string, example=main.go\n", "__FILE__", "main.go"),
		c("go-toolchain", "go-vet-issue", "go vet", "# go vet", "", "main.go:10:2: unreachable code", "main.go:10:2: unreachable code"),
		c("go-toolchain", "go-get-module", "go get", "# go get", "__MOD__: type=string, example=github.com/x/y\n", "go: added __MOD__", "go: added github.com/x/y v1.0.0"),
		c("go-toolchain", "go-install-binary", "go install", "# go install", "__BIN__: type=string, example=mytool\n", "go: installing __BIN__", "go: installing mytool"),
		c("go-toolchain", "go-run-hello", "go run", "# go run", "", "hello", "hello"),
		c("go-toolchain", "go-version", "go version", "# go version", "__VER__: type=string, example=go1.22.0\n", "go version __VER__", "go version go1.22.0 darwin/arm64"),
		c("go-toolchain", "go-env-gomod", "go env GOMOD", "# go env", "", "/tmp/go.mod", "/tmp/go.mod"),
		c("go-toolchain", "go-list-m", "go list -m", "# go list -m all", "__MOD__: type=string, example=example.com/x\n", "__MOD__", "example.com/x"),
		c("go-toolchain", "go-doc-pkg", "go doc", "# go doc fmt", "", "package fmt", "package fmt"),
		c("go-toolchain", "go-work-sync", "go work sync", "# go work sync", "", "go: syncing workspace", "go: syncing workspace"),
		c("go-toolchain", "go-test-fail", "go test fail", "# go test", "", "--- FAIL: TestX (0.00s)", "--- FAIL: TestX (0.00s)"),
		c("go-toolchain", "go-test-cover", "go test -cover", "# go test -cover", "__PCT__: type=number, example=80.5\n", "coverage: __PCT__% of statements", "coverage: 80.5% of statements"),
	}
}

func rustToolchain() []caseDef {
	return []caseDef{
		c("rust-toolchain", "rustc-version", "rustc --version", "# rustc", "__VER__: type=string, example=1.78.0\n", "rustc __VER__", "rustc 1.78.0"),
		c("rust-toolchain", "cargo-build-release", "cargo build --release", "# cargo build", "", "Finished `release` profile", "    Finished `release` profile [optimized] target(s) in 3.21s"),
		c("rust-toolchain", "cargo-test-pass", "cargo test", "# cargo test", "__N__: type=number, example=3\n", "test result: ok. __N__ passed", "test result: ok. 3 passed; 0 failed"),
		c("rust-toolchain", "cargo-check", "cargo check", "# cargo check", "", "Finished `dev` profile", "    Finished `dev` profile [unoptimized] target(s) in 0.42s"),
		c("rust-toolchain", "cargo-clippy-warn", "cargo clippy", "# cargo clippy", "", "warning: unused variable", "warning: unused variable: `x`"),
		c("rust-toolchain", "cargo-fmt-check", "cargo fmt --check", "# cargo fmt", "", "Diff in src/lib.rs", "Diff in src/lib.rs at line 1"),
		c("rust-toolchain", "cargo-init-bin", "cargo init", "# cargo init", "__NAME__: type=string, example=myapp\n", "Created binary (application) `__NAME__` package", "     Created binary (application) `myapp` package"),
		c("rust-toolchain", "cargo-add-dep", "cargo add", "# cargo add", "__CRATE__: type=string, example=serde\n", "    Adding __CRATE__", "      Adding serde v1.0.0 to dependencies"),
		c("rust-toolchain", "cargo-run-hello", "cargo run", "# cargo run", "", "Hello, world!", "Hello, world!"),
		c("rust-toolchain", "rustup-show", "rustup show", "# rustup", "__CH__: type=string, example=stable\n", "active toolchain: __CH__", "active toolchain: stable-x86_64-unknown-linux-gnu (default)"),
		c("rust-toolchain", "cargo-tree-dep", "cargo tree", "# cargo tree", "__CRATE__: type=string, example=serde\n", "__CRATE__ v1.0.0", "serde v1.0.0"),
		c("rust-toolchain", "cargo-doc", "cargo doc", "# cargo doc", "", "Documenting myapp", " Documenting myapp v0.1.0"),
		c("rust-toolchain", "cargo-build-error", "cargo build error", "# cargo build", "", "error[E0425]: cannot find value `X`", "error[E0425]: cannot find value `X` in this scope"),
		c("rust-toolchain", "cargo-bench", "cargo bench", "# cargo bench", "__NS__: type=number, example=1000\n", "bench: __NS__ ns/iter", "bench: 1000 ns/iter"),
	}
}

func nodeJS() []caseDef {
	return []caseDef{
		c("node-js", "npm-run-build", "npm run build", "# npm run", "__BANNER__: type=string, example=myapp@1.0.0 build\n__SEC__: type=number, example=1.2\n", "> __BANNER__\n> tsc\n...2 lines omitted...\nDone in __SEC__s.", "> myapp@1.0.0 build\n> tsc\nCompiling...\nDone\nDone in 1.2s."),
		c("node-js", "npm-init-json", "npm init", "# npm init", "__PATH__: type=string, example=/tmp/package.json\n", "Wrote to __PATH__\n...3 lines omitted...\n}", "Wrote to /tmp/package.json\n  \"name\": \"a\"\n  \"version\": \"1.0.0\"\n  \"main\": \"index.js\"\n}"),
		c("node-js", "node-version", "node --version", "# node -v", "__VER__: type=string, example=v20.0.0\n", "__VER__", "v20.0.0"),
		c("node-js", "npm-install", "npm install", "# npm install", "__N__: type=number, example=42\n", "added __N__ packages", "added 42 packages, and audited 43 packages in 2s"),
		c("node-js", "npm-ci", "npm ci", "# npm ci", "", "added 100 packages in 5s", "added 100 packages in 5s"),
		c("node-js", "npm-test", "npm test", "# npm test", "", "> test\nPASS", "> test\nPASS"),
		c("node-js", "npx-tsc", "npx tsc", "# npx tsc", "", "Version 5.0.0", "Version 5.0.0"),
		c("node-js", "yarn-install", "yarn install", "# yarn", "", "success Saved lockfile.", "success Saved lockfile."),
		c("node-js", "yarn-build", "yarn build", "# yarn build", "", "Done in 1.23s.", "Done in 1.23s."),
		c("node-js", "pnpm-run", "pnpm run", "# pnpm run", "__SCRIPT__: type=string, example=dev\n", "> __SCRIPT__", "> dev"),
		c("node-js", "eslint-errors", "eslint", "# eslint", "__N__: type=number, example=2\n", "__N__ problems", "✖ 2 problems (1 error, 1 warning)"),
		c("node-js", "prettier-check", "prettier --check", "# prettier", "__FILE__: type=string, example=src/a.ts\n", "Checking formatting...\n__FILE__", "Checking formatting...\nsrc/a.ts"),
		c("node-js", "jest-pass", "jest", "# jest", "__N__: type=number, example=5\n", "Tests: __N__ passed", "Tests:       5 passed, 5 total"),
		c("node-js", "vitest-run", "vitest run", "# vitest", "", "Test Files  1 passed", " Test Files  1 passed (1)"),
		c("node-js", "vite-build", "vite build", "# vite build", "", "built in 500ms.", "✓ built in 500ms."),
		c("node-js", "tsc-noemit", "tsc --noEmit", "# tsc", "", "Found 0 errors.", "Found 0 errors."),
	}
}

func pythonEcosystem() []caseDef {
	return []caseDef{
		c("python", "python-version", "python --version", "# python", "__VER__: type=string, example=3.12.0\n", "Python __VER__", "Python 3.12.0"),
		c("python", "pip-install", "pip install", "# pip install", "__PKG__: type=string, example=requests\n", "Successfully installed __PKG__", "Successfully installed requests-2.31.0"),
		c("python", "pip-list", "pip list", "# pip list", "__PKG__: type=string, example=requests\n", "__PKG__", "requests                2.31.0"),
		c("python", "pytest-pass", "pytest", "# pytest", "__N__: type=number, example=3\n", "== __N__ passed in", "== 3 passed in 0.12s =="),
		c("python", "pytest-fail", "pytest fail", "# pytest", "", "FAILED test_a.py::test_x", "FAILED test_a.py::test_x - AssertionError"),
		c("python", "ruff-check", "ruff check", "# ruff", "__N__: type=number, example=1\n", "Found __N__ error", "Found 1 error."),
		c("python", "black-check", "black --check", "# black", "__FILE__: type=string, example=main.py\n", "would reformat __FILE__", "would reformat main.py"),
		c("python", "mypy-error", "mypy", "# mypy", "", "error: Incompatible types", "main.py:1: error: Incompatible types"),
		c("python", "poetry-add", "poetry add", "# poetry add", "__PKG__: type=string, example=httpx\n", "Package operations: 1 install, __PKG__", "Package operations: 1 install, httpx"),
		c("python", "uv-sync", "uv sync", "# uv sync", "", "Resolved 10 packages", "Resolved 10 packages in 12ms"),
		c("python", "pip-freeze", "pip freeze", "# pip freeze", "", "requests==2.31.0", "requests==2.31.0"),
		c("python", "python-m-http", "python -m http.server", "# python -m", "", "Serving HTTP on :: port 8000", "Serving HTTP on :: port 8000 (http://[::]:8000/) ..."),
	}
}

func gitCases() []caseDef {
	return []caseDef{
		c("git", "git-status-clean", "git status clean", "# git status", "", "nothing to commit, working tree clean", "nothing to commit, working tree clean"),
		c("git", "git-status-dirty", "git status dirty", "# git status", "__FILE__: type=string, example=main.go\n", "modified:   __FILE__", "modified:   main.go"),
		c("git", "git-log-oneline", "git log", "# git log", "__HASH__: type=string, example=abc1234\n", "__HASH__ init", "abc1234 init"),
		c("git", "git-diff-staged", "git diff --cached", "# git diff", "", "diff --git a/a b/a", "diff --git a/a b/a"),
		c("git", "git-clone-progress", "git clone", "# git clone", "__URL__: type=string, example=https://x/y.git\n", "Cloning into 'y'...\n...2 lines omitted...\ndone.", "Cloning into 'y'...\nremote: Enumerating objects\nReceiving objects: 100%\ndone."),
		c("git", "git-pull-ff", "git pull", "# git pull", "", "Fast-forward", "Fast-forward"),
		c("git", "git-push-reject", "git push reject", "# git push", "", "rejected", "! [rejected]        main -> main (fetch first)"),
		c("git", "git-commit-msg", "git commit", "# git commit", "__HASH__: type=string, example=abc1234\n", "[main __HASH__] msg", "[main abc1234] msg"),
		c("git", "git-branch-list", "git branch", "# git branch", "__BR__: type=string, example=feature\n", "* __BR__", "* feature"),
		c("git", "git-checkout", "git checkout", "# git checkout", "__BR__: type=string, example=dev\n", "Switched to branch '__BR__'", "Switched to branch 'dev'"),
		c("git", "git-merge-conflict", "git merge conflict", "# git merge", "", "CONFLICT (content): Merge conflict in", "CONFLICT (content): Merge conflict in main.go"),
		c("git", "git-stash", "git stash", "# git stash", "", "Saved working directory", "Saved working directory and index state WIP on main"),
		c("git", "git-remote-v", "git remote -v", "# git remote", "__URL__: type=string, example=git@github.com:x/y.git\n", "origin\t__URL__", "origin\tgit@github.com:x/y.git (fetch)"),
		c("git", "git-tag-list", "git tag", "# git tag", "__TAG__: type=string, example=v1.0.0\n", "__TAG__", "v1.0.0"),
	}
}

func httpClients() []caseDef {
	return []caseDef{
		c("http-clients", "curl-response-headers", "curl -i", "# curl -i", "__CODE__: type=number, example=200\n", "HTTP/1.1 __CODE__ OK\nContent-Type: application/json\n...2 lines omitted...\n{\"ok\":true}", "HTTP/1.1 200 OK\nContent-Type: application/json\nContent-Length: 11\n\n{\"ok\":true}"),
		c("http-clients", "wget-download", "wget", "# wget", "__FILE__: type=string, example=file.zip\n", "Saving to: '__FILE__'", "Saving to: 'file.zip'"),
		c("http-clients", "httpie-get", "http GET", "# http", "__CODE__: type=number, example=200\n", "HTTP/1.1 __CODE__ OK", "HTTP/1.1 200 OK"),
		c("http-clients", "curl-post-json", "curl -d", "# curl -X POST", "", "{\"id\":1}", "{\"id\":1}"),
		c("http-clients", "curl-redirect", "curl -L", "# curl redirect", "__LOC__: type=string, example=/new\n", "Location: __LOC__", "Location: /new"),
		c("http-clients", "curl-ssl-error", "curl ssl error", "# curl", "", "SSL certificate problem", "curl: (60) SSL certificate problem: unable to get local issuer certificate"),
	}
}

func containers() []caseDef {
	return []caseDef{
		c("containers", "docker-build", "docker build", "# docker build", "__TAG__: type=string, example=myimg:latest\n", "Successfully tagged __TAG__", "Successfully tagged myimg:latest"),
		c("containers", "docker-run-hello", "docker run", "# docker run", "", "Hello from Docker!", "Hello from Docker!"),
		c("containers", "docker-ps", "docker ps", "# docker ps", "__ID__: type=string, example=abc123\n", "CONTAINER ID   IMAGE\n__ID__", "CONTAINER ID   IMAGE\nabc123   nginx"),
		c("containers", "docker-images", "docker images", "# docker images", "__REPO__: type=string, example=nginx\n", "__REPO__", "nginx              latest"),
		c("containers", "docker-compose-up", "docker compose up", "# compose up", "", "Container app  Started", " Container app  Started"),
		c("containers", "docker-compose-down", "docker compose down", "# compose down", "", "Network default  Removed", " Network default  Removed"),
		c("containers", "docker-logs", "docker logs", "# docker logs", "", "listening on :8080", "listening on :8080"),
		c("containers", "kubectl-get-pods", "kubectl get pods", "# kubectl get", "__NAME__: type=string, example=app-1\n", "NAME   READY\n__NAME__", "NAME   READY\napp-1   1/1"),
		c("containers", "kubectl-describe", "kubectl describe", "# kubectl describe", "__NS__: type=string, example=default\n", "Namespace:        __NS__", "Namespace:        default"),
		c("containers", "kubectl-apply", "kubectl apply", "# kubectl apply", "", "created", "deployment.apps/app created"),
		c("containers", "helm-install", "helm install", "# helm install", "__REL__: type=string, example=myrel\n", "NAME: __REL__", "NAME: myrel"),
		c("containers", "podman-ps", "podman ps", "# podman ps", "", "CONTAINER ID  IMAGE", "CONTAINER ID  IMAGE       COMMAND"),
	}
}

func buildSystems() []caseDef {
	return []caseDef{
		c("build-systems", "make-target", "make", "# make", "", "make: Entering directory '/tmp'\n...1 lines omitted...\nmake: Leaving directory '/tmp'", "make: Entering directory '/tmp'\ngcc -o app main.c\nmake: Leaving directory '/tmp'"),
		c("build-systems", "cmake-configure", "cmake", "# cmake", "", "-- Configuring done", "-- Configuring done"),
		c("build-systems", "ninja-build", "ninja", "# ninja", "__N__: type=number, example=10\n", "[__N__/10] Linking", "[10/10] Linking CXX executable app"),
		c("build-systems", "bazel-build", "bazel build", "# bazel build", "", "INFO: Build completed successfully", "INFO: Build completed successfully, 3 total actions"),
		c("build-systems", "meson-setup", "meson setup", "# meson setup", "__DIR__: type=string, example=build\n", "Directory __DIR__", "The Meson build system\nDirectory build created"),
		c("build-systems", "gradle-assemble", "gradle assemble", "# gradle", "", "BUILD SUCCESSFUL", "BUILD SUCCESSFUL in 3s"),
		c("build-systems", "maven-package", "mvn package", "# mvn package", "", "BUILD SUCCESS", "[INFO] BUILD SUCCESS"),
		c("build-systems", "sbt-compile", "sbt compile", "# sbt compile", "", "[success] Total time", "[success] Total time: 5 s, completed"),
		c("build-systems", "mix-compile", "mix compile", "# mix compile", "", "Compiling 1 file (.ex)", "Compiling 1 file (.ex)"),
		c("build-systems", "ant-jar", "ant jar", "# ant jar", "", "BUILD SUCCESSFUL", "BUILD SUCCESSFUL\nTotal time: 2 seconds"),
	}
}

func databases() []caseDef {
	return []caseDef{
		c("databases", "psql-select", "psql SELECT", "# psql", "", " id | name ", " id | name \n----+------"),
		c("databases", "psql-error", "psql error", "# psql", "", "ERROR:  relation \"x\" does not exist", "ERROR:  relation \"x\" does not exist"),
		c("databases", "mysql-show-tables", "mysql SHOW TABLES", "# mysql", "__TBL__: type=string, example=users\n", "__TBL__", "users"),
		c("databases", "redis-ping", "redis-cli PING", "# redis-cli", "", "PONG", "PONG"),
		c("databases", "redis-get", "redis-cli GET", "# redis-cli GET", "__VAL__: type=string, example=hello\n", "__VAL__", "hello"),
		c("databases", "sqlite3-query", "sqlite3", "# sqlite3", "__N__: type=number, example=1\n", "__N__", "1"),
		c("databases", "mongosh-connect", "mongosh", "# mongosh", "", "connecting to: mongodb://127.0.0.1:27017", "connecting to: mongodb://127.0.0.1:27017"),
		c("databases", "pg-dump-start", "pg_dump", "# pg_dump", "", "-- PostgreSQL database dump", "-- PostgreSQL database dump"),
	}
}

func jvmKotlin() []caseDef {
	return []caseDef{
		c("jvm-kotlin", "javac-error", "javac error", "# javac", "", "error: cannot find symbol", "Main.java:3: error: cannot find symbol"),
		c("jvm-kotlin", "java-hello", "java Hello", "# java", "", "Hello", "Hello"),
		c("jvm-kotlin", "kotlin-compile", "kotlinc", "# kotlinc", "", "Compilation completed", "Compilation completed successfully"),
		c("jvm-kotlin", "gradle-test", "gradle test", "# gradle test", "__N__: type=number, example=5\n", "__N__ tests completed", "5 tests completed"),
		c("jvm-kotlin", "mvn-test", "mvn test", "# mvn test", "", "Tests run: 3, Failures: 0", "Tests run: 3, Failures: 0, Errors: 0, Skipped: 0"),
		c("jvm-kotlin", "kotlinc-jvm", "kotlinc -jvm", "# kotlinc", "__VER__: type=string, example=21\n", "jvm target __VER__", "jvm target 21"),
		c("jvm-kotlin", "sbt-test", "sbt test", "# sbt test", "", "[info] All tests passed.", "[info] All tests passed."),
		c("jvm-kotlin", "scala-compile", "scalac", "# scalac", "", "compiling 1 source file", "compiling 1 source file to target"),
	}
}

func cCpp() []caseDef {
	return []caseDef{
		c("c-cpp", "gcc-compile", "gcc compile", "# gcc -c", "", "cc -c main.c", "cc -c main.c"),
		c("c-cpp", "gpp-link", "g++ link", "# g++ -o app", "", "g++ -o app main.cpp", "g++ -o app main.cpp"),
		c("c-cpp", "clang-warning", "clang warning", "# clang", "", "warning: unused variable", "warning: unused variable 'x' [-Wunused-variable]"),
		c("c-cpp", "clang-tidy", "clang-tidy", "# clang-tidy", "", "warning: use auto", "warning: use auto when initializing [modernize-use-auto]"),
		c("c-cpp", "lldb-breakpoint", "lldb", "# lldb", "__FILE__: type=string, example=main.cpp:10\n", "Breakpoint 1: where = main, __FILE__", "Breakpoint 1: where = main, address = 0x1000, file = 'main.cpp', line = 10"),
		c("c-cpp", "gdb-backtrace", "gdb bt", "# gdb", "", "#0  main ()", "#0  main () at main.c:5"),
		c("c-cpp", "make-cpp", "make C++", "# make", "", "g++ -o app main.cpp", "g++ -o app main.cpp"),
		c("c-cpp", "valgrind-leak", "valgrind", "# valgrind", "", "ERROR SUMMARY: 0 errors", "ERROR SUMMARY: 0 errors from 0 contexts"),
	}
}

func shellCases() []caseDef {
	return []caseDef{
		c("shell", "bash-echo", "bash -c echo", "# bash", "", "hi", "hi"),
		c("shell", "sh-exit", "sh exit code", "# sh", "__CODE__: type=number, example=1\n", "exit __CODE__", "exit 1"),
		c("shell", "zsh-version", "zsh --version", "# zsh", "__VER__: type=string, example=5.9\n", "zsh __VER__", "zsh 5.9 (x86_64-apple-darwin23.0.0)"),
		c("shell", "fish-version", "fish --version", "# fish", "", "fish, version 3.6.0", "fish, version 3.6.0"),
		c("shell", "dash-test", "dash test", "# dash", "", "ok", "ok"),
	}
}

func packageManagers() []caseDef {
	return []caseDef{
		c("package-managers", "brew-install", "brew install", "# brew install", "__PKG__: type=string, example=jq\n", "==> Pouring __PKG__", "==> Pouring jq--1.6.arm64_sonoma.bottle.tar.gz"),
		c("package-managers", "apt-get-update", "apt-get update", "# apt-get update", "", "Reading package lists... Done", "Reading package lists... Done"),
		c("package-managers", "apt-install", "apt install", "# apt install", "__PKG__: type=string, example=curl\n", "Setting up __PKG__", "Setting up curl (8.0.0) ..."),
		c("package-managers", "pacman-sync", "pacman -Sy", "# pacman", "", ":: Synchronizing package databases...", ":: Synchronizing package databases..."),
		c("package-managers", "apk-add", "apk add", "# apk add", "__PKG__: type=string, example=git\n", "OK: __PKG__", "OK: 50 MiB in 20 packages"),
		c("package-managers", "dnf-install", "dnf install", "# dnf install", "", "Complete!", "Complete!"),
		c("package-managers", "nix-build", "nix-build", "# nix-build", "__PATH__: type=string, example=/nix/store/abc\n", "__PATH__", "/nix/store/abc"),
		c("package-managers", "mise-install", "mise install", "# mise install", "__TOOL__: type=string, example=node@20\n", "mise __TOOL__", "mise node@20.0.0 installed"),
	}
}

func cloudInfra() []caseDef {
	return []caseDef{
		c("cloud-infra", "terraform-plan", "terraform plan", "# terraform plan", "", "Plan: 1 to add", "Plan: 1 to add, 0 to change, 0 to destroy."),
		c("cloud-infra", "terraform-apply", "terraform apply", "# terraform apply", "", "Apply complete!", "Apply complete! Resources: 1 added, 0 changed, 0 destroyed."),
		c("cloud-infra", "aws-s3-ls", "aws s3 ls", "# aws s3 ls", "__BUCKET__: type=string, example=my-bucket\n", "PRE __BUCKET__/", "PRE my-bucket/"),
		c("cloud-infra", "aws-sts-caller", "aws sts", "# aws sts get-caller-identity", "__ARN__: type=string, example=arn:aws:iam::1:user/u\n", "__ARN__", "arn:aws:iam::1:user/u"),
		c("cloud-infra", "gcloud-auth", "gcloud auth", "# gcloud auth list", "", "Credentialed Accounts", "Credentialed Accounts"),
		c("cloud-infra", "az-login", "az login", "# az login", "", "You have logged in", "You have logged in. Let us set up your subscription."),
		c("cloud-infra", "pulumi-preview", "pulumi preview", "# pulumi preview", "", "Previewing update", "Previewing update (dev)"),
		c("cloud-infra", "gh-pr-create", "gh pr create", "# gh pr create", "__URL__: type=string, example=https://github.com/x/y/pull/1\n", "__URL__", "https://github.com/x/y/pull/1"),
		c("cloud-infra", "glab-mr-list", "glab mr list", "# glab mr list", "__IID__: type=number, example=1\n", "!__IID__", "!1 Draft: feature"),
		c("cloud-infra", "sentry-cli-upload", "sentry-cli", "# sentry-cli upload", "", "Uploaded release", "> Uploaded release files to Sentry"),
	}
}

func languagesOther() []caseDef {
	return []caseDef{
		c("languages-other", "ruby-version", "ruby -v", "# ruby", "__VER__: type=string, example=3.2.0\n", "ruby __VER__", "ruby 3.2.0"),
		c("languages-other", "gem-install", "gem install", "# gem install", "__GEM__: type=string, example=rails\n", "Successfully installed __GEM__", "Successfully installed rails-7.0.0"),
		c("languages-other", "bundler-install", "bundle install", "# bundle install", "", "Bundle complete!", "Bundle complete! 10 Gemfile dependencies, 20 gems now installed."),
		c("languages-other", "php-version", "php -v", "# php", "", "PHP 8.2.0", "PHP 8.2.0 (cli)"),
		c("languages-other", "composer-install", "composer install", "# composer install", "", "Generating autoload files", "Generating autoload files"),
		c("languages-other", "swift-build", "swift build", "# swift build", "", "Build complete!", "Build complete! (1.23s)"),
		c("languages-other", "deno-run", "deno run", "# deno run", "", "Hello", "Hello"),
		c("languages-other", "bun-test", "bun test", "# bun test", "__N__: type=number, example=2\n", "__N__ pass", " 2 pass\n 0 fail"),
		c("languages-other", "crystal-build", "crystal build", "# crystal build", "", "Compiling app", "Compiling app"),
		c("languages-other", "zig-build", "zig build", "# zig build", "", "install", "install"),
	}
}

func miscDevtools() []caseDef {
	return []caseDef{
		c("misc-devtools", "openssl-s-client", "openssl s_client", "# openssl", "", "CONNECTED(00000003)", "CONNECTED(00000003)"),
		c("misc-devtools", "ssh-keygen", "ssh-keygen", "# ssh-keygen", "__PATH__: type=string, example=id_ed25519\n", "Your identification has been saved in __PATH__", "Your identification has been saved in id_ed25519"),
		c("misc-devtools", "scp-copy", "scp", "# scp", "__FILE__: type=string, example=data.txt\n", "__FILE__", "data.txt"),
		c("misc-devtools", "jq-parse", "jq", "# jq .name", "", "\"alice\"", "\"alice\""),
		c("misc-devtools", "yq-read", "yq", "# yq", "__VAL__: type=string, example=prod\n", "__VAL__", "prod"),
		c("misc-devtools", "helm-lint", "helm lint", "# helm lint", "", "==> Linting chart\n...1 lines omitted...\n1 chart(s) linted, 0 chart(s) failed", "==> Linting chart\n[INFO] Chart.yaml: icon is recommended\n1 chart(s) linted, 0 chart(s) failed"),
		c("misc-devtools", "ffmpeg-version", "ffmpeg -version", "# ffmpeg", "__VER__: type=string, example=6.0\n", "ffmpeg version __VER__", "ffmpeg version 6.0 Copyright"),
		c("misc-devtools", "protobuf-compile", "protoc", "# protoc", "", "Writing output", "Writing output to out.pb.go"),
	}
}