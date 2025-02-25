env
===
[![Latest release][latest-release-img]][latest-release-url]
[![Build status][build-status-img]][build-status-url]
[![Go Report Card][report-img]][report-url]
[![Documentation][doc-img]][doc-url]

[latest-release-img]: https://img.shields.io/github/release/go-pogo/env.svg?label=latest

[latest-release-url]: https://github.com/go-pogo/env/releases

[build-status-img]: https://github.com/go-pogo/env/actions/workflows/test.yml/badge.svg

[build-status-url]: https://github.com/go-pogo/env/actions/workflows/test.yml

[report-img]: https://goreportcard.com/badge/github.com/go-pogo/env

[report-url]: https://goreportcard.com/report/github.com/go-pogo/env

[doc-img]: https://godoc.org/github.com/go-pogo/env?status.svg

[doc-url]: https://pkg.go.dev/github.com/go-pogo/env


Package `env` reads and parses environment variables from various sources. 
It supports unmarshaling into any type and (over)loading the variables into the 
system's environment.

Included features are:
* Reading environment variables from various sources;
* Decoding environment variables into any type;
* Encoding environment variables from any type;
* Loading and overloading into the system's environment variables.

<hr>

```sh
go get github.com/go-pogo/env
```

```sh
import "github.com/go-pogo/env"
```

## Usage

Below example demonstrates how to decode system environment variables into a struct.

```go
package main

import (
    "github.com/davecgh/go-spew/spew"
    "github.com/go-pogo/env"
    "time"
)

func main() {
    type Config struct {
        Foo     string
        Timeout time.Duration `default:"10s"`
    }

    var conf Config
    if err := env.NewDecoder(env.System()).Decode(&conf); err != nil {
        panic(err)
    }

    spew.Dump(conf)
    // Output:
    // (env.Config) {
    //  Foo: (string) "",
    //  Timeout: (time.Duration) 10s
    // }
}
```

## Usage with `dotenv`

This example reads .env files from the _example_ directory and decodes the found variables into a struct.
```go
package main

import (
    "github.com/davecgh/go-spew/spew"
    "github.com/go-pogo/env"
    "github.com/go-pogo/env/dotenv"
    "time"
)

func main() {
    type Config struct {
        Foo     string
        Timeout time.Duration `default:"10s"`
    }

    var conf Config
    if err := env.NewDecoder(dotenv.Read("example", dotenv.None)).Decode(&conf); err != nil {
        panic(err)
    }

    spew.Dump(conf)
    // Output:
    // (dotenv.Config) {
    //  Foo: (string) (len=3) "bar",
    //  Timeout: (time.Duration) 2s
    // }
}
```

## Documentation

Additional detailed documentation is available at [pkg.go.dev][doc-url]

## Created with

<a href="https://www.jetbrains.com/?from=go-pogo" target="_blank"><img src="https://resources.jetbrains.com/storage/products/company/brand/logos/GoLand_icon.png" width="35" /></a>

## License

Copyright © 2020-2025 [Roel Schut](https://roelschut.nl). All rights reserved.

This project is governed by a BSD-style license that can be found in the [LICENSE](LICENSE) file.
