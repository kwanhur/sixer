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
	// Version marked to show version
	Version bool
	// Verbose marked to show verbose
	Verbose bool

	candidate string
	commitID  string
	announcer string
	timeout   uint
)

// BindVerFlags add Version Verbose flags
func BindVerFlags(flags *pflag.FlagSet) {
	flags.BoolVarP(&Version, "version", "v", false, "Show sixer version number")
	flags.BoolVarP(&Verbose, "verbose", "V", false, "Show sixer verbose information")
}

// BindGlobalFlags bind persistent flags
func BindGlobalFlags(flags *pflag.FlagSet) {
	flags.UintVarP(&timeout, "timeout", "t", 0, "Specify request link timeout, unit: second")
	flags.StringVarP(&candidate, "candidate", "c", "", "Specify release candidate version,like 0.2.0")
	flags.StringVarP(&announcer, "announcer", "a", "", "Specify release candidate announcer")
	flags.StringVarP(&commitID, "commit", "C", "", "Specify release commit id")
}
