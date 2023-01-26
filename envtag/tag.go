// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"reflect"
	"strings"
)

const (
	EnvKey     = "env"
	DefaultKey = "default"

	ignore   = "-"
	noPrefix = "noprefix"
)

type Tag struct {
	Name     string
	Default  string
	Ignore   bool
	NoPrefix bool
}

func (t Tag) IsEmpty() bool {
	return t.Name == "" && !t.Ignore && !t.NoPrefix
}

func ParseTag(str string) Tag {
	var t Tag
	parse(&t, str)
	return t
}

func ParseStructTag(tag reflect.StructTag) Tag {
	var t Tag
	if str, found := tag.Lookup(EnvKey); found {
		parse(&t, str)
	}
	if !t.Ignore {
		if def := tag.Get(DefaultKey); def != "" {
			t.Default = def
		}
	}
	return t
}

func parse(tag *Tag, str string) {
	if str == "" {
		return
	}

	if str == ignore || strings.HasPrefix(str, ignore) {
		tag.Ignore = true
		return
	}

	split := strings.Split(str, ",")
	tag.Name = split[0]

	for _, s := range split[1:] {
		s = strings.TrimSpace(s)
		switch s {
		case noPrefix:
			tag.NoPrefix = true
		}
	}
}
