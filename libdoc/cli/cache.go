package cli

import (
	"fmt"

	"github.com/xhd2015/doctest/libdoc/cache"
)

const cacheUsage = `Usage: doctest cache [--clean] [--dry-run]

Show durable doctest cache usage, or remove the cache tree.

Default (no flags):
  Print Cache home, Doctest root, per-bucket human sizes (B/K/M/G), and Total.

Options:
  --clean       Remove $CacheHome/doctest (and override roots outside it)
  --dry-run     With --clean: print [dry-run] would remove lines; do not delete
  -h, --help    Show this help

Environment:
  DOCTEST_CACHE_HOME   Base cache home (default: user cache dir)
  DOCTEST_LEAF_CACHE   Override leaf pass-store root (cleaned when outside main root)
  DOCTEST_METRICS_ROOT Metrics cache base (metrics tree cleaned when outside main root)

Examples:
  doctest cache
  doctest cache --clean --dry-run
  doctest cache --clean
`

func runCache(io stdio, args []string) error {
	if len(args) > 0 && (args[0] == "-h" || args[0] == "--help") {
		fmt.Fprint(io.Out(), cacheUsage)
		return nil
	}
	// Also accept --help later in argv for lessflags-style callers.
	for _, a := range args {
		if a == "-h" || a == "--help" {
			fmt.Fprint(io.Out(), cacheUsage)
			return nil
		}
	}
	opts, err := cache.OptionsFromEnv()
	if err != nil {
		return err
	}
	return cache.Run(io.Out(), args, opts)
}
