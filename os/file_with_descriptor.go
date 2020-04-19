package os

import "os"

// FileWithDescriptor combines a file with its file descriptor
type FileWithDescriptor struct {
	*os.File
	fd int
}

func (f FileWithDescriptor) Fd() int {
	return f.fd
}

func (f FileWithDescriptor) ProcFilePath() string {
	return f.procFilePath()
}
