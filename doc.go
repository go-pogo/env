// Copyright (c) 2020, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package env implements encoding and decoding of environment variables from
files, the OS or other data sources. It supports unmarshalling to structs and
maps. The mapping between environment variables and Go values is described in
the documentation for the Marshal and Unmarshal functions.

# Supported types

This package uses the rawconv package to parse/unmarshal any string value to
its Go type equivalent. Additionally, custom types may implement the Unmarshaler
interface to implement its own unmarshalling rules, or register an unmarshaler
function using rawconv.Register.

# Load and overload

Additional os.Environ entries can be loaded using the ReadAndLoad, OpenAndLoad,
ReadAndOverload and OpenAndOverload functions. The source is read any bash style
variables are replaced before being set to the system using Setenv.

# Dotenv

The dotenv package supports reading and loading environment variables from .env
files based on active environment (e.g. prod, dev etc.).

# Writing

This package can also write environment variables to an io.Writer.
*/
package env
