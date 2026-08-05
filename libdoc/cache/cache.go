// Package cache implements `doctest cache [--clean] [--dry-run]`: inventory and
// wipe of durable doctest cache roots under $CacheHome/doctest (plus override
// roots that live outside that tree).
//
// Roots are injectable via Options so L2 tests can pass t.TempDir paths without
// process Setenv. Production CLI resolves from DOCTEST_CACHE_HOME,
// DOCTEST_LEAF_CACHE, and DOCTEST_METRICS_ROOT via OptionsFromEnv.
package cache

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/leafcache"
	"github.com/xhd2015/doctest/libdoc/metrics"
)

// Options holds injectable cache roots for Info/Clean.
// Empty fields are not resolved from env — use OptionsFromEnv for production.
type Options struct {
	// CacheHome is the base cache directory (product: DOCTEST_CACHE_HOME or UserCacheDir).
	// Doctest root is always CacheHome/doctest.
	CacheHome string

	// LeafCache is an optional absolute leaf-cache store root (product: DOCTEST_LEAF_CACHE).
	// Included in clean plan when set and not under DoctestRoot.
	LeafCache string

	// MetricsRoot is an optional metrics base (product: DOCTEST_METRICS_ROOT).
	// When set and the metrics tree lies outside DoctestRoot, that tree is an
	// extra clean target ($MetricsRoot/doctest/metrics).
	MetricsRoot string
}

// OptionsFromEnv resolves production roots from environment variables.
func OptionsFromEnv() (Options, error) {
	home, err := core.CacheHome()
	if err != nil {
		return Options{}, err
	}
	opts := Options{CacheHome: home}
	if v := strings.TrimSpace(os.Getenv(leafcache.EnvLeafCache)); v != "" {
		abs, err := filepath.Abs(v)
		if err != nil {
			return Options{}, err
		}
		opts.LeafCache = abs
	}
	if v := strings.TrimSpace(os.Getenv(metrics.EnvMetricsRoot)); v != "" {
		abs, err := filepath.Abs(v)
		if err != nil {
			return Options{}, err
		}
		opts.MetricsRoot = abs
	}
	return opts, nil
}

// DoctestRoot returns $CacheHome/doctest.
func DoctestRoot(cacheHome string) string {
	return filepath.Join(cacheHome, "doctest")
}

// Bucket is a first-level subdirectory under the doctest root with its size.
type Bucket struct {
	Name string
	Size int64
}

// Scan lists first-level buckets under doctest root and their sizes.
// Missing root yields empty buckets and zero total (not an error).
func Scan(doctestRoot string) (buckets []Bucket, total int64, err error) {
	ents, err := os.ReadDir(doctestRoot)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	for _, e := range ents {
		if !e.IsDir() {
			// Count loose files toward total but not as named buckets.
			info, err := e.Info()
			if err != nil {
				return nil, 0, err
			}
			total += info.Size()
			continue
		}
		name := e.Name()
		sz, err := DirSize(filepath.Join(doctestRoot, name))
		if err != nil {
			return nil, 0, err
		}
		buckets = append(buckets, Bucket{Name: name, Size: sz})
		total += sz
	}
	sort.Slice(buckets, func(i, j int) bool {
		return buckets[i].Name < buckets[j].Name
	})
	return buckets, total, nil
}

// DirSize walks path and sums file sizes. Missing path is size 0.
func DirSize(path string) (int64, error) {
	var total int64
	err := filepath.Walk(path, func(_ string, info os.FileInfo, err error) error {
		if err != nil {
			if os.IsNotExist(err) {
				return nil
			}
			return err
		}
		if info == nil || info.IsDir() {
			return nil
		}
		total += info.Size()
		return nil
	})
	if err != nil && os.IsNotExist(err) {
		return 0, nil
	}
	return total, err
}

// FormatSize renders n bytes with human units B/K/M/G.
func FormatSize(n int64) string {
	if n < 0 {
		n = 0
	}
	const (
		k = 1024
		m = 1024 * 1024
		g = 1024 * 1024 * 1024
	)
	switch {
	case n < k:
		return fmt.Sprintf("%dB", n)
	case n < m:
		return formatUnit(float64(n)/float64(k), "K")
	case n < g:
		return formatUnit(float64(n)/float64(m), "M")
	default:
		return formatUnit(float64(n)/float64(g), "G")
	}
}

func formatUnit(v float64, unit string) string {
	if v < 10 {
		return fmt.Sprintf("%.1f%s", v, unit)
	}
	return fmt.Sprintf("%.0f%s", v, unit)
}

// Info writes a human summary of cache home, doctest root, buckets, and total.
func Info(w io.Writer, opts Options) error {
	opts, err := normalizeOpts(opts)
	if err != nil {
		return err
	}
	root := DoctestRoot(opts.CacheHome)
	buckets, total, err := Scan(root)
	if err != nil {
		return err
	}
	fmt.Fprintf(w, "Cache home:   %s\n", opts.CacheHome)
	fmt.Fprintf(w, "Doctest root: %s\n", root)
	fmt.Fprintln(w)
	if len(buckets) == 0 {
		fmt.Fprintf(w, "  (empty — no cache buckets)\n")
		fmt.Fprintln(w)
		fmt.Fprintf(w, "Total: %s  (0 buckets)\n", FormatSize(total))
		return nil
	}
	// Align names in a simple two-column layout.
	nameWidth := 0
	for _, b := range buckets {
		if len(b.Name) > nameWidth {
			nameWidth = len(b.Name)
		}
	}
	if nameWidth < 12 {
		nameWidth = 12
	}
	for _, b := range buckets {
		fmt.Fprintf(w, "  %-*s  %s\n", nameWidth, b.Name, FormatSize(b.Size))
	}
	fmt.Fprintln(w)
	fmt.Fprintf(w, "Total: %s  (%d buckets)\n", FormatSize(total), len(buckets))
	return nil
}

