package build

import (
	"os/exec"
)

func NeedsBuildVCSFlag(dir string) bool {
	git, err := exec.LookPath("git")
	if err != nil {
		return true
	}
	cmd := exec.Command(git, "-C", dir, "rev-parse", "--is-inside-work-tree")
	out, err := cmd.Output()
	if err != nil {
		return true
	}
	return string(out) != "true\n" && string(out) != "true\r\n"
}
