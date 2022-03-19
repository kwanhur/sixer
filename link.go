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
	"net/http"
	"time"

	"github.com/parnurzeal/gorequest"
)

// Linker validate link
type Linker struct {
	timeout uint
}

// SetTimeout set timeout when request
func (l *Linker) SetTimeout(timeout uint) {
	l.timeout = timeout
}

// Head use http.HEAD to valid
func (l *Linker) Head(link string) (bool, error) {
	var valid bool
	var err error

	r := gorequest.New()
	sa := r.Timeout(time.Duration(l.timeout) * time.Second)

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
