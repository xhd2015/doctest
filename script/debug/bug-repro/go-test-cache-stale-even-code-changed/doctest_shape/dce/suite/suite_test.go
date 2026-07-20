package suite

import (
	"testing"

	_ "dt.local/leaf"
	"dt.local/registry"
)

func TestDoctestSuite(t *testing.T) {
	for _, e := range registry.All() {
		e := e
		t.Run(e.Path, func(t *testing.T) { e.Fn(t) })
	}
}
