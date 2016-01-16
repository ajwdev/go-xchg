// +build amd64,linux

package xchg

// For some reason this constant is defined in Go for arm64 and mips64 but not
// for x86{,-64} so we need to define it ourselves.
const SYS_RENAMEAT2 = 316
