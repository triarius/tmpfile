package os

import "os"

// setStickyBit adds ModeSticky to the permission bits of path, non atomic.
func setStickyBit(name string) error {
	fi, err := os.Stat(name)
	if err != nil {
		return err
	}

	return os.Chmod(name, fi.Mode()|os.ModeSticky)
}
