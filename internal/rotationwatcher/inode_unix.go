//go:build !windows

package rotationwatcher

import (
	"os"
	"syscall"
)

// inode returns the inode number of the file described by info.
func inode(info os.FileInfo) uint64 {
	if stat, ok := info.Sys().(*syscall.Stat_t); ok {
		return stat.Ino
	}
	return 0
}
