// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"io/fs"
	"os"

	"github.com/go-pogo/errors"
)

type Loader struct {
	fsys     fs.FS
	fallback Lookupper
	found    Map

	Workdir string
}

func Open(fsys fs.FS, dir string) *Loader {
	if fsys == nil {
		fsys = os.DirFS(dir)
		dir = ""
	}

	l := NewLoader(fsys)
	l.Workdir = dir
	l.TryFile(".env")
	l.TryFile(".env.local")
	return l
}

func OpenEnv(fsys fs.FS, dir, environment string) *Loader {
	l := Open(fsys, dir)
	if environment != "" {
		l.TryFile(".env." + environment)
		l.TryFile(".env." + environment + ".local")
	}
	return l
}

func NewLoader(fsys fs.FS) *Loader {
	return &Loader{
		fsys:  fsys,
		found: make(Map, 4),
	}
}

func (l *Loader) File(file string) error {
	fr, err := l.fsys.Open(l.Workdir + file)
	if err != nil {
		return err
	}

	defer errors.AppendFunc(&err, fr.Close)
	return scanAll(NewScanner(fr), l.found, true)
}

func (l *Loader) TryFile(file string) bool {
	return l.File(file) == nil
}

// Fallback Lookupper which is called when a key cannot be found in any of the
// specified files.
func (l *Loader) Fallback(fallback Lookupper) {
	if l == fallback {
		panic(panicSelfAsFallback)
	}
	l.fallback = fallback
}

func (l *Loader) Lookup(key string) (Value, bool) {
	if val, ok := l.found[key]; ok {
		return val, true
	}
	if l.fallback != nil {
		return l.fallback.Lookup(key)
	}
	return "", false
}

func (l *Loader) Decoder(opts Option) *Decoder {
	return &Decoder{
		scanner: new(nilScanner),
		found:   l.found,
		opts:    opts,
	}
}

func (l *Loader) DefaultDecoder() *Decoder {
	dec := l.Decoder(DefaultOptions)
	if l.fallback == nil {
		setDefaultFallback(dec)
	} else {
		dec.Fallback(l.fallback)
	}
	return dec
}
