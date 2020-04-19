package tmpfile_test

import (
	"io/ioutil"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/triarius/tmpfile"
)

func TestNew(t *testing.T) {
	f, err := tmpfile.New("", "tempTest", "")
	if err != nil {
		t.Fatalf("Error creating tempfile: %+v\n", err)
	}

	if strings.HasPrefix(f.Name(), "/tmp") {
		t.Fatalf("File handle is not in proc FS: name: %s", f.Name())
	}

	if _, err = f.WriteString("test"); err != nil {
		t.Fatalf("Could not write to %+v", f)
	}

	content, err := ioutil.ReadFile(f.Name())
	if err != nil {
		t.Fatalf("could not read: " + f.Name())
	}

	f.Close()

	if string(content) != "test" {
		t.Fatalf("Read content mismatch, expected: %s, found: %s", "test", string(content))
	}
}

func TestChans(t *testing.T) {
	fChan := make(chan *os.File)

	const numFiles = 10

	for i := 0; i < numFiles; i++ {
		go func(i int, fChan chan *os.File) {
			var f *os.File

			defer func() { fChan <- f }()

			f, err := tmpfile.New("", "tempTest", "")
			if err != nil {
				return
			}

			if _, err = f.WriteString(strconv.Itoa(i)); err != nil {
				return
			}
		}(i, fChan)
	}

	var results [numFiles]bool

	for i := 0; i < numFiles; i++ {
		f := <-fChan
		if f == nil {
			t.Fatalf("No file received.")
		}

		content, err := ioutil.ReadFile(f.Name())
		if err != nil {
			t.Fatalf("could not read: " + f.Name())
		}

		f.Close()

		text := string(content)
		pos, err := strconv.Atoi(text)
		if err != nil {
			t.Fatalf("Content not a number: %s", content)
		}
		results[pos] = true
	}

	for i := 0; i < numFiles; i++ {
		if !results[i] {
			t.Errorf("File %d was not read back.", i)
		}
	}
}
