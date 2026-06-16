package runner

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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
	return processArgs(args, "test", parseTestOptions, func(dir string, opts core.Options) error {
		err := runnerbuild.Test(dir, opts)
		if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
			return ErrNoTestsFound
		}
		return err
	})
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
		return path_resolve.RunForDirs(path_resolve.ExtractBasePath(arg), func(dir string) error {
			root, _ := path_resolve.ResolveRoot(dir)
			if root == "" {
				root = dir
			}
			o := opts
			o.SubDir = dir
			return fn(root, o)
		})
	}
	targetDir, _ := filepath.Abs(arg)
	root, ok := path_resolve.ResolveRoot(targetDir)
	if !ok {
		root = targetDir
	}
	opts.SubDir = targetDir
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
		targetDir, _ := filepath.Abs(arg)
		root, ok := path_resolve.ResolveRoot(targetDir)
		if !ok {
			root = targetDir
		}
		o := opts
		o.SubDir = targetDir
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
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}

func parseTestOptions(args []string) (core.Options, []string, error) {
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
		Bool("--color", &colorFlag).
		Bool("--no-color", &noColorFlag).
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
	return opts, remainArgs, nil
}

func parseVetOptions(args []string) (core.Options, []string, error) {
	opts := core.Options{Stderr: os.Stderr}
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}
