// +build arm64,linux

package xchg

import (
	"syscall"
)

// Go has a constant for renameat2 on arm64 already so just reuse it.
const SYS_RENAMEAT2 = syscall.SYS_RENAMEAT2
