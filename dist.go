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
	"bytes"
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/ProtonMail/go-crypto/openpgp"
	"github.com/ProtonMail/go-crypto/openpgp/armor"
	"github.com/ProtonMail/go-crypto/openpgp/packet"

	"github.com/parnurzeal/gorequest"
	"github.com/spf13/cobra"
)

const (
	baseLink    = "https://dist.apache.org/repos/dist/dev/apisix/"
	pkgPrefix   = "apisix"
	pkgPrefix2  = "apache-apisix"
	keyFilename = ".key"
)

// A Candidate represents package with specified version
type Candidate struct {
	pkg string // package name, like: dashboard
	rc  string // release candidate version, like: 0.2.0
}

// PackageLink complete URL for package directory
func (c *Candidate) PackageLink() string {
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

// A Dist repo include package and its asc sha512
type Dist struct {
	Candidate
	timeout   int
	force     bool // force recover existed packages
	announcer string
}

// SetTimeout set dist request timeout, unit second
func (d *Dist) SetTimeout(t int) {
	d.timeout = t
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

// ValidAllLinks validate URL links, include package and its src asc sha512
func (d *Dist) ValidAllLinks() {
	links := []string{d.PackageLink(), d.SrcLink(), d.SrcAscLink(), d.SrcSha512Link()}
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

func (d *Dist) fetchSrc() error {
	if d.force {
		if err := os.Remove(d.srcTgz()); err != nil && !os.IsNotExist(err) {
			return err
		}
	} else {
		if f, err := os.Stat(d.srcTgz()); err != nil && !os.IsNotExist(err) {
			return err
		} else if f != nil {
			return nil
		}
	}

	var err error
	r := gorequest.New()
	sa := r.Timeout(time.Duration(d.timeout) * time.Second)

	sa.Get(d.SrcLink()).EndBytes(func(res gorequest.Response, body []byte, errs []error) {
		if len(errs) != 0 {
			err = errs[0]
			return
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("non-expected response status %s", res.Status)
			return
		}

		if len(body) == 0 {
			err = fmt.Errorf("response body size zero")
			return
		}

		err = os.WriteFile(d.srcTgz(), body, 0644)
	})

	return err
}

func (d *Dist) fetchSrcChecksum() ([]byte, error) {
	var checksum []byte
	var err error

	if !d.force {
		checksum, err = os.ReadFile(d.srcTgzSha512())

		if os.IsNotExist(err) {
			goto download
		}
		return checksum, nil
	}

	if d.force {
		if err := os.Remove(d.srcTgzSha512()); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

download:
	r := gorequest.New()
	sa := r.Timeout(time.Duration(d.timeout) * time.Second)

	sa.Get(d.SrcSha512Link()).EndBytes(func(res gorequest.Response, body []byte, errs []error) {
		if len(errs) != 0 {
			err = errs[0]
			return
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("non-expected response status %s", res.Status)
			return
		}

		if len(body) == 0 {
			err = fmt.Errorf("response checksum size zero")
			return
		}

		checksum = body
	})

	return checksum, nil
}

func (d *Dist) validKey() (bool, error) {
	key, err := os.Open(keyFilename)
	if err != nil {
		return false, err
	}
	defer key.Close()

	entities, err := openpgp.ReadArmoredKeyRing(key)
	if err != nil {
		return false, err
	}

	if len(entities) != 0 {
		return false, fmt.Errorf("should be one entity in key")
	}

	id := entities[0].PrimaryIdentity()
	if id == nil {
		return false, fmt.Errorf("there's no primary identity")
	}

	return strings.HasPrefix(id.Name, d.announcer), nil
}

func (d *Dist) fetchKey() (*os.File, error) {
	if d.force {
		if err := os.Remove(keyFilename); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		goto export
	} else {
		if _, err := os.Stat(keyFilename); err != nil && !os.IsNotExist(err) {
			return nil, err
		} else if os.IsNotExist(err) {
			goto export
		} else {
			if ok, err := d.validKey(); err != nil {
				return nil, err
			} else if !ok {
				goto export
			}
			return os.Open(keyFilename)
		}
	}

export:
	cmd := exec.Command("gpg", "--export", d.announcer, "--output", keyFilename)
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	if !cmd.ProcessState.Success() {
		return nil, fmt.Errorf("gpg export %s keyfile failed", d.announcer)
	}

	return os.Open(keyFilename)
}

func (d *Dist) fetchSrcSignature() (*os.File, error) {
	var err error

	if !d.force {
		_, err = os.Stat(d.srcTgzAsc())

		if os.IsNotExist(err) {
			goto download
		}
		return os.Open(d.srcTgzAsc())
	}

	if d.force {
		if err := os.Remove(d.srcTgzAsc()); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
	}

download:
	r := gorequest.New()
	sa := r.Timeout(time.Duration(d.timeout) * time.Second)

	sa.Get(d.SrcAscLink()).EndBytes(func(res gorequest.Response, body []byte, errs []error) {
		if len(errs) != 0 {
			err = errs[0]
			return
		}

		if res.StatusCode != http.StatusOK {
			err = fmt.Errorf("non-expected response status %s", res.Status)
			return
		}

		if len(body) == 0 {
			err = fmt.Errorf("response signature size zero")
			return
		}

		err = os.WriteFile(d.srcTgzAsc(), body, 0644)
	})

	if err != nil {
		return nil, err
	}

	return os.Open(d.srcTgzAsc())
}

func (d *Dist) checksum(src []byte, body []byte) (bool, error) {
	var err error

	sums := bytes.Split(body, []byte("  "))
	if len(sums) != 2 {
		return false, fmt.Errorf("invalid checksum body")
	}

	if len(src) == 0 {
		filename := strings.TrimSpace(string(sums[1]))
		src, err = os.ReadFile(filename)
		if err != nil {
			return false, err
		}
	}

	hash := sha512.New()
	hash.Write(src)
	checksum := fmt.Sprintf("%x", hash.Sum(nil))
	ret := strings.Compare(string(sums[0]), checksum)

	return ret == 0, nil
}

func (d *Dist) signature(src []byte, sign *os.File, key *os.File) error {
	block, err := armor.Decode(sign)
	if err != nil {
		return err
	}

	if block.Type != "PGP SIGNATURE" {
		return fmt.Errorf("not an armor signature")
	}

	pkg, err := packet.Read(block.Body)
	if err != nil {
		return err
	}

	signature, ok := pkg.(*packet.Signature)
	if !ok {
		return fmt.Errorf("not a valid signature file")
	}

	kblock, err := armor.Decode(key)
	if err != nil {
		return err
	}

	if kblock.Type != "PGP PUBLIC KEY BLOCK" {
		return fmt.Errorf("not an armored public key")
	}

	kpkg, err := packet.Read(kblock.Body)
	if err != nil {
		return err
	}

	pubKey, ok := kpkg.(*packet.PublicKey)
	if !ok {
		return fmt.Errorf("not a valid public key file")
	}

	hash := signature.Hash.New()
	_, err = hash.Write(src)
	if err != nil {
		return err
	}

	return pubKey.VerifySignature(hash, signature)
}

// ValidChecksum validate from sha512 checksum file
func (d *Dist) ValidChecksum() (bool, error) {
	if err := d.fetchSrc(); err != nil {
		return false, err
	}

	checksum, err := d.fetchSrcChecksum()
	if err != nil {
		return false, err
	}

	src, err := os.ReadFile(d.srcTgz())
	if err != nil {
		return false, err
	}

	return d.checksum(src, checksum)
}

// ValidSignature validate from asc file
func (d *Dist) ValidSignature() (bool, error) {
	src, err := os.ReadFile(d.srcTgz())
	if err != nil {
		return false, err
	}

	sign, err := d.fetchSrcSignature()
	if err != nil {
		return false, err
	}
	defer sign.Close()

	key, err := d.fetchKey()
	if err != nil {
		return false, err
	}
	defer key.Close()

	err = d.signature(src, sign, key)
	return err == nil, err
}

// NewDashboardDist dashboard dist
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
	Use:              "dashboard",
	Short:            "apisix dashboard package verifier",
	PersistentPreRun: sixerPreRun,
	Run:              dashboardRun,
}

// verify package
// 1. check links
// 2. download packages
// 3. verify checksum and signature
// 4. untar then check LICENSE and NOTICE
func dashboardRun(cmd *cobra.Command, args []string) {
	dist := NewDashboardDist()
	dist.ValidAllLinks()

	if ok, err := dist.ValidChecksum(); err != nil {
		log.Fatalf("dist validate checksum failed: %s\n", err)
	} else if ok {
		log.Println("dist validate checksum successfully")
	} else {
		log.Fatalln("dist validate checksum failed")
	}

	if ok, err := dist.ValidSignature(); err != nil {
		log.Fatalf("dist validate signature failed:%s", err)
	} else if ok {
		log.Println("dist validate signature successfully")
	} else {
		log.Fatalln("dist validate signature failed")
	}
}
