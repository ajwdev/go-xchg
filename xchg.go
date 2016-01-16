// +build linux
// +build amd64 arm64

package xchg

import (
	"golang.org/x/sys/unix"
	"unsafe"
)

const (
	// These values were pulled from here
	// http://lxr.free-electrons.com/source/include/uapi/linux/fs.h#L38

	// Don't overwrite target
	NOREPLACE = 1
	// Exchange source and destination
	EXCHANGE = 2
	// Whiteout source
	WHITEOUT = 4
)

// XXX This is a trick that is used in the Go syscall library to trick
// the compiler into keeping an object alive for the length of a
// syscall. It shouldn't be nessecary for Go 1.6 so we should consider
// removing this around Go 1.7. See accompanying 'use.s' assembly file
// as well as the following commit message for the technical details.
// The tl;dr is that our unsafe.Pointer could be garbage collected right
// after we get its location via the `uintptr` call which means we'd be
// passing garbage to kernel for the syscall.
//
// https://github.com/golang/go/commit/cf622d758cd51cfa09f5b503d323c81ed3a5541e#diff-264325f3e02feb10326578a776808fb7
func use(p unsafe.Pointer)

func Renameat2(olddirfd int, oldpath string, newdirfd int, newpath string, flags int) (err error) {
	oldPtr, err := unix.BytePtrFromString(oldpath)
	if err != nil {
		return
	}
	newPtr, err := unix.BytePtrFromString(newpath)
	if err != nil {
		return
	}
	_, _, errno := unix.Syscall6(
		SYS_RENAMEAT2,
		uintptr(olddirfd),
		uintptr(unsafe.Pointer(oldPtr)),
		uintptr(newdirfd),
		uintptr(unsafe.Pointer(newPtr)),
		uintptr(flags),
		0,
	)
	use(unsafe.Pointer(oldPtr))
	use(unsafe.Pointer(newPtr))

	if errno != 0 {
		// In the syscall module the authors box a couple of common errors
		// (i.e EAGAIN, EINVAL, and ENOENT). Is that worth doing here?
		err = errno
	}

	return
}

func Exchange(oldpath string, newpath string) (err error) {
	// TODO Should we return an os.LinkError on error? What about on ENOSYS?
	return Renameat2(unix.AT_FDCWD, oldpath, unix.AT_FDCWD, newpath, EXCHANGE)
}
