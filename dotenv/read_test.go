// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"github.com/go-pogo/env"
	"github.com/stretchr/testify/assert"
	"testing"
	"testing/fstest"
)

func TestReader_Lookup(t *testing.T) {
	fsys := fstest.MapFS{
		".env.prod": &fstest.MapFile{
			Data: []byte("FOO=BAR\nQUX=XOO"),
		},
		".env.dev": &fstest.MapFile{
			Data: []byte("FOO=BAZ"),
		},
		".env.dev.local": &fstest.MapFile{
			Data: []byte("FOO=BAZZZ"),
		},
	}

	tests := map[ActiveEnvironment]map[string]struct {
		Key     string
		Want    env.Value
		WantErr error
	}{
		None: {
			"none loaded": {
				Key:     "FOO",
				WantErr: ErrNoFilesLoaded,
			},
		},
		Production: {
			"has FOO": {
				Key:  "FOO",
				Want: "BAR",
			},
			"has QUX": {
				Key:  "QUX",
				Want: "XOO",
			},
		},
		Development: {
			"has FOO": {
				Key:  "FOO",
				Want: "BAZZZ",
			},
			"QUX not found": {
				Key:     "QUX",
				WantErr: env.ErrNotFound,
			},
		},
	}
	for ae, tt := range tests {
		t.Run(ae.String(), func(t *testing.T) {
			for name, tc := range tt {
				t.Run(name, func(t *testing.T) {
					have, haveErr := ReadFS(fsys, "", ae).Lookup(tc.Key)
					assert.Equal(t, tc.Want, have)
					if tc.WantErr == nil {
						assert.NoError(t, haveErr)
					} else {
						assert.ErrorIs(t, haveErr, tc.WantErr)
					}
				})
			}
		})
	}
}

func TestReader_Environ(t *testing.T) {
	fsys := fstest.MapFS{
		".env.prod": &fstest.MapFile{
			Data: []byte("FOO=BAR\nQUX=XOO"),
		},
		".env.dev": &fstest.MapFile{
			Data: []byte("FOO=BAZ"),
		},
		".env.dev.local": &fstest.MapFile{
			Data: []byte("FOO=BAZZZ"),
		},
	}

	tests := map[ActiveEnvironment]struct {
		Want    env.Map
		WantErr error
	}{
		None: {
			WantErr: ErrNoFilesLoaded,
		},
		Production: {
			Want: env.Map{
				"FOO": "BAR",
				"QUX": "XOO",
			},
		},
		Development: {
			Want: env.Map{"FOO": "BAZZZ"},
		},
	}
	for ae, tc := range tests {
		t.Run(ae.String(), func(t *testing.T) {
			have, haveErr := ReadFS(fsys, "", ae).Environ()
			assert.Equal(t, tc.Want, have)
			if tc.WantErr == nil {
				assert.NoError(t, haveErr)
			} else {
				assert.ErrorIs(t, haveErr, tc.WantErr)
			}
		})
	}
}
