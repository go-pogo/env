// Copyright (c) 2022, Roel Schut. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package dotenv

// ActiveEnvironment is the active environment.
type ActiveEnvironment string

func (ae ActiveEnvironment) FileName() string {
	return ".env." + ae.String()
}

func (ae ActiveEnvironment) String() string { return string(ae) }

const (
	None        ActiveEnvironment = ""
	Development ActiveEnvironment = "dev"
	Testing     ActiveEnvironment = "test"
	Production  ActiveEnvironment = "prod"
)
