// Package rotationwatcher detects log file rotation by monitoring inode
// changes or file truncation, allowing logslice to reopen the file and
// continue reading from the new log without missing entries.
package rotationwatcher

import (
	"errors"
	"os"
	"sync"
	"time"
)

// Watcher monitors a file for rotation events.
type Watcher struct {
	mu       sync.Mutex
	path     string
	pollInterval time.Duration
	inode    uint64
	size     int64
	rotated  bool
	stop     chan struct{}
}

// New creates a Watcher for the given file path and poll interval.
// Returns an error if the file cannot be stat'd initially.
func New(path string, pollInterval time.Duration) (*Watcher, error) {
	if path == "" {
		return nil, errors.New("rotationwatcher: path must not be empty")
	}
	if pollInterval <= 0 {
		return nil, errors.New("rotationwatcher: poll interval must be positive")
	}
	w := &Watcher{
		path:         path,
		pollInterval: pollInterval,
		stop:         make(chan struct{}),
	}
	if err := w.snapshot(); err != nil {
		return nil, err
	}
	return w, nil
}

// Start begins polling for rotation in the background.
func (w *Watcher) Start() {
	go func() {
		ticker := time.NewTicker(w.pollInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				_ = w.check()
			case <-w.stop:
				return
			}
		}
	}()
}

// Stop halts background polling.
func (w *Watcher) Stop() {
	close(w.stop)
}

// Rotated reports whether a rotation has been detected since the last Reset.
func (w *Watcher) Rotated() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return w.rotated
}

// Reset clears the rotation flag and re-snapshots the file.
func (w *Watcher) Reset() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	w.rotated = false
	return w.snapshot()
}

func (w *Watcher) snapshot() error {
	info, err := os.Stat(w.path)
	if err != nil {
		return err
	}
	w.inode = inode(info)
	w.size = info.Size()
	return nil
}

func (w *Watcher) check() error {
	info, err := os.Stat(w.path)
	if err != nil {
		return err
	}
	w.mu.Lock()
	defer w.mu.Unlock()
	if inode(info) != w.inode || info.Size() < w.size {
		w.rotated = true
	}
	w.size = info.Size()
	return nil
}
