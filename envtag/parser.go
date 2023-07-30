// Copyright (c) 2023, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"reflect"
	"strings"
	"sync"
)

const (
	EnvKey     = "env"
	DefaultKey = "default"
)

type StructParser struct {
	once       sync.Once
	EnvKey     string
	DefaultKey string

	// TagsOnly ignores fields that do not have an `env` tag when set to true.
	TagsOnly bool
}

func (p *StructParser) envKey() string {
	p.once.Do(func() {
		if p.EnvKey == "" && p.DefaultKey == "" {
			p.EnvKey = EnvKey
			p.DefaultKey = DefaultKey
		}
	})

	if p.EnvKey != "" {
		return p.EnvKey
	}
	return EnvKey
}

func (p *StructParser) ParseStructField(field reflect.StructField) Tag {
	tag := p.ParseStructTag(field.Tag)
	if (tag.IsEmpty() && p.TagsOnly) || tag.Ignore {
		return tag
	}
	if tag.Name == "" {
		tag.Name = strings.ToUpper(field.Name)
	}
	return tag
}

func (p *StructParser) ParseStructTag(tag reflect.StructTag) Tag {
	var t Tag
	if str, found := tag.Lookup(p.envKey()); found {
		parse(&t, str)
	}
	if p.DefaultKey != "" {
		t.Default = tag.Get(p.DefaultKey)
	}
	return t
}
