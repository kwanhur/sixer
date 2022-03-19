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
