go-env
======

[![Latest release][latest-release-img]][latest-release-url]
[![Travis build status][travis-build-img]][travis-build-url]
[![Go Report Card][go-report-img]][go-report-url]
[![GoDoc documentation][go-doc-img]][go-doc-url]
![Minimal Go version][go-version-img]

[latest-release-img]: https://img.shields.io/github/release/roeldev/go-env.svg?label=latest
[latest-release-url]: https://github.com/roeldev/go-env/releases
[travis-build-img]: https://img.shields.io/travis/roeldev/go-env.svg
[travis-build-url]: https://travis-ci.org/roeldev/go-env
[go-report-img]: https://goreportcard.com/badge/github.com/roeldev/go-env
[go-report-url]: https://goreportcard.com/report/github.com/roeldev/go-env
[go-doc-img]: https://godoc.org/github.com/roeldev/go-env?status.svg
[go-doc-url]: https://pkg.go.dev/github.com/roeldev/go-env
[go-version-img]: https://img.shields.io/github/go-mod/go-version/roeldev/go-env

Go package for parsing environment variables. It can be used to read from any `io.Reader`, `os.Environ()` or (CLI) arguments.


## Install
```sh
go get github.com/roeldev/go-env
```


## Import
```go
import "github.com/roeldev/go-env"
```

## Usage

### Reading a .env file
``go
func main() {
}
``

### CLI arguments

## Supported formats

- `key=value`; value without quotes
- `key="value"`; value with double quotes
- `key='value'`; value with single quotes


## Documentation
Additional detailed documentation is available at [go.dev][go-doc-url]


### Created with
<a href="https://www.jetbrains.com/?from=roeldev/go-env" target="_blank"><img src="https://pbs.twimg.com/profile_images/1206615658638856192/eiS7UWLo_400x400.jpg" width="35" /></a>


## License
[GPL-3.0+](LICENSE) © 2020 [Roel Schut](https://roelschut.nl)
