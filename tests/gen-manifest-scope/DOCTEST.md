# Gen manifest scope — ledger, prune, warm, and scope isolation

## Version
0.0.3

# DSN (Domain Specific Notion)

### Participants

- **Gen root** — isolated `--gen-dir` for one fixture module.
- **doctest.gen-manifest** — managed path→hash ledger (not a full FS snapshot).
- **Desired set** — paths emitted this run; orphans = in-scope ∧ in-manifest ∧ ¬desired.
- **gen-plan result** — `# new` / `# modified` / `# unchanged` / `# deleted` + summary.

### Behaviors (locked by leaves)

1. **Subset** does not shrink sibling ledger / packages.
2. **Warm same args** → managed `deleted=0`.
3. **Unmanaged plant** survives; not in ledger; not `# deleted`.
4. **Managed file missing** → rewrite; result `# new`.
5. **Out-of-scope missing** not healed by subset; heal when that tree is run.
6. **Source-driven content change** → `# modified` on re-emit.
7. **Source leaf removed** (single-tree) → gen path `# deleted` + leave manifest.
8. **`-a`** wipes gen root → cold-like many `# new`.
9. **Bookkeeping `go.mod` removed** → recreated on next generate.

## Decision Tree

```
gen-manifest-scope/
├── full-then-subset-keeps-sibling/
├── warm-same-args-no-delete/
├── unmanaged-plant-not-deleted/
├── managed-missing-rewrites-as-new/
├── out-of-scope-missing-not-healed/
├── source-change-modified/
├── source-leaf-removed-deleted/
├── force-a-rewrites-as-new/
└── bookkeeping-gomod-recreated/
```

## How to Run

```sh
doctest test --label-all -count=1 ./tests/gen-manifest-scope
```

