package tmpfile_test

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"testing"

	"github.com/triarius/tmpfile"
)

func TestNew(t *testing.T) {
	f, err := tmpfile.New("", "tempTest", "")
	if err != nil {
		t.Fatalf("Error creating tempfile: %+v\n", err)
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

	for i := 0; i < 10; i++ {
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

	for i := 0; i < 10; i++ {
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
		fmt.Println(text)
	}
}
