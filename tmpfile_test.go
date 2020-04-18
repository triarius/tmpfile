package tmpfile_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"syscall"
	"testing"

	"github.com/triarius/tmpfile"
)

func testChans(t *testing.T) {
	fChan := make(chan *os.File)
	for i := 0; i < 10; i++ {
		go func(i int, fChan chan *os.File) {
			var f *os.File
			defer func() { fChan <- f }()

			f, err := tmpfile.TempFile("", "tempTest", "")
			if err != nil {
				fmt.Printf("Error creating tempfile: %e\n", err)
				return
			}

			syscall.Write(fd, []byte(fmt.Sprintf("Wrote to file: %d, name: %s", i, name)))

			fmt.Printf("%d: %s\n", fd, name)
		}(i, fdChan)
	}

	for i := 0; i < 10; i++ {
		//fd := <-fdChan

		content, err := ioutil.ReadFile(name)
		if err != nil {
			panic("could not read: " + name)
		}
		syscall.Close(fd)
		text := string(content)
		fmt.Println(text)
	}
}
