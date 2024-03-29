// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package envfile provides tools to read and load environment variables from
// files.
package envfile

const (
	panicNilFile = "envfile: file must not be nil"
	panicNilFsys = "envfile: fs.FS must not be nil"
)
