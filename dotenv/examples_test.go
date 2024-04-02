// Copyright (c) 2024, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"embed"
	"github.com/davecgh/go-spew/spew"
	"github.com/go-pogo/env"
	"time"
)

//go:embed example/*
var fsys embed.FS

// This example reads .env files from the "example" directory and decodes the
// found variables into a struct.
func ExampleRead() {
	type Config struct {
		Foo     string
		Timeout time.Duration `default:"10s"`
	}

	var conf Config
	if err := env.NewDecoder(ReadFS(fsys, "example", None)).Decode(&conf); err != nil {
		panic(err)
	}

	spew.Dump(conf)
	// Output:
	// (dotenv.Config) {
	//  Foo: (string) (len=3) "bar",
	//  Timeout: (time.Duration) 2s
	// }
}
