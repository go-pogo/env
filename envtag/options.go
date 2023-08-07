// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"strings"
)

type NormalizeFunc func(str string) string

type Options struct {
	// EnvKey is used to lookup the env tag string from a reflect.StructTag.
	EnvKey string
	// DefaultKey is used to lookup the default value string from a
	// reflect.StructTag.
	DefaultKey string
	// Normalize is used to normalize a reflect.StructField's name
	// when no name is provided in the env tag string.
	Normalize NormalizeFunc
	// TagsOnly ignores fields that do not have an env tag when set to true.
	TagsOnly bool
}

// Defaults sets the default values for Options.
func (o *Options) Defaults() {
	o.EnvKey = EnvKey
	o.DefaultKey = DefaultKey
	o.Normalize = NormalizeFieldName
	o.TagsOnly = false
}

// NormalizeFieldName normalizes struct field names to a format that is commonly
// used in environment variables, eg. MyFieldName becomes MY_FIELD_NAME.
func NormalizeFieldName(fieldName string) string {
	var buf strings.Builder
	buf.Grow(len(fieldName) + 3)

	var prevChar byte // previous character to write
	var lastUpper int
	var lastLower int
	var lastDigits int

	for i := 0; i < len(fieldName); i++ {
		c := fieldName[i] // current character

		switch true {
		case c >= 'A' && c <= 'Z':
			if i > 2 && (lastLower == i-2 || lastDigits == i-2) {
				buf.WriteRune('_')
			}
			lastUpper = i

		case c >= 'a' && c <= 'z':
			if i > 1 && (lastUpper == i-1 || lastDigits == i-1) {
				buf.WriteRune('_')
			}
			// convert lowercase to uppercase
			c += 'A'
			c -= 'a'
			lastLower = i

		case c >= '0' && c <= '9':
			if i > 2 && (lastUpper == i-2 || lastLower == i-2) {
				buf.WriteRune('_')
			}
			lastDigits = i

			// ignore all other characters
		default:
			lastUpper += 1
			lastLower += 1
			lastDigits += 1
			continue
		}

		if i > 0 {
			buf.WriteByte(prevChar)
		}
		prevChar = c
	}
	buf.WriteByte(prevChar)
	return buf.String()
}