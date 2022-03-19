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
