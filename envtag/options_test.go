// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeFieldName(t *testing.T) {
	tests := map[string]string{
		"Foo":          "FOO",
		"FOO":          "FOO",
		"FOO123":       "FOO_123",
		"FOOBar":       "FOO_BAR",
		"FooBar":       "FOO_BAR",
		"FooBAR":       "FOO_BAR",
		"Foo_BAR":      "FOO_BAR",
		"Foo123BAR":    "FOO_123_BAR",
		"FooBarBaz":    "FOO_BAR_BAZ",
		"FOOBarBaz":    "FOO_BAR_BAZ",
		"FooBARBaz":    "FOO_BAR_BAZ",
		"FooBAR123Baz": "FOO_BAR_123_BAZ",
		"FooBarBAZ":    "FOO_BAR_BAZ",
	}
	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, want, NormalizeFieldName(name))
		})
	}
}
