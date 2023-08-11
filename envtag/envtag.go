// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Package envtag provides tools to parse tags from strings or struct fields.
package envtag

type Tag struct {
	// Name of the tag which will be used to construct the environment
	// variable's full name.
	Name string
	// Default is an optional default value for the field.
	Default string
	// Ignore indicates the field should be ignored.
	Ignore bool
	// Inline indicates the field will be "flattened" when it is a struct. This
	// means it's child fields will be treated as if they are part of the parent
	// struct.
	Inline bool
	// Include indicates the structs child fields should be included, even if
	// they do not have an env tag and Options.StrictTag is set to true.
	Include bool
	// NoPrefix indicates the tag's name should not be prefixed with any of it's
	// parent's names.
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
