// +build linux

package tmpfile

import (
	"os"
	"strconv"

	"github.com/pkg/errors"
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
		err = errors.Wrapf(err, "Could not create temp file in: %s", dir)
		return
	}

	f = os.NewFile(uintptr(fd), "/dev/fd/"+strconv.Itoa(fd))

	return f, err
}
