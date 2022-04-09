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
	"strings"
)

const (
	githubApacheOgz = "https://github.com/apache"
)

// Git represent github repo's info
type Git struct {
	Repo    string
	Commit  string
	Release string
	Tag     string
}

// MarkdownID fetch markdown from release, 2.11.0 -> 2110
func (g *Git) MarkdownID() string {
	return strings.ReplaceAll(g.Release, ".", "")
}

// GitHub validator for github link
type GitHub struct {
	Linker

	git *Git
}

// NewGitHub GitHub instance
func NewGitHub(g *Git) (*GitHub, error) {
	if g == nil {
		return nil, fmt.Errorf("git was nill")
	}

	return &GitHub{
		Linker: Linker{
			timeout: timeout,
		},
		git: g,
	}, nil
}

func (g *GitHub) releaseNoteLink() string {
	return fmt.Sprintf("%s/%s/blob/release/%s/CHANGELOG.md#%s", githubApacheOgz, g.git.Repo, g.git.Tag, g.git.MarkdownID())
}

func (g *GitHub) releaseCommitLink() string {
	return fmt.Sprintf("%s/%s/commit/%s", githubApacheOgz, g.git.Repo, g.git.Commit)
}

// ValidLinks validate release links
func (g *GitHub) ValidLinks() error {
	links := []string{g.releaseNoteLink(), g.releaseCommitLink()}
	for _, link := range links {
		if ok, err := g.Linker.Head(link); err != nil {
			log.Printf("github %s validate bad ❌ %s\n", link, err)
			return err
		} else if ok {
			log.Printf("github %s validate ok ✅\n", link)
		} else {
			log.Printf("github %s validate bad ❌\n", link)
		}
	}

	return nil
}
