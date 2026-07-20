package mid

import (
	"testing"

	"dt.local/droot"
)

func Setup(t *testing.T, req *droot.Request) error {
	req.WorkDir = "MID_V2"
	return nil
}
