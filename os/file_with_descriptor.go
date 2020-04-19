package os

import "os"

// FileWithDescriptor combines a file with its file descriptor
type FileWithDescriptor struct {
	*os.File
	Fd int
}
