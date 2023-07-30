// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"strings"
)

const (
	ignore   = "-"
	inline   = "inline"   // inline means ignore name of parent field
	noPrefix = "noprefix" // noprefix means ignore all parent field names
)

// ParseTag parses str into a Tag.
func ParseTag(str string) Tag {
	var t Tag
	parse(&t, str)
	return t
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
