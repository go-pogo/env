// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFieldNameNormalizer(t *testing.T) {
	tests := map[string]string{
		"Foo":          "FOO",
		"FOO":          "FOO",
		"FOO123":       "FOO123",
		"FOOBar":       "FOO_BAR",
		"FooBar":       "FOO_BAR",
		"FooBAR":       "FOO_BAR",
		"Foo_BAR":      "FOO_BAR",
		"Foo123BAR":    "FOO123_BAR",
		"FooBarBaz":    "FOO_BAR_BAZ",
		"FOOBarBaz":    "FOO_BAR_BAZ",
		"FooBARBaz":    "FOO_BAR_BAZ",
		"FooBAR123Baz": "FOO_BAR123_BAZ",
		"FooBarBAZ":    "FOO_BAR_BAZ",
		"ApiV1":        "API_V1",
		"ApiV1Test":    "API_V1_TEST",

		// below special cases probably result in an unexpected normalized name
		// but it's not worth the effort to fix this right now. a workaround
		// would be to set the right name using an env tag on the struct field.
		"SomeGraphQL": "SOME_GRAPH_QL", // SOME_GRAPHQL
		"PostgreSQL":  "POSTGRE_SQL",   // POSTGRESQL
	}

	var n FieldNameNormalizer
	for name, want := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, want, n.Normalize(name, ""))
		})
	}
}
