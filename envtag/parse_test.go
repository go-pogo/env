// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package envtag

import (
	"github.com/stretchr/testify/assert"
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
