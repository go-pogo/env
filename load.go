// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"io/fs"
	"os"
)

type Loader struct {
	fsys     fs.FS
	fallback Lookupper
	found    Map

	Workdir string
}

func Open(fsys fs.FS, path string) *Loader {
	if fsys == nil {
		fsys = os.DirFS(path)
		path = ""
	}

	l := NewLoader(fsys)
	l.Workdir = path
	l.TryFile(".env")
	l.TryFile(".env.local")
	return l
}

func OpenEnv(fsys fs.FS, path, environment string) *Loader {
	l := Open(fsys, path)
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

	defer fr.Close()
	return scanAll(NewScanner(fr), l.found, true)
}

func (l *Loader) TryFile(file string) bool {
	if l.File(file) != nil {
		return false
	}
	return true
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
		options: opts,
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

// type file struct {
// 	fsys     fs.FS
// 	path     string
// 	required bool
// }
//
// func (f *file) Configure(v interface{}) error {
// 	fr, err := f.fsys.Open(f.path)
// 	if err != nil {
// 		if f.required || !errors.Is(err, fs.ErrNotExist) {
// 			return errors.WithStack(err)
// 		}
// 		// the DecoderFunc must figure out what to do with a nil io.Reader
// 		return Decode(nil, v)
// 	}
// 	err = Decode(fr, v)
// 	errors.Append(&err, fr.Close())
// 	return err
// }
//
// func absPath(path string) string {
// 	if !filepath.IsAbs(path) {
// 		if cwd, err := os.Getwd(); err == nil {
// 			path = filepath.Join(cwd, path)
// 		}
// 	}
// 	return path
// }
