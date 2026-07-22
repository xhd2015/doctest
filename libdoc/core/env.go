package core

import (
	"os"
	"strings"
)

// ChildEnv builds a child-process environment by starting from base (or
// os.Environ() when base is nil) and applying KEY=value overrides with
// key-replace semantics: any existing KEY entries are stripped, then each
// override is appended. Empty values still set KEY= (replace). Override
// entries without '=' are ignored.
//
// Prefer this over blind append(os.Environ(), "KEY=val") so child processes
// do not inherit stale duplicate keys (GOCACHE, DOCTEST_SESSION_ID, etc.).
func ChildEnv(base []string, overrides ...string) []string {
	if base == nil {
		base = os.Environ()
	}
	if len(overrides) == 0 {
		out := make([]string, len(base))
		copy(out, base)
		return out
	}

	strip := make(map[string]struct{}, len(overrides))
	applied := make([]string, 0, len(overrides))
	for _, o := range overrides {
		k, _, ok := strings.Cut(o, "=")
		if !ok || k == "" {
			continue
		}
		// Last override for a key wins: drop earlier applied for same key.
		if _, seen := strip[k]; seen {
			for i := len(applied) - 1; i >= 0; i-- {
				ak, _, _ := strings.Cut(applied[i], "=")
				if ak == k {
					applied = append(applied[:i], applied[i+1:]...)
					break
				}
			}
		}
		strip[k] = struct{}{}
		applied = append(applied, o)
	}
	if len(strip) == 0 {
		out := make([]string, len(base))
		copy(out, base)
		return out
	}

	out := make([]string, 0, len(base)+len(applied))
	for _, e := range base {
		k, _, _ := strings.Cut(e, "=")
		if _, drop := strip[k]; drop {
			continue
		}
		out = append(out, e)
	}
	out = append(out, applied...)
	return out
}
