package build

import (
	"fmt"

	"github.com/xhd2015/doctest/libdoc/core"
)

// appendOptsGoTestFlags appends user-facing go test flags from opts onto base.
// Does not include "test", "-mod=mod", package patterns, -json, or xgo extras.
func appendOptsGoTestFlags(flagArgs []string, opts core.Options) []string {
	if opts.Verbose {
		flagArgs = append(flagArgs, "-v")
	}
	if opts.Count > 0 {
		flagArgs = append(flagArgs, fmt.Sprintf("-count=%d", opts.Count))
	}
	if opts.ForceWithFlagA {
		// go build/test -a: force rebuilding packages that are already up-to-date.
		flagArgs = append(flagArgs, "-a")
	}
	// nil = omit (go default 10m); non-nil including 0 = pass -timeout=…
	if opts.Timeout != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-timeout=%s", *opts.Timeout))
	}
	if opts.CPUProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-cpuprofile=%s", opts.CPUProfile))
	}
	if opts.MemProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-memprofile=%s", opts.MemProfile))
	}
	if opts.MemProfileRate != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-memprofilerate=%d", *opts.MemProfileRate))
	}
	if opts.BlockProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-blockprofile=%s", opts.BlockProfile))
	}
	if opts.BlockProfileRate != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-blockprofilerate=%d", *opts.BlockProfileRate))
	}
	if opts.MutexProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-mutexprofile=%s", opts.MutexProfile))
	}
	if opts.MutexProfileFraction != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-mutexprofilefraction=%d", *opts.MutexProfileFraction))
	}
	if opts.Trace != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-trace=%s", opts.Trace))
	}
	if opts.OutputDir != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-outputdir=%s", opts.OutputDir))
	}
	if opts.CoverProfile != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-coverprofile=%s", opts.CoverProfile))
	}
	if opts.CoverMode != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-covermode=%s", opts.CoverMode))
	}
	if opts.CoverPkg != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-coverpkg=%s", opts.CoverPkg))
	}
	// -cover when explicit or implied by covermode/coverpkg (coverprofile alone
	// is enough for go; still emit -cover when user asked or mode/pkg set).
	if opts.Cover || opts.CoverMode != "" || opts.CoverPkg != "" {
		flagArgs = append(flagArgs, "-cover")
	}
	if opts.Race {
		flagArgs = append(flagArgs, "-race")
	}
	if opts.Short {
		flagArgs = append(flagArgs, "-short")
	}
	if opts.FailFast {
		flagArgs = append(flagArgs, "-failfast")
	}
	if opts.Parallel != nil {
		flagArgs = append(flagArgs, fmt.Sprintf("-parallel=%d", *opts.Parallel))
	}
	if opts.Shuffle != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-shuffle=%s", opts.Shuffle))
	}
	if opts.Tags != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-tags=%s", opts.Tags))
	}
	if opts.Gcflags != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-gcflags=%s", opts.Gcflags))
	}
	if opts.Ldflags != "" {
		flagArgs = append(flagArgs, fmt.Sprintf("-ldflags=%s", opts.Ldflags))
	}
	return flagArgs
}

// errCoverProfileMultiPackage is returned when -coverprofile is set with more
// than one package arg (go would reject or we would overwrite a single file).
func errCoverProfileMultiPackage(n int) error {
	return fmt.Errorf("-coverprofile cannot be used with multiple packages (%d); run a single package/tree or omit -coverprofile (multi-package profile merge is not supported)", n)
}

// checkCoverProfilePackages rejects multi-package -coverprofile (no merge yet).
func checkCoverProfilePackages(opts core.Options, packageArgs []string) error {
	if opts.CoverProfile != "" && len(packageArgs) > 1 {
		return errCoverProfileMultiPackage(len(packageArgs))
	}
	return nil
}
