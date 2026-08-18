package core

import (
	"os"
	"os/signal"
	"sync"
)

// Expose session: gen roots with outstanding product expose files. A single
// process-wide SIGINT handler is armed while the set is non-empty (or a build
// temp hook is registered) and signal.Stop'd when idle.
//
// On SIGINT the handler serializes with materialize/record/cleanup (same
// mutex), strips product files, runs the optional temp-dir hook, then
// os.Exit(130) only when a CLI session has EnableExposeInterruptExit() held.
// Library callers and go test of this repo do not Exit — generate-only Close
// may leave files for a later leftover sweep. The first SIGINT still cleans;
// the handler then stops so a leftover cannot swallow every later Ctrl+C.
var (
	exposeMu            sync.Mutex
	exposeRoots         = map[string]struct{}{}
	exposeSigCh         chan os.Signal
	exposeStopCh        chan struct{}
	exposeExitHolders   int
	exposeInterruptHook func()
)

// EnableExposeInterruptExit makes the expose SIGINT handler os.Exit(130) after
// cleanup. CLI test/build/vet hold this for the session. Nested calls refcount;
// the returned func pops once. Default is off so library PrepareTree and unit
// tests cannot kill the host process.
func EnableExposeInterruptExit() func() {
	exposeMu.Lock()
	exposeExitHolders++
	exposeMu.Unlock()
	var once sync.Once
	return func() {
		once.Do(func() {
			exposeMu.Lock()
			if exposeExitHolders > 0 {
				exposeExitHolders--
			}
			exposeMu.Unlock()
		})
	}
}

// ExposeInterruptExitEnabled reports whether a CLI session currently wants
// os.Exit(130) on SIGINT after expose cleanup.
func ExposeInterruptExitEnabled() bool {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	return exposeExitHolders > 0
}

// SetExposeInterruptHook registers a callback run after expose files are
// stripped on SIGINT (build --rm temp removal). Pass nil to clear. The hook
// must not be invoked while exposeMu is held.
func SetExposeInterruptHook(fn func()) {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	exposeInterruptHook = fn
}

// ArmExposeInterrupt starts the process SIGINT handler even when no expose
// list exists yet (so build --rm can remove temps before the first expose).
func ArmExposeInterrupt() {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	ensureExposeInterruptLocked()
}

// DisarmExposeInterruptIfIdle signal.Stop's the handler when no gen root is
// still tracked. Safe when already idle or still dirty.
func DisarmExposeInterruptIfIdle() {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	if len(exposeRoots) == 0 {
		stopExposeInterruptLocked()
	}
}

// ExposeInterruptArmed reports whether the process SIGINT handler is installed.
func ExposeInterruptArmed() bool {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	return exposeSigCh != nil
}

func registerExposeGenRoot(genRoot string) {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	registerExposeGenRootLocked(genRoot)
}

func registerExposeGenRootLocked(genRoot string) {
	key := absGenRoot(genRoot)
	if key == "" || key == "." {
		return
	}
	exposeRoots[key] = struct{}{}
	ensureExposeInterruptLocked()
}

func unregisterExposeGenRoot(genRoot string) {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	unregisterExposeGenRootLocked(genRoot)
}

func unregisterExposeGenRootLocked(genRoot string) {
	key := absGenRoot(genRoot)
	if key == "" || key == "." {
		return
	}
	delete(exposeRoots, key)
	if len(exposeRoots) == 0 {
		stopExposeInterruptLocked()
	}
}

func exposeGenRootTracked(genRoot string) bool {
	key := absGenRoot(genRoot)
	exposeMu.Lock()
	defer exposeMu.Unlock()
	_, ok := exposeRoots[key]
	return ok
}

func exposeTrackedRootsSnapshotLocked() []string {
	out := make([]string, 0, len(exposeRoots))
	for r := range exposeRoots {
		out = append(out, r)
	}
	return out
}

func ensureExposeInterruptLocked() {
	if exposeSigCh != nil {
		return
	}
	exposeSigCh = make(chan os.Signal, 1)
	exposeStopCh = make(chan struct{})
	signal.Notify(exposeSigCh, os.Interrupt)
	sig, stop := exposeSigCh, exposeStopCh
	go func() {
		for {
			select {
			case <-sig:
				hook, exit := takeExposeInterrupt()
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

// takeExposeInterrupt strips every tracked root under exposeMu, then returns
// the temp hook and exit flag so the caller can run the hook without the lock
// (generateContext.lifecycleMu) and Exit after temps are gone.
func takeExposeInterrupt() (hook func(), exit bool) {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	_ = cleanupAllExposeMaterializedLocked()
	return finishExposeInterruptLocked()
}

// finishExposeInterruptLocked returns the temp hook and CLI-exit flag.
// When exit is false (library / go test), the process handler is stopped even
// if leftovers remain tracked — otherwise SIGINT stays consumed forever.
func finishExposeInterruptLocked() (hook func(), exit bool) {
	exit = exposeExitHolders > 0
	hook = exposeInterruptHook
	if !exit {
		stopExposeInterruptLocked()
	}
	return hook, exit
}

func stopExposeInterruptLocked() {
	if exposeSigCh == nil {
		return
	}
	signal.Stop(exposeSigCh)
	close(exposeStopCh)
	exposeSigCh = nil
	exposeStopCh = nil
}

// CleanupAllExposeMaterialized strips expose product files for every gen root
// that still has a materialized list. Used on SIGINT so mixed gen roots are
// not left behind when the first handler exits.
func CleanupAllExposeMaterialized() error {
	exposeMu.Lock()
	defer exposeMu.Unlock()
	return cleanupAllExposeMaterializedLocked()
}

func cleanupAllExposeMaterializedLocked() error {
	var first error
	for _, root := range exposeTrackedRootsSnapshotLocked() {
		if err := cleanupExposeMaterializedLocked(root); err != nil && first == nil {
			first = err
		}
	}
	return first
}
