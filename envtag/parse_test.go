// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseTag(t *testing.T) {
	tests := map[string]Tag{
		"":                 {},
		"-":                {Ignore: true},
		"-,noprefix":       {Ignore: true},
		"foo":              {Name: "foo"},
		"foo,inline":       {Name: "foo", Inline: true},
		"foo,noprefix":     {Name: "foo", NoPrefix: true},
		"FOOBAR":           {Name: "FOOBAR"},
		"noprefix":         {Name: "noprefix"},
		",noprefix":        {NoPrefix: true},
		",inline,noprefix": {Inline: true, NoPrefix: true},
	}

	for tag, want := range tests {
		t.Run(tag, func(t *testing.T) {
			assert.Equal(t, want, ParseTag(tag))
		})
	}
}

func TestParseStructField(t *testing.T) {
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
		foo string //nolint:unused
	}

	tests := map[string]struct {
		opts  *Options
		field reflect.StructField
		want  Tag
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
			want:  Tag{Ignore: true},
		},
		"tags only": {
			opts:  &Options{TagsOnly: true},
			field: reflect.TypeOf(fixtureBasic{}).Field(0),
			want:  Tag{Ignore: true},
		},
		"unexported": {
			field: reflect.TypeOf(fixtureUnexported{}).Field(0),
			want:  Tag{Ignore: true},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if tc.opts == nil {
				tc.opts = new(Options)
				tc.opts.Defaults()
			}

			have := ParseStructField(*tc.opts, tc.field)
			assert.Equal(t, tc.want, have)
		})
	}
}
