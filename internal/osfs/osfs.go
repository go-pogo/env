package osfs

import (
	"io/fs"
	"os"
)

var _ fs.FS = (*FS)(nil)

// FS is a fs.FS compatible wrapper around os.Open.
type FS struct{}

func (FS) Open(name string) (fs.File, error) { return os.Open(name) }
