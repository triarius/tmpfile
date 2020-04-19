// +build darwin

package os

// On OS X, sigaction(2) doesn't guarantee that SA_RESTART will cause open(2) to be restarted for
// regular files. This is easy to reproduce on fuse file systems (see https://golang.org/issue/11180).
const retrySAResart = true
