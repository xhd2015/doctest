package build

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
)

// GenPlanArg describes one CLI path argument for gen-plan printing.
type GenPlanArg struct {
	// Arg is the original CLI argument string.
	Arg string
	// TreeRel is the slash path of the doctest tree relative to the module root
	// (prefix of package paths under the gen root). Empty / "." means packages
	// live at the gen root.
	TreeRel string
}

// PrintGenPlanBanner writes the short gen-plan debug banner to w.
func PrintGenPlanBanner(w io.Writer) {
	if w == nil {
		w = os.Stderr
	}
	fmt.Fprintln(w, "doctest: DOCTEST_DEBUG gen-plan=1")
}

// PrintGenPlanInvocation writes the invocation header.
func PrintGenPlanInvocation(w io.Writer, args []string, opts core.Options, genRoot string, multi bool) {
	if w == nil {
		w = os.Stderr
	}
	mode := "single-tree"
	if multi {
		mode = "multi-arg"
	}
	fmt.Fprintln(w, "gen-plan: invocation")
	fmt.Fprintf(w, "  args: %s\n", strings.Join(args, " "))
	switch {
	case opts.LabelAll:
		fmt.Fprintln(w, "  labels: --label-all")
	case len(opts.LabelExprs) > 0:
		fmt.Fprintf(w, "  labels: %s\n", strings.Join(opts.LabelExprs, ", "))
	default:
		fmt.Fprintln(w, "  labels: (none)")
	}
	if genRoot == "" {
		genRoot = opts.GenDir
	}
	if genRoot == "" {
		genRoot = "(default mapping-gen)"
	}
	fmt.Fprintf(w, "  gen-root: %s\n", genRoot)
	fmt.Fprintf(w, "  mode: %s\n", mode)
}

// PrintGenPlanAndResult prints plan hierarchies then the result tree with
// statuses. multi selects per-arg package-only + merged vs single-arg full tree.
func PrintGenPlanAndResult(w io.Writer, opts core.Options, args []GenPlanArg, genRoot string, multi bool) {
	if w == nil {
		w = os.Stderr
	}
	if opts.GenBatch == nil || genRoot == "" {
		return
	}
	desired := opts.GenBatch.Desired(genRoot)
	outcomes := opts.GenBatch.Outcomes(genRoot)
	if len(desired) == 0 && len(outcomes) == 0 {
		// Still print markers so tests can detect gen-plan emit path.
		n := len(args)
		if n == 0 {
			n = 1
		}
		for i, a := range args {
			fmt.Fprintf(w, "gen-plan: arg[%d/%d]  %s\n", i+1, n, a.Arg)
		}
		if multi {
			fmt.Fprintln(w, "gen-plan: merged")
		}
		fmt.Fprintln(w, "gen-plan: result")
		fmt.Fprintln(w, "  summary: new=0 modified=0 unchanged=0 deleted=0")
		return
	}

	// Union path set: desired + deleted outcomes.
	allPaths := map[string]struct{}{}
	for p := range desired {
		allPaths[p] = struct{}{}
	}
	for p, st := range outcomes {
		if st == core.EmitDeleted {
			allPaths[p] = struct{}{}
		}
	}

	n := len(args)
	if n == 0 {
		// Fall back to a single synthetic arg covering everything.
		args = []GenPlanArg{{Arg: ".", TreeRel: "."}}
		n = 1
		multi = false
	}

	// --- plan ---
	for i, a := range args {
		fmt.Fprintf(w, "gen-plan: arg[%d/%d]  %s\n", i+1, n, a.Arg)
		var paths []string
		if multi {
			paths = filterTreePaths(allPaths, a.TreeRel, false)
		} else {
			// Single-arg: bookkeeping + tree packages.
			paths = filterTreePaths(allPaths, a.TreeRel, true)
		}
		// Plan: structure only (no status tags/color).
		printPathTree(w, paths, nil, colorStyle{enabled: false}, false)
	}
	if multi {
		fmt.Fprintln(w, "gen-plan: merged")
		// Full set: bookkeeping + all trees + __workspace.
		paths := sortedPaths(allPaths)
		printPathTree(w, paths, nil, colorStyle{enabled: false}, false)
	}

	// --- result ---
	fmt.Fprintln(w, "gen-plan: result")
	style := newColorStyle(opts.Color, w)
	legend := "  # gray=unchanged  green=new|modified  red=deleted  (tags always printed)"
	if style.enabled {
		fmt.Fprintln(w, style.gray(legend))
	} else {
		fmt.Fprintln(w, legend)
	}
	// Result: annotate + optional color (colorize=true).
	if multi {
		paths := sortedPaths(allPaths)
		printPathTree(w, paths, outcomes, style, true)
	} else {
		paths := filterTreePaths(allPaths, args[0].TreeRel, true)
		printPathTree(w, paths, outcomes, style, true)
	}
	nw, mod, unc, del := countOutcomes(allPaths, outcomes)
	summary := fmt.Sprintf("  summary: new=%d modified=%d unchanged=%d deleted=%d", nw, mod, unc, del)
	if style.enabled {
		summary = style.gray(summary)
	}
	fmt.Fprintln(w, summary)
}

func filterTreePaths(all map[string]struct{}, treeRel string, includeBookkeeping bool) []string {
	treeRel = filepath.ToSlash(filepath.Clean(treeRel))
	if treeRel == "" {
		treeRel = "."
	}
	var out []string
	for p := range all {
		p = filepath.ToSlash(p)
		if rootBookkeepingPath(p) {
			if includeBookkeeping {
				out = append(out, p)
			}
			continue
		}
		if treeRel == "." {
			// Everything except we already handled bookkeeping above when include.
			out = append(out, p)
			continue
		}
		if p == treeRel || strings.HasPrefix(p, treeRel+"/") {
			out = append(out, p)
		}
	}
	sort.Strings(out)
	return out
}

