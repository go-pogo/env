// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"net"
	"net/url"
	"time"

	"github.com/davecgh/go-spew/spew"
)

// Below example demonstrates how to decode system environment variables into a
// struct.
func ExampleDecoder_Decode() {
	type Config struct {
		Foo     string
		Timeout time.Duration `default:"10s"`
	}

	var conf Config
	if err := NewDecoder(System()).Decode(&conf); err != nil {
		panic(err)
	}

	spew.Dump(conf)
	// Output:
	// (env.Config) {
	//  Foo: (string) "",
	//  Timeout: (time.Duration) 10s
	// }
}

func ExampleUnmarshal() {
	type Envs struct {
		Foo string
		Bar struct {
			Url url.URL
		} `env:",inline"`
		Timeout time.Duration `default:"10s"`
		Ip      net.IP
	}

	var data = `
FOO=bar
# ignore me
URL=http://example.com
IP=192.168.1.1`

	var envs Envs
	if err := Unmarshal([]byte(data), &envs); err != nil {
		panic(err)
	}

	spew.Dump(envs)
	// Output:
	// (env.Envs) {
	//  Foo: (string) (len=3) "bar",
	//  Bar: (struct { Url url.URL }) {
	//   Url: (url.URL) http://example.com
	//  },
	//  Timeout: (time.Duration) 10s,
	//  Ip: (net.IP) (len=16 cap=16) 192.168.1.1
	// }
}

func ExampleMarshal() {
	type Envs struct {
		Foo string
		Bar struct {
			Url url.URL `default:"https://example.com"`
		} `env:",inline"`
		Timeout time.Duration `default:"10s"`
		Ip      net.IP
	}

	b, err := Marshal(Envs{})
	if err != nil {
		panic(err)
	}

	fmt.Println(string(b))
	// Output:
	// FOO=
	// URL=https://example.com
	// TIMEOUT=10s
	// IP=
}
