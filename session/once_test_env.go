package session

import "syscall"

func lookupSessionEnvImpl() (string, bool) {
	return syscall.Getenv(DoctestSessionIDEnv)
}
