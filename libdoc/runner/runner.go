package runner

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/xhd2015/less-flags"
	runnerbuild "github.com/xhd2015/doctest/libdoc/build"
	"github.com/xhd2015/doctest/libdoc/core"
	"github.com/xhd2015/doctest/libdoc/path_resolve"
	"github.com/xhd2015/doctest/libdoc/validate"
)

var ErrNoTestsFound = path_resolve.ErrNoTestsFound

func Build(dir string) error {
	return runnerbuild.Build(dir, core.Options{RemoveTemp: false})
}

func BuildArgs(args []string) error {
	return processArgs(args, "build", parseBuildOptions, func(dir string, opts core.Options) error {
		err := runnerbuild.Build(dir, opts)
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	})
}

func Test(args []string) error {
	opts, remainArgs, err := parseTestOptions(args)
	if err != nil {
		return err
	}
	if len(remainArgs) < 1 {
		return fmt.Errorf("test requires <dir>")
	}

	// One session id for the whole CLI invocation so parallel trees share
	// session.Once / testbin materialization when nested self-tests run.
	if v, ok := syscall.Getenv(core.DoctestSessionIDEnv); !ok || v == "" {
		_ = os.Setenv(core.DoctestSessionIDEnv, core.NewDoctestSessionID())
	}

	opts.SuppressResultSummary = true
	start := time.Now()
	var stats runnerbuild.TestRunStats
	var statsMu sync.Mutex
	runFn := func(dir string, o core.Options) error {
		o.SuppressResultSummary = true
		s, err := runnerbuild.TestWithStats(dir, o)
		statsMu.Lock()
		stats.Passed += s.Passed
		stats.Total += s.Total
		stats.Skipped = append(stats.Skipped, s.Skipped...)
		if s.NoTestsChanged {
			stats.NoTestsChanged = true
		}
		statsMu.Unlock()
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	}

	var runErr error
	if len(remainArgs) == 1 {
		runErr = processSingleArg(remainArgs[0], opts, runFn)
	} else {
		runErr = processMultiArg(remainArgs, opts, runFn)
	}

	if len(stats.Skipped) > 0 {
		runnerbuild.PrintSkippedSummary(stats.Skipped)
	}
	if stats.Total > 0 {
		stats.Elapsed = time.Since(start)
		runnerbuild.PrintResultSummary(opts, stats)
	}

	if runErr != nil {
		if errors.Is(runErr, ErrNoTestsFound) && stats.Total == 0 && len(stats.Skipped) == 0 {
			return ErrNoTestsFound
		}
		if errors.Is(runErr, ErrNoTestsFound) && len(stats.Skipped) > 0 {
			return nil
		}
		return runErr
	}
	if stats.Total == 0 {
		if stats.NoTestsChanged || len(stats.Skipped) > 0 {
			return nil
		}
		return ErrNoTestsFound
	}
	if stats.Passed < stats.Total {
		return fmt.Errorf("%d of %d tests passed", stats.Passed, stats.Total)
	}
	return nil
}

func VetArgs(args []string) error {
	return processArgs(args, "vet", parseVetOptions, func(dir string, opts core.Options) error {
		return validate.RunWithOptions(dir, opts)
	})
}

func processArgs(args []string, cmdName string, parseFn func([]string) (core.Options, []string, error), processDirFn func(string, core.Options) error) error {
	opts, remainArgs, err := parseFn(args)
	if err != nil {
		return err
	}
	if len(remainArgs) < 1 {
		return fmt.Errorf("%s requires <dir>", cmdName)
	}
	if len(remainArgs) == 1 {
		return processSingleArg(remainArgs[0], opts, processDirFn)
	}
	return processMultiArg(remainArgs, opts, processDirFn)
}

func processSingleArg(arg string, opts core.Options, fn func(string, core.Options) error) error {
	if arg == "..." {
		return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
	}
	if path_resolve.IsDotDotDotPattern(arg) {
		// Parallel trees: buffer each tree's streams so progress lines do not interleave.
		var printMu sync.Mutex
		stdoutBase := opts.Stdout
		if stdoutBase == nil {
			stdoutBase = os.Stdout
		}
		stderrBase := opts.Stderr
		if stderrBase == nil {
			stderrBase = os.Stderr
		}
		return path_resolve.RunForDirs(path_resolve.ExtractBasePath(arg), func(dir string) error {
			root, _ := path_resolve.ResolveRoot(dir)
			if root == "" {
				root = dir
			}
			var outBuf, errBuf bytes.Buffer
			o := opts
			o.SubDir = dir
			o.ExplicitLeaf = false
			o.Stdout = &outBuf
			o.Stderr = &errBuf
			err := fn(root, o)
			// stderr first: tree header / "cd ..." then progress on stdout
			printMu.Lock()
			_, _ = io.Copy(stderrBase, &errBuf)
			_, _ = io.Copy(stdoutBase, &outBuf)
			printMu.Unlock()
			return err
		})
	}
	targetDir, explicitLeaf := resolveTestTarget(arg)
	root, ok := path_resolve.ResolveRoot(targetDir)
	if !ok {
		root = targetDir
	}
	opts.SubDir = targetDir
	opts.ExplicitLeaf = explicitLeaf
	return fn(root, opts)
}

