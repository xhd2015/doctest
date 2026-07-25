// Package debug parses DOCTEST_DEBUG, a GODEBUG-style env for engine-internal
// diagnostics (not part of the public leaf harness / d.DOCTEST_* surface).
//
// Format: comma-separated key=value pairs, e.g.
//
//	DOCTEST_DEBUG=bypass-go-test=1,cpuprofile=/tmp/cpu.pprof
//
// Unknown keys are errors (fail closed).
package debug

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"runtime/pprof"
	"strings"
)

// EnvName is the process environment variable name.
const EnvName = "DOCTEST_DEBUG"

// Settings holds parsed DOCTEST_DEBUG keys.
type Settings struct {
	// BypassGoTest skips host-driven go test exec after generate / workspace
	// write+tidy. Prepare and hub fan-in still run.
	BypassGoTest bool

	// GenPlan prints generate plan and result hierarchy trees on stderr
	// (DOCTEST_DEBUG=gen-plan=1). Never pollutes test stdout / JSON.
	GenPlan bool

	// Host process profiles (paths; empty = off). Written on Stop from StartProfiles.
	CPUProfile   string
	MemProfile   string
	BlockProfile string
}

// Parse reads a GODEBUG-style string (comma-separated key=value).
// Empty input yields zero Settings and nil error.
func Parse(s string) (Settings, error) {
	var out Settings
	s = strings.TrimSpace(s)
	if s == "" {
		return out, nil
	}
	for _, part := range strings.Split(s, ",") {
		part = strings.TrimSpace(part)
		if part == "" {
			continue
		}
		key, val, ok := strings.Cut(part, "=")
		key = strings.TrimSpace(key)
		if key == "" {
			return Settings{}, fmt.Errorf("%s: empty key in %q", EnvName, part)
		}
		if !ok {
			return Settings{}, fmt.Errorf("%s: %q must be key=value", EnvName, part)
		}
		val = strings.TrimSpace(val)
		switch key {
		case "bypass-go-test":
			on, err := parseBool(val)
			if err != nil {
				return Settings{}, fmt.Errorf("%s: bypass-go-test: %w", EnvName, err)
			}
			out.BypassGoTest = on
		case "gen-plan":
			on, err := parseBool(val)
			if err != nil {
				return Settings{}, fmt.Errorf("%s: gen-plan: %w", EnvName, err)
			}
			out.GenPlan = on
		case "cpuprofile":
			if val == "" {
				return Settings{}, fmt.Errorf("%s: cpuprofile requires a non-empty path", EnvName)
			}
			out.CPUProfile = val
		case "memprofile":
			if val == "" {
				return Settings{}, fmt.Errorf("%s: memprofile requires a non-empty path", EnvName)
			}
			out.MemProfile = val
		case "blockprofile":
			if val == "" {
				return Settings{}, fmt.Errorf("%s: blockprofile requires a non-empty path", EnvName)
			}
			out.BlockProfile = val
		default:
			return Settings{}, fmt.Errorf("%s: unknown key %q", EnvName, key)
		}
	}
	return out, nil
}

// FromEnv parses os.Getenv(EnvName).
func FromEnv() (Settings, error) {
	return Parse(os.Getenv(EnvName))
}

func parseBool(val string) (bool, error) {
	switch strings.ToLower(val) {
	case "1", "true", "t", "yes", "y", "on":
		return true, nil
	case "0", "false", "f", "no", "n", "off":
		return false, nil
	default:
		return false, fmt.Errorf("invalid bool %q (want 0|1|true|false)", val)
	}
}

// StartProfiles enables host CPU / mem / block profiling for the process.
// Returns a stop function that must be deferred; it writes mem/block profiles
// and stops CPU profiling. Safe to call with zero settings (stop is a no-op).
//
// Paths are abs-resolved against the process cwd. Parent directories are created.
func StartProfiles(s Settings) (stop func(), err error) {
	var (
		cpuFile   *os.File
		cpuOn     bool
		blockRate int
	)
	stop = func() {
		if cpuOn {
			pprof.StopCPUProfile()
			cpuOn = false
		}
		if cpuFile != nil {
			_ = cpuFile.Close()
			cpuFile = nil
		}
		if s.MemProfile != "" {
			if werr := writeMemProfile(s.MemProfile); werr != nil {
				fmt.Fprintf(os.Stderr, "doctest: %s memprofile: %v\n", EnvName, werr)
			}
		}
		if s.BlockProfile != "" {
			if werr := writeBlockProfile(s.BlockProfile); werr != nil {
				fmt.Fprintf(os.Stderr, "doctest: %s blockprofile: %v\n", EnvName, werr)
			}
			if blockRate != 0 {
				runtime.SetBlockProfileRate(0)
			}
		}
	}

	if s.CPUProfile == "" && s.MemProfile == "" && s.BlockProfile == "" {
		return stop, nil
	}

	// Resolve paths up front so start failures are clean.
	if s.CPUProfile != "" {
		p, rerr := resolveProfilePath(s.CPUProfile)
		if rerr != nil {
			return nil, fmt.Errorf("%s: cpuprofile: %w", EnvName, rerr)
		}
		s.CPUProfile = p
	}
	if s.MemProfile != "" {
		p, rerr := resolveProfilePath(s.MemProfile)
		if rerr != nil {
			return nil, fmt.Errorf("%s: memprofile: %w", EnvName, rerr)
		}
		s.MemProfile = p
	}
	if s.BlockProfile != "" {
		p, rerr := resolveProfilePath(s.BlockProfile)
		if rerr != nil {
			return nil, fmt.Errorf("%s: blockprofile: %w", EnvName, rerr)
		}
		s.BlockProfile = p
	}

	if s.BlockProfile != "" {
		blockRate = 1
		runtime.SetBlockProfileRate(blockRate)
		fmt.Fprintf(os.Stderr, "doctest: %s blockprofile=%s\n", EnvName, s.BlockProfile)
	}
	if s.CPUProfile != "" {
		if err := os.MkdirAll(filepath.Dir(s.CPUProfile), 0755); err != nil {
			stop()
			return nil, fmt.Errorf("%s: cpuprofile mkdir: %w", EnvName, err)
		}
		f, oerr := os.Create(s.CPUProfile)
		if oerr != nil {
			stop()
			return nil, fmt.Errorf("%s: cpuprofile create: %w", EnvName, oerr)
		}
		if err := pprof.StartCPUProfile(f); err != nil {
			_ = f.Close()
			stop()
			return nil, fmt.Errorf("%s: start cpuprofile: %w", EnvName, err)
		}
		cpuFile = f
		cpuOn = true
		fmt.Fprintf(os.Stderr, "doctest: %s cpuprofile=%s\n", EnvName, s.CPUProfile)
	}
	if s.MemProfile != "" {
		fmt.Fprintf(os.Stderr, "doctest: %s memprofile=%s\n", EnvName, s.MemProfile)
	}
	return stop, nil
}

func resolveProfilePath(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("empty path")
	}
	if filepath.IsAbs(p) {
		return filepath.Clean(p), nil
	}
	abs, err := filepath.Abs(p)
	if err != nil {
		return "", err
	}
	return abs, nil
}

func writeMemProfile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	runtime.GC()
	return pprof.WriteHeapProfile(f)
}

func writeBlockProfile(path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	return pprof.Lookup("block").WriteTo(f, 0)
}
