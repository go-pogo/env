// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

import (
	"flag"
	"io"
)

func GetEnvironment(args []string) ActiveEnvironment {
	fs := flag.NewFlagSet("dotenv", flag.ContinueOnError)
	fs.SetOutput(io.Discard)

	var e string
	fs.StringVar(&e, "active-env", "", "")
	_ = fs.Parse(args)
	return ActiveEnvironment(e)
}

func GetEnvironmentOrDefault(args []string, def ActiveEnvironment) ActiveEnvironment {
	if e := GetEnvironment(args); e != "" {
		return e
	}
	return def
}