// CleanTarget is one path planned for removal.
type CleanTarget struct {
	Path string
	Size int64
}

// PlanClean returns paths to remove: main doctest root plus outside overrides.
// Main root is always listed (even if missing) so dry-run/live messaging is stable.
// Missing main root still appears; size is 0 when absent.
func PlanClean(opts Options) ([]CleanTarget, error) {
	opts, err := normalizeOpts(opts)
	if err != nil {
		return nil, err
	}
	root := DoctestRoot(opts.CacheHome)
	if err := checkSafeMainRoot(root); err != nil {
		return nil, err
	}
	sz, err := DirSize(root)
	if err != nil {
		return nil, err
	}
	targets := []CleanTarget{{Path: root, Size: sz}}
	extras, err := outsideOverrides(opts, root)
	if err != nil {
		return nil, err
	}
	for _, p := range extras {
		if err := checkSafeOverride(p); err != nil {
			return nil, err
		}
		esz, err := DirSize(p)
		if err != nil {
			return nil, err
		}
		targets = append(targets, CleanTarget{Path: p, Size: esz})
	}
	return targets, nil
}

// Clean removes (or dry-runs) the planned clean targets.
// On live remove error, returns non-zero-style error including the path (hard-fail).
func Clean(w io.Writer, opts Options, dryRun bool) error {
	targets, err := PlanClean(opts)
	if err != nil {
		return err
	}
	for _, t := range targets {
		human := FormatSize(t.Size)
		if dryRun {
			fmt.Fprintf(w, "[dry-run] would remove: %s  (%s)\n", t.Path, human)
			continue
		}
		if err := os.RemoveAll(t.Path); err != nil {
			return fmt.Errorf("remove %s: %w", t.Path, err)
		}
		fmt.Fprintf(w, "Removed %s  (%s)\n", t.Path, human)
	}
	return nil
}

// Run executes the cache command for args after "cache" (flags only).
// Writers receive info/clean output on stdout; flag/help errors return as error
// (CLI prints them to stderr via process main / harness).
func Run(stdout io.Writer, args []string, opts Options) error {
	clean, dryRun, err := parseFlags(args)
	if err != nil {
		return err
	}
	if dryRun && !clean {
		return fmt.Errorf("--dry-run requires --clean")
	}
	if clean {
		return Clean(stdout, opts, dryRun)
	}
	return Info(stdout, opts)
}

func parseFlags(args []string) (clean, dryRun bool, err error) {
	for _, a := range args {
		switch a {
		case "--clean":
			clean = true
		case "--dry-run":
			dryRun = true
		case "-h", "--help":
			// Caller should intercept help; treat as unexpected here.
			return false, false, fmt.Errorf("unexpected help flag")
		default:
			if strings.HasPrefix(a, "-") {
				return false, false, fmt.Errorf("unknown flag: %s", a)
			}
			return false, false, fmt.Errorf("unexpected argument: %s", a)
		}
	}
	return clean, dryRun, nil
}

func normalizeOpts(opts Options) (Options, error) {
	if strings.TrimSpace(opts.CacheHome) == "" {
		return Options{}, fmt.Errorf("cache home is empty")
	}
	abs, err := filepath.Abs(opts.CacheHome)
	if err != nil {
		return Options{}, err
	}
	opts.CacheHome = abs
	if opts.LeafCache != "" {
		lc, err := filepath.Abs(opts.LeafCache)
		if err != nil {
			return Options{}, err
		}
		opts.LeafCache = lc
	}
	if opts.MetricsRoot != "" {
		mr, err := filepath.Abs(opts.MetricsRoot)
		if err != nil {
			return Options{}, err
		}
		opts.MetricsRoot = mr
	}
	return opts, nil
}

// outsideOverrides returns clean targets not covered by deleting doctest root.
func outsideOverrides(opts Options, doctestRoot string) ([]string, error) {
	var out []string
	seen := map[string]bool{filepath.Clean(doctestRoot): true}

	add := func(p string) {
		p = filepath.Clean(p)
		if p == "" || seen[p] {
			return
		}
		if isUnderOrEqual(p, doctestRoot) {
			return
		}
		seen[p] = true
		out = append(out, p)
	}

	if opts.LeafCache != "" {
		add(opts.LeafCache)
	}
	if opts.MetricsRoot != "" {
		// Product metrics tree: $MetricsRoot/doctest/metrics
		mt := filepath.Join(opts.MetricsRoot, "doctest", "metrics")
		add(mt)
	}
	return out, nil
}

func isUnderOrEqual(path, root string) bool {
	path = filepath.Clean(path)
	root = filepath.Clean(root)
	if path == root {
		return true
	}
	sep := string(filepath.Separator)
	return strings.HasPrefix(path, root+sep)
}

// checkSafeMainRoot refuses dangerous remove targets for the main doctest tree.
func checkSafeMainRoot(path string) error {
	if err := checkSafePath(path); err != nil {
		return err
	}
	if filepath.Base(path) != "doctest" {
		return fmt.Errorf("refuse remove: main cache root base name must be doctest (got %q)", filepath.Base(path))
	}
	return nil
}

func checkSafeOverride(path string) error {
	return checkSafePath(path)
}

func checkSafePath(path string) error {
	if strings.TrimSpace(path) == "" {
		return fmt.Errorf("refuse remove: empty path")
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("refuse remove: path is not absolute: %s", path)
	}
	clean := filepath.Clean(path)
	if clean == string(filepath.Separator) || clean == filepath.VolumeName(clean)+string(filepath.Separator) {
		return fmt.Errorf("refuse remove: filesystem root")
	}
	return nil
}
