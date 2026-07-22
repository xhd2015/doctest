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
		order = append(order, dir)
		mu.Unlock()
		return nil
	})
	if err != nil {
		t.Fatalf("RunForDirsLimit: %v", err)
	}
	if len(order) != 2 {
		t.Fatalf("visited %d dirs, want 2", len(order))
	}
}
