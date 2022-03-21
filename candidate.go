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

import "fmt"

const (
	baseLink   = "https://dist.apache.org/repos/dist/dev/apisix/"
	pkgPrefix2 = "apache"
)

// A Candidate represents package with specified version
type Candidate struct {
	pkg string // package name, like: apisix-dashboard
	rc  string // release candidate version, like: 0.2.0
	sub bool   // sub-project
}

// PackageLink complete URL for package directory
func (c *Candidate) PackageLink() string {
	return fmt.Sprintf("%s%s", baseLink, c.Package())
}

// Package a package name with prefix "apisix"
func (c *Candidate) Package() string {
	if c.sub {
		return fmt.Sprintf("%s-%s", c.pkg, c.rc)
	}

	return c.rc
}

// Package2 a package name with prefix "apache-"
func (c *Candidate) Package2() string {
	return fmt.Sprintf("%s-%s-%s", pkgPrefix2, c.pkg, c.rc)
}

func (c *Candidate) srcTgz() string {
	return fmt.Sprintf("%s-src.tgz", c.Package2())
}

// SrcLink source package URL
func (c *Candidate) SrcLink() string {
	return fmt.Sprintf("%s/%s", c.PackageLink(), c.srcTgz())
}

func (c *Candidate) srcTgzAsc() string {
	return fmt.Sprintf("%s-src.tgz.asc", c.Package2())
}

// SrcAscLink source package asc URL
func (c *Candidate) SrcAscLink() string {
	return fmt.Sprintf("%s/%s", c.PackageLink(), c.srcTgzAsc())
}

func (c *Candidate) srcTgzSha512() string {
	return fmt.Sprintf("%s-src.tgz.sha512", c.Package2())
}

// SrcSha512Link source package sha512 URL
func (c *Candidate) SrcSha512Link() string {
	return fmt.Sprintf("%s/%s", c.PackageLink(), c.srcTgzSha512())
}
