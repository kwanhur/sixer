// Copyright 2022 kwanhur
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
package main

import "github.com/spf13/pflag"

var (
	Version bool
	Verbose bool

	ReleaseCandidate string
)

// BindVerFlags add Version Verbose flags
func BindVerFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&Version, "version", "v", false, "Show apisixer's version number")
	flags.BoolVarP(&Verbose, "verbose", "V", false, "Show apisixer's verbose information")
}

func BindRCFlag(flag *pflag.FlagSet) {
	flag.StringVarP(&ReleaseCandidate, "release-candidate", "r", "", "Specify apisix release candidate version,like 0.2.0")
}
