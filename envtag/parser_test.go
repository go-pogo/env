// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestStructParser_ParseStructField(t *testing.T) {
	type fixtureBasic struct {
		Foo string
	}
	type fixtureNamed struct {
		Foo string `env:"BAR"`
	}
	type fixtureIgnore struct {
		Foo string `env:"-" default:"some value"`
	}
	type fixtureUnexported struct {
		foo string
	}

	tests := map[string]struct {
		parser *StructParser
		field  reflect.StructField
		want   Tag
	}{
		"basic": {
			field: reflect.TypeOf(fixtureBasic{}).Field(0),
			want:  Tag{Name: "FOO"},
		},
		"named": {
			field: reflect.TypeOf(fixtureNamed{}).Field(0),
			want:  Tag{Name: "BAR"},
		},
		"ignore": {
			field: reflect.TypeOf(fixtureIgnore{}).Field(0),
			want:  Tag{Ignore: true, Default: "some value"},
		},
		"tags only": {
			parser: &StructParser{TagsOnly: true},
			field:  reflect.TypeOf(fixtureBasic{}).Field(0),
			want:   Tag{},
		},
		"unexported": {
			field: reflect.TypeOf(fixtureUnexported{}).Field(0),
			want:  Tag{Name: "FOO"},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.parser == nil {
				tc.parser = new(StructParser)
			}

			have := tc.parser.ParseStructField(tc.field)
			assert.Equal(t, tc.want, have)
		})
	}
}
