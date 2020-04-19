// +build aix darwin dragonfly freebsd linux netbsd openbsd solaris

package tmpfile

import "strconv"

func procFilePath(fd int) string {
	return "/proc/self/fd/" + strconv.Itoa(fd)
}