```go
import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

const genManifestName = "doctest.gen-manifest"

// Request drives 1–3 product CLI generate phases with optional between-run ops.
type Request struct {
	Bin     string
	WorkDir string
	GenDir  string
	TreeA   string
	TreeB   string
	Env     []string
	Timeout time.Duration

	// ArgsFull / ArgsSubset / ArgsThird: CLI argv after binary (include "test").
	ArgsFull   []string
	ArgsSubset []string // run2; default ArgsFull
	ArgsThird  []string // optional run3

	DebugEnv string // empty → bypass-go-test=1

	// Between run1 and run2:
	PlantRel          string   // write unmanaged file under GenDir
	DeleteGenRels     []string // remove managed (or bookkeeping) paths under GenDir
	RewriteSourceRels map[string]string // workdir-rel path → new file contents
	RemoveSourceRels  []string // remove source paths under WorkDir (dirs ok)
}

// Response captures up to three runs and post-op disk/ledger state.
type Response struct {
	FullExitCode   int
	FullStdout     string
	FullStderr     string
	FullErr        error
	SubsetExitCode int
	SubsetStdout   string
	SubsetStderr   string
	SubsetErr      error
	ThirdExitCode  int
	ThirdStdout    string
	ThirdStderr    string
	ThirdErr       error

	ManifestAfterFull   string
	ManifestAfterSubset string
	ManifestAfterThird  string
	SiblingGenDirExists bool

	PlantAbs             string
	PlantExistsAfter     bool
	PlantInManifestAfter bool

	// First DeleteGenRels path abs + existence after run2 / run3.
	DeletedGenAbs          string
	DeletedGenExistsAfter2 bool
	DeletedGenExistsAfter3 bool
	GoModExistsAfter2      bool
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	if req.Timeout <= 0 {
		req.Timeout = 2 * time.Minute
	}
	if req.Bin == "" {
		return nil, fmt.Errorf("req.Bin required")
	}
	if req.GenDir == "" || req.WorkDir == "" {
		return nil, fmt.Errorf("GenDir and WorkDir required")
	}
	if len(req.ArgsFull) == 0 {
		return nil, fmt.Errorf("ArgsFull required")
	}
	args2 := req.ArgsSubset
	if len(args2) == 0 {
		args2 = req.ArgsFull
	}

	resp := &Response{}

	r1, err := runProduct(t, req, req.ArgsFull)
	if err != nil && r1 == nil {
		return nil, err
	}
	if r1 != nil {
		resp.FullExitCode = r1.exitCode
		resp.FullStdout = r1.stdout
		resp.FullStderr = r1.stderr
		resp.FullErr = r1.err
	}
	resp.ManifestAfterFull = readFileOrEmpty(filepath.Join(req.GenDir, genManifestName))

	if err := applyBetweenRun1And2(t, req, resp); err != nil {
		return resp, err
	}

	r2, err2 := runProduct(t, req, args2)
	if err2 != nil && r2 == nil {
		return resp, err2
	}
	if r2 != nil {
		resp.SubsetExitCode = r2.exitCode
		resp.SubsetStdout = r2.stdout
		resp.SubsetStderr = r2.stderr
		resp.SubsetErr = r2.err
	}
	resp.ManifestAfterSubset = readFileOrEmpty(filepath.Join(req.GenDir, genManifestName))
	fillAfterRun2(req, resp)

	if len(req.ArgsThird) > 0 {
		r3, err3 := runProduct(t, req, req.ArgsThird)
		if err3 != nil && r3 == nil {
			return resp, err3
		}
		if r3 != nil {
			resp.ThirdExitCode = r3.exitCode
			resp.ThirdStdout = r3.stdout
			resp.ThirdStderr = r3.stderr
			resp.ThirdErr = r3.err
		}
		resp.ManifestAfterThird = readFileOrEmpty(filepath.Join(req.GenDir, genManifestName))
		if resp.DeletedGenAbs != "" {
			if _, err := os.Stat(resp.DeletedGenAbs); err == nil {
				resp.DeletedGenExistsAfter3 = true
			}
		}
	}
	return resp, nil
}

func applyBetweenRun1And2(t *testing.T, req *Request, resp *Response) error {
	t.Helper()
	for _, rel := range req.DeleteGenRels {
		rel = filepath.ToSlash(rel)
		abs := filepath.Join(req.GenDir, filepath.FromSlash(rel))
		if resp.DeletedGenAbs == "" {
			resp.DeletedGenAbs = abs
		}
		if err := os.Remove(abs); err != nil && !os.IsNotExist(err) {
			// Also allow removing dirs for future; files first.
			if err2 := os.RemoveAll(abs); err2 != nil {
				return fmt.Errorf("delete gen %s: %w", rel, err)
			}
		}
	}
	if req.PlantRel != "" {
		rel := filepath.ToSlash(req.PlantRel)
		abs := filepath.Join(req.GenDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return fmt.Errorf("plant mkdir: %w", err)
		}
		if err := os.WriteFile(abs, []byte("package unused\n"), 0o644); err != nil {
			return fmt.Errorf("plant write: %w", err)
		}
		resp.PlantAbs = abs
	}
	for rel, content := range req.RewriteSourceRels {
		abs := filepath.Join(req.WorkDir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
			return fmt.Errorf("rewrite mkdir %s: %w", rel, err)
		}
		if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
			return fmt.Errorf("rewrite %s: %w", rel, err)
		}
	}
	for _, rel := range req.RemoveSourceRels {
		abs := filepath.Join(req.WorkDir, filepath.FromSlash(rel))
		if err := os.RemoveAll(abs); err != nil {
			return fmt.Errorf("remove source %s: %w", rel, err)
		}
	}
	return nil
}

func fillAfterRun2(req *Request, resp *Response) {
	sib := filepath.Join(req.GenDir, "tree-b")
	if st, err := os.Stat(sib); err == nil && st.IsDir() {
		resp.SiblingGenDirExists = true
	}
	if resp.PlantAbs != "" {
		if _, err := os.Stat(resp.PlantAbs); err == nil {
			resp.PlantExistsAfter = true
		}
		resp.PlantInManifestAfter = manifestHasExactRel(resp.ManifestAfterSubset, filepath.ToSlash(req.PlantRel))
	}
	if resp.DeletedGenAbs != "" {
		if _, err := os.Stat(resp.DeletedGenAbs); err == nil {
			resp.DeletedGenExistsAfter2 = true
		}
	}
	if _, err := os.Stat(filepath.Join(req.GenDir, "go.mod")); err == nil {
		resp.GoModExistsAfter2 = true
	}
}

type cliResult struct {
	exitCode int
	stdout   string
	stderr   string
	err      error
}

func runProduct(t *testing.T, req *Request, args []string) (*cliResult, error) {
	t.Helper()
	ctx, cancel := context.WithTimeout(context.Background(), req.Timeout)
	defer cancel()
	cmd := exec.CommandContext(ctx, req.Bin, args...)
	cmd.Dir = req.WorkDir
	base := filterEnv(os.Environ(), "DOCTEST_DEBUG")
	env := append(append([]string{}, base...), req.Env...)
	dbg := req.DebugEnv
	if dbg == "" {
		dbg = "bypass-go-test=1"
	}
	env = append(env, "DOCTEST_DEBUG="+dbg)
	cmd.Env = env
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	out := &cliResult{stdout: stdout.String(), stderr: stderr.String(), err: err}
	if err == nil {
		return out, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		out.exitCode = exitErr.ExitCode()
		return out, nil
	}
	if ctx.Err() != nil {
		return out, ctx.Err()
	}
	return out, err
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

func readFileOrEmpty(path string) string {
	b, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return string(b)
}

func manifestHasTreePrefix(manifest, treeRel string) bool {
	prefix := strings.Trim(filepath.ToSlash(treeRel), "/") + "/"
	for _, line := range strings.Split(manifest, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "version") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) < 1 {
			continue
		}
		if strings.HasPrefix(fields[0], prefix) || fields[0] == strings.TrimSuffix(prefix, "/") {
			return true
		}
	}
	return false
}

func manifestHasExactRel(manifest, rel string) bool {
	rel = filepath.ToSlash(rel)
	for _, line := range strings.Split(manifest, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "version") {
			continue
		}
		fields := strings.Fields(line)
		if len(fields) >= 1 && fields[0] == rel {
			return true
		}
	}
	return false
}

func parseSummaryCount(stderr, key string) int {
	re := regexp.MustCompile(key + `\s*=\s*(\d+)`)
	m := re.FindStringSubmatch(stderr)
	if m == nil {
		return -1
	}
	n, _ := strconv.Atoi(m[1])
	return n
}

func parseDeletedFromGenPlan(stderr string) int {
	return parseSummaryCount(stderr, "deleted")
}

func parseNewFromGenPlan(stderr string) int {
	return parseSummaryCount(stderr, "new")
}

func parseModifiedFromGenPlan(stderr string) int {
	return parseSummaryCount(stderr, "modified")
}

// genPlanHasTag reports a result line containing basename and "# <tag>".
func genPlanHasTag(stderr, basename, tag string) bool {
	needle := "# " + tag
	for _, line := range strings.Split(stderr, "\n") {
		if strings.Contains(line, basename) && strings.Contains(line, needle) {
			return true
		}
	}
	return false
}

func requireExit0(t *testing.T, label string, code int, stderr string) {
	t.Helper()
	if code != 0 {
		t.Fatalf("%s exit=%d\nstderr:\n%s", label, code, stderr)
	}
}
```
