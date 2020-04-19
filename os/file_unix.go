// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package os

import (
	"os"
	"runtime"
	"syscall"
)

// OpenFile opens a file an returns a struct containing the file and its file descriptor
func OpenFile(name string, flag int, perm os.FileMode) (file *os.File, err error) {
	file, _, err = openFile(name, flag, perm)
	return
}

func OpenFileWithDescriptor(name string, flag int, perm os.FileMode) (fwd FileWithDescriptor, err error) {
	fwd.File, fwd.fd, err = openFile(name, flag, perm)
	return
}

func openFile(name string, flag int, perm os.FileMode) (f *os.File, fd int, err error) {
	setSticky := false

	if !supportsCreateWithStickyBit && flag&os.O_CREATE != 0 && perm&os.ModeSticky != 0 {
		if _, err = os.Stat(name); os.IsNotExist(err) {
			setSticky = true
		}
	}

	// explicitly allow leaking file discriptor to subprocesses
	flag &= ^syscall.O_CLOEXEC
	mode := syscallMode(perm)

	for fd, err = syscall.Open(name, flag, mode); err != nil; fd, err = syscall.Open(name, flag, mode) {
		// On OS X, sigaction(2) doesn't guarantee that SA_RESTART will cause open(2) to be restarted
		// for regular files. This is easy to reproduce on fuse file systems (see https://golang.org/issue/11180).
		if runtime.GOOS == "darwin" && err == syscall.EINTR {
			continue
		} else {
			err = &os.PathError{Op: "open", Path: name, Err: err}
			return
		}
	}

	// open(2) itself won't handle the sticky bit on *BSD and Solaris
	if setSticky {
		if err = setStickyBit(name); err != nil {
			return
		}
	}

	f = os.NewFile(uintptr(fd), name)

	return f, fd, err
}
