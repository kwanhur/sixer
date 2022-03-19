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

import "testing"

func TestGitHub_ValidLinks(t *testing.T) {
	type fields struct {
		Linker Linker
		git    *Git
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "valid links",
			fields: fields{
				Linker: Linker{
					timeout: 3,
				},
				git: &Git{
					Commit:  "2c563dc15c54a8deb3ba08707594d4d15da76b1b",
					Repo:    pkgApisixDashboard,
					Release: "2.11.0",
				}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GitHub{
				Linker: tt.fields.Linker,
				git:    tt.fields.git,
			}
			if err := g.ValidLinks(); (err != nil) != tt.wantErr {
				t.Errorf("ValidLinks() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
