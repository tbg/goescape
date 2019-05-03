// Copyright 2019 @tbg
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package goescape_test

import (
	"bytes"
	"flag"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tbg/goescape"
)

var rewrite = flag.Bool("rewrite", false, "regenerate the text fixtures from the test output")

// TestFixtures lints ./examples/... and verifies that for each .go file, the
// linter messages match exactly the corresponding .go.out file.
func TestFixtures(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	fs, err := goescape.Lint("./examples/...", cwd)
	if err != nil {
		t.Fatal(err)
	}

	byFile := map[string][]string{}

	for _, f := range fs {
		byFile[f.File] = append(byFile[f.File], f.Error())
	}

	fixtures, err := filepath.Glob("./examples/*.go.out")
	if err != nil {
		t.Fatal(err)
	}

	byFixture := map[string][]string{}
	for _, fixture := range fixtures {
		b, err := ioutil.ReadFile(fixture)
		if err != nil {
			t.Fatal(err)
		}
		filename := fixture[:len(fixture)-4] // remove ".out"
		byFixture[filename] = strings.Split(string(bytes.TrimSpace(b)), "\n")
	}

	if *rewrite {
		for filename := range byFixture {
			if err := os.Remove(filename + ".out"); err != nil {
				t.Fatal(err)
			}
		}
		for filename, lines := range byFile {
			if err := ioutil.WriteFile(filename+".out", []byte(strings.Join(lines, "\n")), 0644); err != nil {
				t.Fatal(err)
			}
		}
	}

	assert.Equal(t, byFixture, byFile,
		"mismatch detected, use -rewrite if this is expected",
	)
}

func TestList(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	pkgs, err := goescape.List("./...", cwd)
	if err != nil {
		t.Fatal(err)
	}

	exp := []string{".", "./examples"}
	assert.Equal(t, exp, pkgs)
}
