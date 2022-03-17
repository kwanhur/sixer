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
	"archive/tar"
	"bytes"
	"compress/gzip"
	"crypto/sha512"
	"fmt"
	"io"
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
	keyFilename = ".key"
)

// A Dist repo include package and its asc sha512
type Dist struct {
	Candidate
	timeout   uint
	force     bool // force recover existed packages
	announcer string
}

func (d *Dist) validAttrs() (bool, error) {
	if d.announcer == "" {
		return false, fmt.Errorf("announcer not specified")
	}

	return true, nil
}

// SetTimeout set dist request timeout, unit second
func (d *Dist) SetTimeout(t uint) {
	d.timeout = t
}

func (d *Dist) validLink(link string) (bool, error) {
	var valid bool
	var err error

	r := gorequest.New()
	sa := r.Timeout(time.Duration(d.timeout) * time.Second)

	sa.Head(link).End(func(res gorequest.Response, body string, errs []error) {
		for _, e := range errs {
			if e != nil {
				err = e
			}
		}
		if err != nil {
			return
		}

		if res.StatusCode == http.StatusOK {
			valid = true
		}
	})

	return valid, err
}

// ValidAllLinks validate URL links, include package and its src asc sha512
func (d *Dist) ValidAllLinks() error {
	links := []string{d.PackageLink(), d.SrcLink(), d.SrcAscLink(), d.SrcSha512Link()}
	for _, link := range links {
		if ok, err := d.validLink(link); err != nil {
			log.Printf("dist %s validate bad❌:%s\n", link, err)
			return err
		} else if ok {
			log.Printf("dist %s validate ok✅\n", link)
		} else {
			log.Printf("dist %s validate bad❌\n", link)
		}
	}

	return nil
}

func (d *Dist) fetchSrcTgz() error {
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
		for _, e := range errs {
			if e != nil {
				err = e
			}
		}
		if err != nil {
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
			err = nil
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
		for _, e := range errs {
			if e != nil {
				err = e
			}
		}
		if err != nil {
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

		err = os.WriteFile(d.srcTgzSha512(), body, 0644)
		if err != nil {
			return
		}

		checksum = body
	})

	if err != nil {
		return nil, err
	}

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

	if len(entities) != 1 {
		return false, fmt.Errorf("should be one entity in key")
	}

	id := entities[0].PrimaryIdentity()
	if id == nil {
		return false, fmt.Errorf("there's no primary identity")
	}

	return strings.HasPrefix(id.Name, d.announcer), nil
}

func (d *Dist) fetchKey() (*os.File, error) {
	if _, err := d.validAttrs(); err != nil {
		return nil, err
	}

	if d.force {
		if err := os.Remove(keyFilename); err != nil && !os.IsNotExist(err) {
			return nil, err
		}
		goto export
	} else {
		if _, err := os.Stat(keyFilename); err != nil {
			if !os.IsNotExist(err) {
				return nil, err
			}

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
	cmd := exec.Command("gpg", "--armor", "--output", keyFilename, "--export", d.announcer)
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
		for _, e := range errs {
			if e != nil {
				err = e
			}
		}
		if err != nil {
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

func (d *Dist) checkExtras() (bool, error) {
	f, err := os.Open(d.srcTgz())
	if err != nil {
		return false, err
	}
	defer f.Close()

	gzf, err := gzip.NewReader(f)
	if err != nil {
		return false, err
	}
	defer gzf.Close()

	tf := tar.NewReader(gzf)
	for {
		hdr, err := tf.Next()
		if err == io.EOF {
			break
		}

		switch hdr.Typeflag {
		case tar.TypeReg:
			switch hdr.Name {
			case "./LICENSE":
				log.Println("LICENSE ok ✅")
			case "./NOTICE":
				log.Println("NOTICE ok ✅")
			}
		}
	}

	return true, nil
}

// Clean cleans download files
func (d *Dist) Clean() error {
	if err := os.Remove(d.srcTgzSha512()); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Remove(d.srcTgzAsc()); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Remove(d.srcTgz()); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := os.Remove(keyFilename); err != nil && !os.IsNotExist(err) {
		return err
	}

	return nil
}

// Verify package
// 1. check links
// 2. download packages
// 3. verify checksum and signature
// 4. untar then check LICENSE and NOTICE
func (d *Dist) Verify() {
	if err := d.fetchSrcTgz(); err != nil {
		log.Printf("dist fetch src tgz bad❌:%s\n", err)
	}

	if ok, err := d.ValidChecksum(); err != nil {
		log.Printf("dist validate checksum bad❌: %s\n", err)
	} else if ok {
		log.Println("dist validate checksum ok✅")
	} else {
		log.Println("dist validate checksum bad❌")
	}

	if ok, err := d.ValidSignature(); err != nil {
		log.Printf("dist validate signature bad❌:%s", err)
	} else if ok {
		log.Println("dist validate signature ok✅")
	} else {
		log.Println("dist validate signature bad❌")
	}

	if _, err := d.checkExtras(); err != nil {
		log.Printf("dist check extras bad❌:%s\n", err)
	}
}

// NewDashboardDist dashboard dist
func NewDashboardDist() Dist {
	return Dist{
		Candidate: Candidate{
			pkg: "dashboard",
			rc:  candidate,
		},
		force:     force,
		announcer: announcer,
		timeout:   timeout,
	}
}

var dashboardCmd = &cobra.Command{
	Use:              "dashboard",
	Short:            "apisix dashboard package verifier",
	PersistentPreRun: sixerPreRun,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		dist := NewDashboardDist()
		return dist.ValidAllLinks()
	},
	Run: func(cmd *cobra.Command, args []string) {
		dist := NewDashboardDist()
		dist.Verify()
	},
	PostRunE: func(cmd *cobra.Command, args []string) error {
		dist := NewDashboardDist()
		return dist.Clean()
	},
}
