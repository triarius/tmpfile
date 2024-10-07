//go:build linux
// +build linux

package tmpfile

import (
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sys/unix"
)

// New creates a temporary file and returns a file descriptor and name
// If dir is empty, it will use the os's Temp directory.
func New(dir string, leakToSubProc bool) (f *os.File, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	flags := unix.O_RDWR | unix.O_TMPFILE | unix.O_EXCL
	if !leakToSubProc {
		flags |= unix.O_CLOEXEC
	}

	fd, err := unix.Open(dir, flags, 0600)
	if err != nil {
		return nil, fmt.Errorf("Could not create temp file in: %s: %w", dir, err)
	}

	return os.NewFile(uintptr(fd), "/dev/fd/"+strconv.Itoa(fd)), nil
}
