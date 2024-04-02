// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package osfs

import (
	"io/fs"
	"os"
	"path/filepath"
)

var _ fs.FS = (*FS)(nil)

// FS is a fs.FS compatible wrapper around os.Open.
type FS struct{}

func (FS) Open(name string) (fs.File, error) { return os.Open(name) }

func (FS) JoinFilePath(elem ...string) string { return filepath.Join(elem...) }
