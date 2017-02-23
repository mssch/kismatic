package check

import (
	"fmt"
	"syscall"
)

// FreeSpaceCheck checks the available disk space on a path
type FreeSpaceCheck struct {
	MinimumBytes uint64
	Path         string
}

// Check returns true if the path has enough free space. Otherwise return false.
func (c FreeSpaceCheck) Check() (bool, error) {
	var stat syscall.Statfs_t
	if err := syscall.Statfs(c.Path, &stat); err != nil {
		return false, fmt.Errorf("failed to check free space at path %s: %v", c.Path, err)
	}
	// Available blocks * size per block = available space in bytes
	availableBytes := stat.Bavail * uint64(stat.Bsize)
	return availableBytes >= c.MinimumBytes, nil
}
