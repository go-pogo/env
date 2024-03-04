// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Package dotenv supports reading and loading environment variables from .env
files based on active environment (e.g. prod, dev etc.). The order or reading
is as follows:
- .env
- .env.local
- .env.{active-env}
- .env.{active-env}.local

It is recommended to not commit any .local files to the repository as these
represent variables that are specific to your local environment.
*/
package dotenv
