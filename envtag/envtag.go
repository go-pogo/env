// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

type Tag struct {
	Name     string
	Default  string
	Ignore   bool
	Inline   bool
	NoPrefix bool
}

// IsEmpty indicates if Tag is considered empty.
func (t Tag) IsEmpty() bool {
	return t.Name == "" && !t.Ignore && !t.Inline && !t.NoPrefix
}

// ShouldIgnore indicates if Tag should be ignored.
func (t Tag) ShouldIgnore() bool {
	return t.Name == "" || t.Ignore
}
