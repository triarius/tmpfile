// +build linux

package os

import "strconv"

func (f FileWithDescriptor) procFilePath() string {
	return "/proc/self/fd/" + strconv.Itoa(f.fd)
}
