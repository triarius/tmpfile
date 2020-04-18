// +build aix darwin dragonfly freebsd js,wasm linux netbsd openbsd solaris

package tmpfile

import (
	"os"
	"runtime"
	"syscall"
)

// OpenFile opens a file an return a struct containing the file and it's file descriptor
func OpenFile(name string, flag int, perm os.FileMode) (fwd FileWithDescriptor, err error) {
	setSticky := false
	if !supportsCreateWithStickyBit && flag&os.O_CREATE != 0 && perm&os.ModeSticky != 0 {
		if _, err = os.Stat(name); os.IsNotExist(err) {
			setSticky = true
		}
	}

	for {
		fwd.Fd, err = syscall.Open(name, flag|syscall.O_CLOEXEC, syscallMode(perm))
		if err == nil {
			break
		}

		// On OS X, sigaction(2) doesn't guarantee that SA_RESTART will cause
		// open(2) to be restarted for regular files. This is easy to reproduce on
		// fuse file systems (see https://golang.org/issue/11180).
		if runtime.GOOS == "darwin" && err == syscall.EINTR {
			continue
		}

		err = &os.PathError{Op: "open", Path: name, Err: err}
		return
	}

	// open(2) itself won't handle the sticky bit on *BSD and Solaris
	if setSticky {
		setStickyBit(name)
	}

	// There's a race here with fork/exec, which we are
	// content to live with. See ../syscall/exec_unix.go.
	if !supportsCloseOnExec {
		syscall.CloseOnExec(fwd.Fd)
	}

	fwd.File = os.NewFile(uintptr(fwd.Fd), name)

	return
}
