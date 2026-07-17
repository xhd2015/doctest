package metrics

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// EnvMetricsRoot overrides the metrics cache root when set.
const EnvMetricsRoot = "DOCTEST_METRICS_ROOT"

// ResolveMetricsRoot returns MetricsRoot override, DOCTEST_METRICS_ROOT, or user cache dir.
func ResolveMetricsRoot(override string) string {
	if override != "" {
		return override
	}
	if v := os.Getenv(EnvMetricsRoot); v != "" {
		return v
	}
	base, err := os.UserCacheDir()
	if err != nil {
		return ""
	}
	return base
}

// ProjectMetricsDir returns $cacheDir/doctest/metrics/<projectID>.
func ProjectMetricsDir(cacheDir, projectID string) string {
	return filepath.Join(cacheDir, "doctest", "metrics", projectID)
}

// ProjectRunsDir returns $cacheDir/doctest/metrics/<projectID>/runs.
func ProjectRunsDir(cacheDir, projectID string) string {
	return filepath.Join(cacheDir, "doctest", "metrics", projectID, "runs")
}

// ProjectIDForDir resolves project id from git origin in dir, or nogit fallback.
func ProjectIDForDir(dir string) string {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	origin := GitRemoteOrigin(abs)
	if id := ProjectIDFromOrigin(origin); id != "" {
		return id
	}
	return ProjectIDFallback(abs)
}

// GitRemoteOrigin returns `git remote get-url origin` for dir, or empty.
func GitRemoteOrigin(dir string) string {
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}
