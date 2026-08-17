package core

import (
	"os"
	"os/signal"
	"sync"
)

// kindB session: gen roots with outstanding product expose files. A single
// process-wide SIGINT handler is armed while the set is non-empty (or a build
// temp hook is registered) and signal.Stop'd when idle.
//
// On SIGINT the handler serializes with materialize/record/cleanup (same
// mutex), strips product files, runs the optional temp-dir hook, then
// os.Exit(130) only when a CLI session has EnableKindBInterruptExit() held.
// Library callers and go test of this repo do not Exit — generate-only Close
// may leave files (and the handler) for a later leftover sweep, but Ctrl+C
// only cleans; the next SIGINT uses the default terminate.
var (
	kindBMu            sync.Mutex
	kindBRoots         = map[string]struct{}{}
	kindBSigCh         chan os.Signal
	kindBStopCh        chan struct{}
	kindBExitHolders   int
	kindBInterruptHook func()
)

// EnableKindBInterruptExit makes the Kind B SIGINT handler os.Exit(130) after
// cleanup. CLI test/build/vet hold this for the session. Nested calls refcount;
// the returned func pops once. Default is off so library PrepareTree and unit
// tests cannot kill the host process.
func EnableKindBInterruptExit() func() {
	kindBMu.Lock()
	kindBExitHolders++
	kindBMu.Unlock()
	var once sync.Once
	return func() {
		once.Do(func() {
			kindBMu.Lock()
			if kindBExitHolders > 0 {
				kindBExitHolders--
			}
			kindBMu.Unlock()
		})
	}
}

// KindBInterruptExitEnabled reports whether a CLI session currently wants
// os.Exit(130) on SIGINT after Kind B cleanup.
func KindBInterruptExitEnabled() bool {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	return kindBExitHolders > 0
}

// SetKindBInterruptHook registers a callback run after Kind B files are
// stripped on SIGINT (build --rm temp removal). Pass nil to clear. The hook
// must not be invoked while kindBMu is held.
func SetKindBInterruptHook(fn func()) {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	kindBInterruptHook = fn
}

// ArmKindBInterrupt starts the process SIGINT handler even when no Kind B
// list exists yet (so build --rm can remove temps before the first expose).
func ArmKindBInterrupt() {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	ensureKindBInterruptLocked()
}

// DisarmKindBInterruptIfIdle signal.Stop's the handler when no gen root is
// still tracked. Safe when already idle or still dirty.
func DisarmKindBInterruptIfIdle() {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	if len(kindBRoots) == 0 {
		stopKindBInterruptLocked()
	}
}

// KindBInterruptArmed reports whether the process SIGINT handler is installed.
func KindBInterruptArmed() bool {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	return kindBSigCh != nil
}

func registerKindBGenRoot(genRoot string) {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	registerKindBGenRootLocked(genRoot)
}

func registerKindBGenRootLocked(genRoot string) {
	key := absGenRoot(genRoot)
	if key == "" || key == "." {
		return
	}
	kindBRoots[key] = struct{}{}
	ensureKindBInterruptLocked()
}

func unregisterKindBGenRoot(genRoot string) {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	unregisterKindBGenRootLocked(genRoot)
}

func unregisterKindBGenRootLocked(genRoot string) {
	key := absGenRoot(genRoot)
	if key == "" || key == "." {
		return
	}
	delete(kindBRoots, key)
	if len(kindBRoots) == 0 {
		stopKindBInterruptLocked()
	}
}

func kindBGenRootTracked(genRoot string) bool {
	key := absGenRoot(genRoot)
	kindBMu.Lock()
	defer kindBMu.Unlock()
	_, ok := kindBRoots[key]
	return ok
}

func kindBTrackedRootsSnapshotLocked() []string {
	out := make([]string, 0, len(kindBRoots))
	for r := range kindBRoots {
		out = append(out, r)
	}
	return out
}

func ensureKindBInterruptLocked() {
	if kindBSigCh != nil {
		return
	}
	kindBSigCh = make(chan os.Signal, 1)
	kindBStopCh = make(chan struct{})
	signal.Notify(kindBSigCh, os.Interrupt)
	sig, stop := kindBSigCh, kindBStopCh
	go func() {
		for {
			select {
			case <-sig:
				hook, exit := takeKindBInterrupt()
				if hook != nil {
					hook()
				}
				if exit {
					os.Exit(130)
				}
			case <-stop:
				return
			}
		}
	}()
}

// takeKindBInterrupt strips every tracked root under kindBMu, then returns
// the temp hook and exit flag so the caller can run the hook without the lock
// (generateContext.lifecycleMu) and Exit after temps are gone.
func takeKindBInterrupt() (hook func(), exit bool) {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	_ = cleanupAllKindBMaterializedLocked()
	return kindBInterruptHook, kindBExitHolders > 0
}

func stopKindBInterruptLocked() {
	if kindBSigCh == nil {
		return
	}
	signal.Stop(kindBSigCh)
	close(kindBStopCh)
	kindBSigCh = nil
	kindBStopCh = nil
}

// CleanupAllKindBMaterialized strips Kind B product files for every gen root
// that still has a materialized list. Used on SIGINT so mixed gen roots are
// not left behind when the first handler exits.
func CleanupAllKindBMaterialized() error {
	kindBMu.Lock()
	defer kindBMu.Unlock()
	return cleanupAllKindBMaterializedLocked()
}

func cleanupAllKindBMaterializedLocked() error {
	var first error
	for _, root := range kindBTrackedRootsSnapshotLocked() {
		if err := cleanupKindBMaterializedLocked(root); err != nil && first == nil {
			first = err
		}
	}
	return first
}
