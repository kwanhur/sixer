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
	"log"
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
)

const (
	baseLink   = "https://dist.apache.org/repos/dist/dev/apisix/"
	pkgPrefix  = "apisix"
	pkgPrefix2 = "apache-apisix"
)

// A Candidate represents package with specified version
type Candidate struct {
	pkg string // package name, like: dashboard
	rc  string // release candidate version, like: 0.2.0
}

func (c *Candidate) PkgLink() string {
	return fmt.Sprintf("%s%s", baseLink, c.Package())
}

// Package a package name with prefix "apisix"
func (c *Candidate) Package() string {
	return fmt.Sprintf("%s-%s-%s", pkgPrefix, c.pkg, c.rc)
}

// Package2 a package name with prefix "apache-apisix"
func (c *Candidate) Package2() string {
	return fmt.Sprintf("%s-%s-%s", pkgPrefix2, c.pkg, c.rc)
}

func (c *Candidate) SrcLink() string {
	return fmt.Sprintf("%s/%s-src.tgz", c.PkgLink(), c.Package2())
}

func (c *Candidate) SrcAscLink() string {
	return fmt.Sprintf("%s/%s-src.tgz.asc", c.PkgLink(), c.Package2())
}

func (c *Candidate) SrcSha512Link() string {
	return fmt.Sprintf("%s/%s-src.tgz.sha512", c.PkgLink(), c.Package2())
}

type Dist struct {
	Candidate
	timeout int
}

func (d *Dist) keysLink() string {
	return fmt.Sprintf("%s%s", baseLink, "KEYS")
}

func (d *Dist) validLink(link string) (bool, error) {
	var valid bool
	var err error

	r := gorequest.New()
	sa := r.Timeout(time.Duration(d.timeout) * time.Second)

	sa.Head(link).End(func(res gorequest.Response, body string, errs []error) {
		if len(errs) != 0 {
			err = errs[0]
			return
		}

		if res.StatusCode == http.StatusOK {
			valid = true
		}
	})

	return valid, err
}

func (d Dist) ValidAllLinks() {
	links := []string{d.PkgLink(), d.SrcLink(), d.SrcAscLink(), d.SrcSha512Link()}
	for _, link := range links {
		if ok, err := d.validLink(link); err != nil {
			log.Fatalf("dist %s validate failed:%s\n", link, err)
		} else if ok {
			log.Printf("dist %s validate successfully\n", link)
		} else {
			log.Printf("dist %s validate failed\n", link)
		}
	}
}

func NewDashboardDist() Dist {
	return Dist{
		Candidate: Candidate{
			pkg: "dashboard",
			rc:  ReleaseCandidate,
		},
		timeout: 3,
	}
}

var dashboardCmd = &cobra.Command{
	Use:   "dashboard",
	Short: "apisix dashboard package verifier",
	Run: func(cmd *cobra.Command, args []string) {
		verifyDashboard()
	},
}

// verify package
// 1. check links
// 2. download packages
// 3. verify checksum and signature
// 4. untar then check LICENSE and NOTICE
func verifyDashboard() {
	dist := NewDashboardDist()
	dist.ValidAllLinks()
}