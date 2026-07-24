package mid

import (
	"testing"

	"dt.local/droot"
)

func Setup(t *testing.T, d *session.Doctest, req *droot.Request) error {
	req.WorkDir = "MID_V1"
	return nil
}
