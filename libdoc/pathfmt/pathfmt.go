package pathfmt

import (
	"os"
	"path/filepath"
	"strings"
)

func evalPath(path string) string {
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		return resolved
	}
	return path
}

// DisplayPath returns a shortened form of path for human-readable CLI output only.
// It must not be used for file I/O or exec.Command.Dir.
func DisplayPath(path string) string {
	if path == "" {
		return path
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return path
	}
	if cwd, err := os.Getwd(); err == nil {
		if cwdAbs, err := filepath.Abs(cwd); err == nil {
			cwdEval := evalPath(cwdAbs)
			absEval := evalPath(abs)
			if absEval == cwdEval {
				return "."
			}
			rel, err := filepath.Rel(cwdEval, absEval)
			if err == nil && rel != "." && !strings.HasPrefix(rel, "..") {
				return "." + string(filepath.Separator) + rel
			}
		}
	}
	if home, err := os.UserHomeDir(); err == nil {
		if abs == home {
			return "~"
		}
		homePrefix := home + string(filepath.Separator)
		if strings.HasPrefix(abs, homePrefix) {
			return "~" + strings.TrimPrefix(abs, home)
		}
	}
	return abs
}

// ExpandDisplayPath converts a display path (with ~ prefix) back to an absolute path
// for filesystem operations. Non-display paths are returned unchanged.
func ExpandDisplayPath(path string) string {
	if path == "" || !strings.HasPrefix(path, "~") {
		return path
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return path
	}
	if path == "~" {
		return home
	}
	suffix := strings.TrimPrefix(path, "~")
	suffix = strings.TrimPrefix(suffix, string(filepath.Separator))
	if suffix == "" {
		return home
	}
	return filepath.Join(home, suffix)
}