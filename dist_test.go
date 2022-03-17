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
	"os"
	"testing"

	"github.com/ProtonMail/go-crypto/openpgp"
)

func TestDist_Validate(t *testing.T) {
	candidate = "2.11.0"

	tests := []struct {
		name string
		dist Dist
	}{
		{
			name: "test dashboard legal rc",
			dist: NewDashboardDist(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.dist.ValidAllLinks()
		})
	}
}

func TestDist_Validate2(t *testing.T) {
	candidate = "2.11.1"

	tests := []struct {
		name string
		dist Dist
	}{
		{
			name: "test dashboard illegal rc",
			dist: NewDashboardDist(),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Run(tt.name, func(t *testing.T) {
				tt.dist.ValidAllLinks()
			})
		})
	}
}

func TestDist_ValidChecksum(t *testing.T) {
	candidate = "2.11.0"

	dist := NewDashboardDist()
	ok, err := dist.ValidChecksum()
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.FailNow()
	}
}

func TestDist_downloadSrc(t *testing.T) {
	candidate = "2.11.0"
	force = true
	//timeout = 3

	tests := []struct {
		name    string
		dist    Dist
		wantErr bool
	}{
		{
			name:    "download dashboard src package",
			dist:    NewDashboardDist(),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			//tt.dist.SetTimeout(10)
			if err := tt.dist.fetchSrcTgz(); (err != nil) != tt.wantErr {
				t.Errorf("fetchSrcTgz() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestDist_publicKey(t *testing.T) {
	key, err := os.Open("key")
	if err != nil {
		t.Error(err)
	}
	defer key.Close()

	entities, err := openpgp.ReadArmoredKeyRing(key)
	if err != nil {
		t.Error(err)
	}

	if len(entities) != 1 {
		t.Logf("key ring length should be one")
		t.FailNow()
	}

	entity := entities[0]
	if len(entity.Identities) != 1 {
		t.Logf("identity should be one")
		t.FailNow()
	}

	iden := entity.PrimaryIdentity()
	t.Logf("%s", iden.Name)

	pubKey := entity.PrimaryKey

	t.Logf("%s", pubKey.KeyIdShortString())
	t.Logf("%s", pubKey.KeyIdString())
}

func TestDist_fetchKey(t *testing.T) {
	dist := NewDashboardDist()
	dist.announcer = "kwanhur"

	f, err := dist.fetchKey()
	if err != nil {
		t.Error(err)
	}
	f.Close()
}

func TestDist_validKey(t *testing.T) {
	dist := NewDashboardDist()
	dist.announcer = "kwanhur"

	ok, err := dist.validKey()
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Logf("announcer %s valid key fail", dist.announcer)
		t.FailNow()
	}

	dist.announcer = "Zeping Bai"
	ok, err = dist.validKey()
	if err != nil {
		t.Error(err)
	}

	if !ok {
		t.Logf("announcer %s valid expect fail ok", dist.announcer)
	}
}
