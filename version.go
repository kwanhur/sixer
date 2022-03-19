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

import (
	"fmt"
	"runtime"

	"github.com/spf13/cobra"
)

const (
	version = "v0.0.1"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show sixer version number",
	Run: func(cmd *cobra.Command, args []string) {
		showVersion()
	},
}

func showVersion() {
	fmt.Printf("sixer version:%s\n", version)
}

var verboseCmd = &cobra.Command{
	Use:   "verbose",
	Short: "Show sixer verbose information",
	Run: func(cmd *cobra.Command, args []string) {
		showVerbose()
	},
}

func showVerbose() {
	fmt.Printf("sixer version: %s\n", version)
	fmt.Printf("go version: %s\n", runtime.Version())
	fmt.Printf("git commit: %s\n", commit)
}
