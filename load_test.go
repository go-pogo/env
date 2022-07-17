// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"io/fs"
	"strings"
)

func ExampleOpen() {
	fsys := dummyFS{
		".env":       dummyFile{strings.NewReader("local=")},
		".env.local": dummyFile{strings.NewReader("local=some local value")},
	}

	envs, err := Open(fsys, "").DefaultDecoder().Map()
	if err != nil {
		panic(err)
	}

	fmt.Println(envs)
	// Output: map[local:some local value]
}

func ExampleOpenEnv() {
	fsys := dummyFS{
		".env":       dummyFile{strings.NewReader("local=")},
		".env.local": dummyFile{strings.NewReader("local=some local value")},
		".env.prod":  dummyFile{strings.NewReader("local=value is now `prod`")},
	}

	// for this example we're using an embedded filesystem.
	// in real life you'll probably read a .env relative to the
	// current working directory, or executable's path.
	envs, err := OpenEnv(fsys, "", "prod").DefaultDecoder().Map()
	if err != nil {
		panic(err)
	}

	fmt.Println(envs)
	// Output: map[local:value is now `prod`]
}

type dummyFS map[string]dummyFile

func (f dummyFS) Open(name string) (fs.File, error) {
	file, ok := f[name]
	if !ok {
		return nil, fs.ErrNotExist
	}
	return &file, nil
}

type dummyFile struct {
	*strings.Reader
}

func (f *dummyFile) Stat() (fs.FileInfo, error) { return nil, nil }
func (f *dummyFile) Close() error               { return nil }
