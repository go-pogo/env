// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package env

import (
	"fmt"
	"github.com/davecgh/go-spew/spew"
	"net"
	"net/url"
	"time"
)

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
		panic(fmt.Sprintf("%+v", err))
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
