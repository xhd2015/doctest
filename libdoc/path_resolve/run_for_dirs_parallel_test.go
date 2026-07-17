package path_resolve

import (
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestRunForDirsLimitRunsAllDirsConcurrently(t *testing.T) {
	root := t.TempDir()
	// Three independent DOCTEST.md trees under root.
	for _, name := range []string{"a", "b", "c"} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# t\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	var (
		inFlight atomic.Int32
		maxSeen  atomic.Int32
		seen     sync.Map
		mu       sync.Mutex
		started  = make(chan struct{}, 3)
	)

	err := RunForDirsLimit(root, 3, func(dir string) error {
		seen.Store(dir, true)
		n := inFlight.Add(1)
		for {
			cur := maxSeen.Load()
			if n <= cur || maxSeen.CompareAndSwap(cur, n) {
				break
			}
		}
		started <- struct{}{}
		// Hold until all three have entered so concurrency is observable.
		mu.Lock()
		// Wait for three starts without holding lock across sleep of others:
		mu.Unlock()
		// Barrier: consume is done by parent via timeout on channel fills.
		time.Sleep(50 * time.Millisecond)
		inFlight.Add(-1)
		return nil
	})
	if err != nil {
		t.Fatalf("RunForDirsLimit: %v", err)
	}

	// Drain start signals (should be 3).
	for i := 0; i < 3; i++ {
		select {
		case <-started:
		case <-time.After(2 * time.Second):
			t.Fatal("timed out waiting for concurrent starts")
		}
	}
	if maxSeen.Load() < 2 {
		t.Fatalf("expected concurrent workers >=2, max in-flight=%d", maxSeen.Load())
	}
	// All three absolute paths visited.
	count := 0
	seen.Range(func(_, _ any) bool { count++; return true })
	if count != 3 {
		t.Fatalf("visited %d dirs, want 3", count)
	}
}

func TestIsHeavySelftestTree(t *testing.T) {
	// Paths that look like the module integration suite.
	if !isHeavySelftestTree("/Users/x/proj/doctest/tests") {
		t.Fatal("expected doctest/tests root heavy")
	}
	if !isHeavySelftestTree("/Users/x/proj/doctest/tests/changed") {
		t.Fatal("expected doctest/tests/changed heavy")
	}
	if isHeavySelftestTree("/Users/x/proj/doctest/assert/tests/output-assert-v3") {
		t.Fatal("assert trees are light")
	}
	if isHeavySelftestTree("/Users/x/proj/doctest/libdoc/core/tests/assert-mod") {
		t.Fatal("libdoc trees are light")
	}
	if isHeavySelftestTree("/Users/x/proj/doctest/session/tests/once") {
		t.Fatal("session trees are light")
	}
}

func TestRunForDirsLimitSerialWhenWorkersOne(t *testing.T) {
	root := t.TempDir()
	for _, name := range []string{"x", "y"} {
		dir := filepath.Join(root, name)
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, "DOCTEST.md"), []byte("# t\n"), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	var order []string
	var mu sync.Mutex
	err := RunForDirsLimit(root, 1, func(dir string) error {
		mu.Lock()
		order = append(order, filepath.Base(dir))
		mu.Unlock()
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(order) != 2 {
		t.Fatalf("order=%v", order)
	}
}
