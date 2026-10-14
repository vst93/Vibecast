//go:build !linux && !darwin

package server

import (
	"os"
	"time"
)

// fileBirthTime is a fallback for platforms without birth time support
// in syscall.Stat_t. The created column falls back to the mod time.
func fileBirthTime(_ string, _ os.FileInfo) time.Time {
	return time.Time{}
}
