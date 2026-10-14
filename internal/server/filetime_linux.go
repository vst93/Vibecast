//go:build linux

package server

import (
	"os"
	"syscall"
	"time"
)

// fileBirthTime returns the file's creation time (birth time).
// Returns the zero time if the filesystem doesn't support it.
func fileBirthTime(_ string, info os.FileInfo) time.Time {
	st, ok := info.Sys().(*syscall.Stat_t)
	if !ok {
		return time.Time{}
	}
	return time.Unix(st.Ctim.Sec, st.Ctim.Nsec)
}
