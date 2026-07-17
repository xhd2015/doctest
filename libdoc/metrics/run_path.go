package metrics

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// RunFilePath builds the canonical run JSONL path under the cache directory.
// Layout: $cacheDir/doctest/metrics/<projectID>/runs/YYYY-MM-DD-HH-MM-SS-NN-<suffix>.jsonl
// Time is formatted in UTC; NN is zero-padded to two digits.
func RunFilePath(cacheDir, projectID string, t time.Time, nn int, suffix string) string {
	t = t.UTC()
	name := fmt.Sprintf("%s-%02d-%s.jsonl",
		t.Format("2006-01-02-15-04-05"),
		nn,
		suffix,
	)
	return filepath.Join(cacheDir, "doctest", "metrics", projectID, "runs", name)
}

// CreateRunFile creates a new exclusive run file under the cache layout.
// It tries NN from 00–99 with the given suffix (or a generated one if empty).
// On clash it bumps NN and, if needed, regenerates a unique suffix.
func CreateRunFile(cacheDir, projectID string, at time.Time, suffix string) (string, error) {
	at = at.UTC()
	dir := filepath.Join(cacheDir, "doctest", "metrics", projectID, "runs")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", err
	}

	baseSuffix := suffix
	if baseSuffix == "" {
		var err error
		baseSuffix, err = randomSuffix(8)
		if err != nil {
			return "", err
		}
	}

	trySuffix := baseSuffix
	for attempt := 0; attempt < 200; attempt++ {
		if attempt > 0 && attempt%100 == 0 {
			// Exhausted NN range with this suffix; mint a new one.
			var err error
			trySuffix, err = randomSuffix(8)
			if err != nil {
				return "", err
			}
		}
		nn := attempt % 100
		path := RunFilePath(cacheDir, projectID, at, nn, trySuffix)
		f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
		if err != nil {
			if os.IsExist(err) {
				continue
			}
			return "", err
		}
		if err := f.Close(); err != nil {
			return "", err
		}
		return path, nil
	}
	return "", fmt.Errorf("metrics: failed to create exclusive run file under %s", dir)
}

func randomSuffix(nBytes int) (string, error) {
	b := make([]byte, nBytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b)[:nBytes], nil // nBytes hex chars if we take half; use full hex length
}
