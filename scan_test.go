// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestNewScanner(t *testing.T) {
	t.Run("nil reader", func(t *testing.T) {
		assert.PanicsWithValue(t, panicNilReader, func() { NewScanner(nil) })
	})
}

func TestScanner(t *testing.T) {
	tests := map[string]struct {
		input string
		want  []string
	}{
		"basic": {
			input: "some=value\n",
			want:  []string{"some=value"},
		},
		"multiple lines": {
			input: "some=value\n\n#comment\nfoo=bar",
			want:  []string{"some=value", "foo=bar"},
		},
		"comments": {
			input: "\t#comment\n   #   another comment    \n#final comment",
		},
	}
	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			scan := NewScanner(strings.NewReader(tc.input))
			have := make([]string, 0, len(tc.want))
			for scan.Scan() {
				have = append(have, scan.Text())
			}
			if tc.want == nil {
				assert.Empty(t, have)
			} else {
				assert.Equal(t, tc.want, have)
			}
		})
	}
}

func TestScanLines(t *testing.T) {
	tests := []struct {
		name, input string
		wantErr     error
		wantTok     []byte
		wantAdv     int
	}{
		{
			name: "empty",
		},
		{
			input:   "some=value",
			wantAdv: 10,
			wantTok: []byte("some=value"),
		},
		{
			input:   "  a line with whitespace  ",
			wantAdv: 26,
			wantTok: []byte("a line with whitespace"),
		},
		{
			input:   "a line\nanother line",
			wantAdv: 7,
			wantTok: []byte("a line"),
		},
		{
			input:   "# some comment",
			wantAdv: 14,
			wantTok: []byte{},
		},
		{
			input:   "\t #another comment",
			wantAdv: 18,
			wantTok: []byte{},
		},
	}

	for _, tc := range tests {
		if tc.name == "" {
			tc.name = tc.input
		}

		t.Run(tc.name, func(t *testing.T) {
			haveAdv, haveTok, haveErr := ScanLines([]byte(tc.input), true)
			assert.Equal(t, tc.wantAdv, haveAdv)
			assert.Equal(t, tc.wantTok, haveTok)
			assert.Equal(t, tc.wantErr, haveErr)
		})
	}
}
