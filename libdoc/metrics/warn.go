package metrics

import "time"

// DefaultSuiteWarnThreshold is the wall-clock budget for a default suite
// before a non-fatal WARNING is emitted (3 minutes).
const DefaultSuiteWarnThreshold = 3 * time.Minute

// ShouldWarnDefaultSuiteSlow reports whether the default-suite slow WARNING
// should fire. True only when defaultSuite is true, total > 0, and
// elapsed > threshold (strict greater-than).
func ShouldWarnDefaultSuiteSlow(defaultSuite bool, total int, elapsed, threshold time.Duration) bool {
	if !defaultSuite {
		return false
	}
	if total <= 0 {
		return false
	}
	if threshold <= 0 {
		threshold = DefaultSuiteWarnThreshold
	}
	return elapsed > threshold
}

// FormatDefaultSuiteSlowWarning returns the fixed WARNING banner recommending
// skill:doctest-review-perf for default suite performance review.
func FormatDefaultSuiteSlowWarning() string {
	return "WARNING: doctest default suite should be fast (run within 3 minutes); " +
		"you're strongly recommended to use skill:doctest-review-perf to optimize " +
		"default test suite performance (doctest skill review-perf --show)\n"
}