func rootBookkeepingPath(rel string) bool {
	switch rel {
	case "go.mod", "go.sum", "doctest.gen-manifest", "doctest.tidy-done":
		return true
	default:
		return strings.HasPrefix(rel, "doctest.")
	}
}

func sortedPaths(all map[string]struct{}) []string {
	out := make([]string, 0, len(all))
	for p := range all {
		out = append(out, filepath.ToSlash(p))
	}
	sort.Strings(out)
	return out
}

func countOutcomes(all map[string]struct{}, outcomes map[string]core.EmitStatus) (nw, mod, unc, del int) {
	seen := map[string]struct{}{}
	for p := range all {
		seen[p] = struct{}{}
		st := outcomes[p]
		switch st {
		case core.EmitNew:
			nw++
		case core.EmitModified:
			mod++
		case core.EmitDeleted:
			del++
		case core.EmitUnchanged:
			unc++
		default:
			// Desired without outcome: treat as unchanged (present, not rewritten).
			unc++
		}
	}
	for p, st := range outcomes {
		if _, ok := seen[p]; ok {
			continue
		}
		if st == core.EmitDeleted {
			del++
		}
	}
	return nw, mod, unc, del
}

// treeNode is a simple path trie for indented printing.
type treeNode struct {
	name     string
	children map[string]*treeNode
	file     bool // true if this node is a terminal file path component
	rel      string
}

func printPathTree(w io.Writer, paths []string, outcomes map[string]core.EmitStatus, style colorStyle, colorize bool) {
	root := &treeNode{children: map[string]*treeNode{}}
	for _, p := range paths {
		p = strings.Trim(filepath.ToSlash(p), "/")
		if p == "" {
			continue
		}
		parts := strings.Split(p, "/")
		cur := root
		acc := ""
		for i, part := range parts {
			if acc == "" {
				acc = part
			} else {
				acc = acc + "/" + part
			}
			if cur.children == nil {
				cur.children = map[string]*treeNode{}
			}
			ch, ok := cur.children[part]
			if !ok {
				ch = &treeNode{name: part, children: map[string]*treeNode{}, rel: acc}
				cur.children[part] = ch
			}
			if i == len(parts)-1 {
				ch.file = true
				ch.rel = p
			}
			cur = ch
		}
	}
	printTreeNode(w, root, "  ", outcomes, style, colorize)
}

func printTreeNode(w io.Writer, n *treeNode, indent string, outcomes map[string]core.EmitStatus, style colorStyle, colorize bool) {
	if n.children == nil || len(n.children) == 0 {
		return
	}
	names := make([]string, 0, len(n.children))
	for name := range n.children {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		ch := n.children[name]
		label := name
		// Directory vs file: if children exist, show as directory line.
		isDir := len(ch.children) > 0
		if isDir {
			// Print directory name (default/gray).
			if colorize && style.enabled {
				fmt.Fprintln(w, indent+style.gray(label+"/"))
			} else {
				fmt.Fprintln(w, indent+label+"/")
			}
			printTreeNode(w, ch, indent+"  ", outcomes, style, colorize)
			// If the path is also a file (unusual), still ok.
			continue
		}
		// File leaf. Plan phase (colorize=false): plain name only.
		// Result phase (colorize=true): color basename + always "# <tag>".
		if !colorize {
			fmt.Fprintln(w, indent+label)
			continue
		}
		st := core.EmitUnchanged
		if outcomes != nil {
			if s, ok := outcomes[ch.rel]; ok && s != core.EmitUnknown {
				st = s
			}
		}
		base := label
		if style.enabled {
			base = colorStatusLabel(base, st, style)
		}
		tag := "  # " + st.Tag()
		if style.enabled {
			tag = style.gray(tag)
		}
		fmt.Fprintln(w, indent+base+tag)
	}
}

func colorStatusLabel(base string, st core.EmitStatus, style colorStyle) string {
	switch st {
	case core.EmitNew, core.EmitModified:
		return style.green(base)
	case core.EmitDeleted:
		return style.red(base)
	case core.EmitUnchanged:
		return style.gray(base)
	default:
		return style.gray(base)
	}
}

// ResolveGenPlanArgs maps CLI remain args to GenPlanArg using module-relative
// tree paths. workDir is the process cwd for relative args (may be empty).
func ResolveGenPlanArgs(remainArgs []string, workDir string) []GenPlanArg {
	out := make([]GenPlanArg, 0, len(remainArgs))
	for _, a := range remainArgs {
		arg := a
		// Strip ./... suffix for pattern args.
		base := a
		if strings.HasSuffix(base, "/...") {
			base = strings.TrimSuffix(base, "/...")
		}
		if base == "" {
			base = "."
		}
		abs := base
		if !filepath.IsAbs(abs) {
			if workDir != "" {
				abs = filepath.Join(workDir, abs)
			} else if a, err := filepath.Abs(base); err == nil {
				abs = a
			}
		}
		abs = filepath.Clean(abs)
		treeRel := "."
		if modRoot, _, ok := core.FindModuleRoot(abs); ok && modRoot != "" {
			if rel, err := filepath.Rel(modRoot, abs); err == nil {
				treeRel = filepath.ToSlash(rel)
			}
		} else {
			// Fallback: use base name so tree-a/tree-b style still filters.
			treeRel = filepath.ToSlash(filepath.Base(abs))
		}
		if treeRel == "" {
			treeRel = "."
		}
		out = append(out, GenPlanArg{Arg: arg, TreeRel: treeRel})
	}
	return out
}
