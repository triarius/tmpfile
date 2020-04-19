// +build aix solaris

package os

import (
	"os"
	"strconv"
)

func (f FileWithDescriptor) procFilePath() string {
	return "/proc/" + strconv.Itoa(os.Getpid()) + "/fd/" + strconv.Itoa(f.fd)
}