func processMultiArg(args []string, opts core.Options, fn func(string, core.Options) error) error {
	var errs []string
	allNoTestsFound := true

	for _, arg := range args {
		if arg == "..." {
			return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
		}
		if path_resolve.IsDotDotDotPattern(arg) {
			err := path_resolve.RunForDirs(path_resolve.ExtractBasePath(arg), func(dir string) error {
				root, _ := path_resolve.ResolveRoot(dir)
				if root == "" {
					root = dir
				}
				o := opts
				o.SubDir = dir
				o.ExplicitLeaf = false
				return fn(root, o)
			})
			if err != nil {
				if errors.Is(err, ErrNoTestsFound) {
					continue
				}
				errs = append(errs, err.Error())
				allNoTestsFound = false
			} else {
				allNoTestsFound = false
			}
			continue
		}
		targetDir, explicitLeaf := resolveTestTarget(arg)
		root, ok := path_resolve.ResolveRoot(targetDir)
		if !ok {
			root = targetDir
		}
		o := opts
		o.SubDir = targetDir
		o.ExplicitLeaf = explicitLeaf
		err := fn(root, o)
		if errors.Is(err, ErrNoTestsFound) {
			continue
		}
		if err != nil {
			errs = append(errs, err.Error())
			allNoTestsFound = false
		} else {
			allNoTestsFound = false
		}
	}

	if len(errs) > 0 {
		return fmt.Errorf("test failures:\n%s", strings.Join(errs, "\n"))
	}
	if allNoTestsFound {
		return ErrNoTestsFound
	}
	return nil
}

func parseBuildOptions(args []string) (core.Options, []string, error) {
	opts := core.Options{Stderr: os.Stderr, RemoveTemp: false}
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--rm", &opts.RemoveTemp).
		String("--gen-dir", &opts.GenDir).
		Int("-count", &opts.Count).
		Bool("--changed", &opts.ChangedOnly).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}

func extractLabelFlags(args []string) (labelExprs []string, remain []string, err error) {
	for i := 0; i < len(args); i++ {
		if args[i] == "--label" {
			if i+1 >= len(args) {
				return nil, nil, fmt.Errorf("--label requires an expression argument")
			}
			labelExprs = append(labelExprs, args[i+1])
			i++
			continue
		}
		remain = append(remain, args[i])
	}
	return labelExprs, remain, nil
}

func parseTestOptions(args []string) (core.Options, []string, error) {
	labelExprs, args, err := extractLabelFlags(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	for _, expr := range labelExprs {
		if err := core.ParseLabelExpr(expr); err != nil {
			return core.Options{}, nil, err
		}
	}

	var sawColor, sawNoColor bool
	for _, arg := range args {
		if arg == "--color" {
			sawColor = true
		}
		if arg == "--no-color" {
			sawNoColor = true
		}
	}
	if sawColor && sawNoColor {
		return core.Options{}, nil, fmt.Errorf("--color and --no-color are mutually exclusive")
	}

	opts := core.Options{Stderr: os.Stderr, RemoveTemp: false, Color: core.ColorAuto}
	var colorFlag, noColorFlag bool
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--rm", &opts.RemoveTemp).
		String("--gen-dir", &opts.GenDir).
		Int("-count", &opts.Count).
		Duration("--timeout", &opts.Timeout).
		Bool("--color", &colorFlag).
		Bool("--no-color", &noColorFlag).
		Bool("--changed", &opts.ChangedOnly).
		Bool("--label-all", &opts.LabelAll).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	if colorFlag {
		opts.Color = core.ColorAlways
	}
	if noColorFlag {
		opts.Color = core.ColorNever
	}
	if opts.LabelAll && len(labelExprs) > 0 {
		return core.Options{}, nil, fmt.Errorf("--label-all and --label are mutually exclusive")
	}
	opts.LabelExprs = labelExprs
	return opts, remainArgs, nil
}

func parseVetOptions(args []string) (core.Options, []string, error) {
	opts := core.Options{Stderr: os.Stderr}
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--changed", &opts.ChangedOnly).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}
