// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"strings"
	"sync"
)

type Normalizer interface {
	Normalize(str, prefix string) string
}

type NormalizerFunc func(fieldName, prefix string) string

func (fn NormalizerFunc) Normalize(fieldName, prefix string) string {
	return fn(fieldName, prefix)
}

type Options struct {
	// EnvKey is used to look up the env tag string from a [reflect.StructTag].
	EnvKey string
	// DefaultKey is used to look up the default value string from a
	// [reflect.StructTag].
	DefaultKey string
	// Normalizer is used to normalize a [reflect.StructField]'s name
	// when no name is provided in the env tag string.
	Normalizer Normalizer
	// StrictTags ignores fields that do not have an env tag when set to true.
	StrictTags bool
}

// DefaultOptions returns an [Options] with default values.
func DefaultOptions() Options {
	return Options{
		EnvKey:     EnvKey,
		DefaultKey: DefaultKey,
		Normalizer: defaultNormalizer,
		StrictTags: false,
	}
}

var defaultNormalizer = new(FieldNameNormalizer)

// FieldNameNormalizer normalizes struct field names to a format that is
// commonly used in environment variables, e.g. MyFieldName becomes
// MY_FIELD_NAME.
type FieldNameNormalizer struct {
	normalized sync.Map
}

func (n *FieldNameNormalizer) Normalize(fieldName, prefix string) string {
	var buf strings.Builder
	buf.Grow(len(fieldName) + len(prefix) + 3)

	if len(prefix) > 0 {
		if n.isNormalized(prefix) {
			buf.WriteString(prefix)
		} else {
			n.normalize(&buf, prefix)
		}
		buf.WriteRune('_')
	}

	if n.isNormalized(fieldName) {
		buf.WriteString(fieldName)
	} else {
		n.normalize(&buf, fieldName)
	}
	return buf.String()
}

func (n *FieldNameNormalizer) isNormalized(str string) bool {
	if res, ok := n.normalized.Load(str); ok {
		return res.(bool)
	}

	for i := 0; i < len(str); i++ {
		if isUpper(str[i]) || isDigit(str[i]) || str[i] == '_' {
			continue
		}
		n.normalized.Store(str, false)
		return false
	}
	n.normalized.Store(str, true)
	return true
}

func (n *FieldNameNormalizer) normalize(buf *strings.Builder, str string) {
	var prev byte // previous character to write
	var lastUpper int
	var lastLower int
	var lastDigits int

	for i := 0; i < len(str); i++ {
		cur := str[i] // current character

		switch true {
		case isUpper(cur):
			if i > 2 && (lastLower == i-2 || lastDigits == i-2) {
				buf.WriteRune('_')
			}
			lastUpper = i

		case isLower(cur):
			if i > 1 && (lastUpper == i-1 || lastDigits == i-1) {
				buf.WriteRune('_')
			}
			// convert lowercase to uppercase
			cur += 'A'
			cur -= 'a'
			lastLower = i

		case isDigit(cur):
			if i > 2 && (lastUpper-lastLower == 1 || lastLower-lastUpper == 1) {
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
			buf.WriteByte(prev)
		}
		prev = cur
	}
	buf.WriteByte(prev)
}

func isUpper(c byte) bool { return c >= 'A' && c <= 'Z' }
func isLower(c byte) bool { return c >= 'a' && c <= 'z' }
func isDigit(c byte) bool { return c >= '0' && c <= '9' }
