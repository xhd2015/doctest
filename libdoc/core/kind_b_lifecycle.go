package core

import (
	"os"
	"os/signal"
	"sync"
)

// kindB session: gen roots with outstanding product expose files. A single
// process-wide SIGINT handler is armed while the set is non-empty and
// signal.Stop'd when the last root is fully cleaned — so library callers
// (and go test of this repo) do not keep a stale os.Exit(130) after the run.
var (
	kindBMu     sync.Mutex
	kindBRoots  = map[string]struct{}{}
	kindBSigCh  chan os.Signal
	kindBStopCh chan struct{}
)

func registerKindBGenRoot(genRoot string) {
	key := absGenRoot(genRoot)
	if key == "" || key == "." {
		return
	}
	kindBMu.Lock()
	defer kindBMu.Unlock()
	kindBRoots[key] = struct{}{}
	ensureKindBInterruptLocked()
}

func unregisterKindBGenRoot(genRoot string) {
	key := absGenRoot(genRoot)
	if key == "" || key == "." {
		return
	}
	kindBMu.Lock()
	defer kindBMu.Unlock()
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

func kindBTrackedRootsSnapshot() []string {
	kindBMu.Lock()
	defer kindBMu.Unlock()
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
		select {
		case <-sig:
			_ = CleanupAllKindBMaterialized()
			os.Exit(130)
		case <-stop:
			return
		}
	}()
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
	var first error
	for _, root := range kindBTrackedRootsSnapshot() {
		if err := CleanupKindBMaterialized(root); err != nil && first == nil {
			first = err
		}
	}
	return first
}
