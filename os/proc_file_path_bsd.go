// +build darwin dragonfly freebsd netbsd openbsd

package os

import "strconv"

func (f FileWithDescriptor) procFilePath() string {
	return "/dev/fd/" + strconv.Itoa(f.fd)
}
