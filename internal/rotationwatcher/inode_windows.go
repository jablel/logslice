//go:build windows

package rotationwatcher

import "os"

// inode returns 0 on Windows where inode tracking is not supported.
// Rotation is still detected via file size shrinkage.
func inode(info os.FileInfo) uint64 {
	return 0
}
