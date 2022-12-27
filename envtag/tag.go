// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import "strings"

const (
	Key      = "env"
	ignore   = "-"
	noPrefix = "noprefix"
)

type Tag struct {
	Name     string
	Ignore   bool
	NoPrefix bool
}

func ParseTag(tag string) Tag {
	var t Tag
	if tag == "" {
		return t
	}

	if tag == ignore || strings.HasPrefix(tag, ignore) {
		t.Ignore = true
		return t
	}

	split := strings.Split(tag, ",")
	t.Name = split[0]

	for _, s := range split[1:] {
		s = strings.TrimSpace(s)
		switch s {
		case noPrefix:
			t.NoPrefix = true
		}
	}

	return t
}
