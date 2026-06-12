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
)

var ErrNoTestsFound = errors.New("no tests found")

func Build(dir string) error {
	return runnerbuild.Build(dir, core.Options{RemoveTemp: false})
}

func BuildArgs(args []string) error {
	opts, remainArgs, err := parseBuildOptions(args)
	if err != nil {
		return err
	}
	if len(remainArgs) != 1 {
		return fmt.Errorf("build requires <dir>")
	}
	arg := remainArgs[0]
	if arg == "..." {
		return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
	}
	if isDotDotDotPattern(arg) {
		return runForDirs(extractBasePath(arg), func(dir string) error {
			root, _ := ResolveRoot(dir)
			if root == "" {
				root = dir
			}
			o := opts
			o.SubDir = dir
			err := runnerbuild.Build(root, o)
			if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
				return ErrNoTestsFound
			}
			return err
		})
	}
	targetDir, _ := filepath.Abs(remainArgs[0])
	root, ok := ResolveRoot(targetDir)
	if !ok {
		return ErrNoTestsFound
	}
	opts.SubDir = targetDir
	err = runnerbuild.Build(root, opts)
	if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
		return ErrNoTestsFound
	}
	return err
}

func Test(args []string) error {
	opts, remainArgs, err := parseTestOptions(args)
	if err != nil {
		return err
	}
	if len(remainArgs) != 1 {
		return fmt.Errorf("test requires <dir>")
	}
	arg := remainArgs[0]
	if arg == "..." {
		return fmt.Errorf("bare '...' pattern is not supported; use './...' or 'path/...' instead")
	}
	if isDotDotDotPattern(arg) {
		return runForDirs(extractBasePath(arg), func(dir string) error {
			root, _ := ResolveRoot(dir)
			if root == "" {
				root = dir
			}
			o := opts
			o.SubDir = dir
			err := runnerbuild.Test(root, o)
			if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
				return ErrNoTestsFound
			}
			return err
		})
	}
	targetDir, _ := filepath.Abs(remainArgs[0])
	root, ok := ResolveRoot(targetDir)
	if !ok {
		return ErrNoTestsFound
	}
	opts.SubDir = targetDir
	err = runnerbuild.Test(root, opts)
	if err != nil && strings.Contains(err.Error(), "no runnable test cases found") {
		return ErrNoTestsFound
	}
	return err
}

func runForDirs(basePath string, fn func(dir string) error) error {
	dirs, err := findDotDotDotDirs(basePath)
	if err != nil {
		return err
	}
	if len(dirs) == 0 {
		return ErrNoTestsFound
	}
	var errs []string
	for _, dir := range dirs {
		if err := fn(dir); err != nil {
			if errors.Is(err, ErrNoTestsFound) {
				continue
			}
			errs = append(errs, fmt.Sprintf("%s: %v", dir, err))
		}
	}
	if len(errs) > 0 {
		return fmt.Errorf("test failures:\n%s", strings.Join(errs, "\n"))
	}
	return nil
}

func findDotDotDotDirs(basePath string) ([]string, error) {
	if basePath == "" || basePath == "." {
		dirs, err := FindDOCTestDirs(".")
		if err == nil {
			return dirs, nil
		}
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		return FindDOCTestDirsFromSubdirs(".")
	}
	dirs, err := FindDOCTestDirsWithBase(basePath, basePath)
	if err != nil {
		return nil, err
	}
	return dirs, nil
}

func isDotDotDotPattern(arg string) bool {
	if len(arg) > 0 && arg[0] == '/' {
		return false
	}
	return strings.HasSuffix(arg, "/...")
}

func extractBasePath(arg string) string {
	base := strings.TrimSuffix(arg, "/...")
	base = strings.TrimPrefix(base, "./")
	if base == "" {
		return "."
	}
	return base
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
	opts := core.Options{Stderr: os.Stderr, RemoveTemp: false}
	remainArgs, err := lessflags.Bool("-v,--verbose", &opts.Verbose).
		Bool("--rm", &opts.RemoveTemp).
		Int("-count", &opts.Count).
		Parse(args)
	if err != nil {
		return core.Options{}, nil, err
	}
	return opts, remainArgs, nil
}
