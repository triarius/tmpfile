package tmpfile

import (
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"sync"
	"time"

	tmpos "github.com/triarius/tmpfile/os"
)

const (
	totalRetry     = 10000
	reseedTreshold = 10
	tempFileMode   = 0600
)

// While obtaining random numbers from the source is concurency safe, reseeding is not:
// > Seed should not be called concurrently with any other Rand Method
// https://golang.org/src/math/rand/rand.go?s=2477:2508#L64
// Therefore, we need to mutex all Rand methods
var rngMutex sync.Mutex

func reseed() {
	rngMutex.Lock()
	rand.Seed(time.Now().UnixNano() + int64(os.Getpid()))
	rngMutex.Unlock()
}

func nextRandom() string {
	rngMutex.Lock()
	r := rand.Uint64()
	rngMutex.Unlock()
	return strconv.FormatUint(r, 16)
}

// New creates a temporary file and returns a file descriptor and name
// If dir is empty, it will use the os's Temp directory
func New(dir, prefix, suffix string) (fwd tmpos.FileWithDescriptor, err error) {
	if dir == "" {
		dir = os.TempDir()
	}

	nConflicts := 0
	for i := 0; i < totalRetry; i++ {
		name := filepath.Join(dir, prefix+nextRandom()+suffix)
		fwd, err = tmpos.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, tempFileMode)
		if os.IsExist(err) {
			if nConflicts++; nConflicts > reseedTreshold {
				reseed()
			}
			continue
		}
		break
	}
	return
}
