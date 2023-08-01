// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"reflect"
	"strings"
)

const (
	ignore   = "-"
	inline   = "inline"   // inline means ignore name of parent field
	noPrefix = "noprefix" // noprefix means ignore all parent field names

	EnvKey     = "env"
	DefaultKey = "default"
)

type Options struct {
	EnvKey     string
	DefaultKey string

	// TagsOnly ignores fields that do not have an `env` tag when set to true.
	TagsOnly bool
}

// Defaults sets the default values for Options.
func (o *Options) Defaults() {
	o.EnvKey = EnvKey
	o.DefaultKey = DefaultKey
	o.TagsOnly = false
}

// ParseTag parses str into a Tag.
func ParseTag(str string) Tag {
	var t Tag
	parse(&t, str)
	return t
}

func ParseStructField(opts Options, field reflect.StructField) Tag {
	var tag Tag
	if !field.IsExported() {
		tag.Ignore = true
		return tag
	}

	if str, found := field.Tag.Lookup(opts.EnvKey); found {
		parse(&tag, str)
	} else if !found && opts.TagsOnly {
		tag.Ignore = true
		return tag
	}

	if tag.Ignore {
		return tag
	}
	if tag.Name == "" {
		tag.Name = strings.ToUpper(field.Name)
	}
	if opts.DefaultKey != "" {
		tag.Default = field.Tag.Get(opts.DefaultKey)
	}
	return tag
}

func parse(tag *Tag, str string) {
	if str == "" {
		return
	}
	if str == ignore {
		tag.Ignore = true
		return
	}

	split := strings.Split(str, ",")
	if split[0] == ignore {
		tag.Ignore = true
		return
	}

	tag.Name = split[0]
	for _, s := range split[1:] {
		s = strings.TrimSpace(s)
		switch s {
		case inline:
			tag.Inline = true
		case noPrefix:
			tag.NoPrefix = true
		}
	}
}
