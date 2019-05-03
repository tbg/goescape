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

package goescape

import (
	"fmt"
	"github.com/pkg/errors"
	"sort"
	"strings"
)

// A Failure is a message about a variable that was supposedly pinned to the
// stack, but was in fact found to escape to the heap.
type Failure struct {
	File string
	Line int
	Vars []string
}

// Failure implements error.
func (f *Failure) Error() string {
	return fmt.Sprintf("%s:%d: unexpectedly escaped: %s",
		f.File, f.Line, strings.Join(f.Vars, ", "),
	)
}

// Lint runs List, Parse, Build, and LintMap with the given pkgSpec. The current
// working dir helps create paths in the output that are short but can be opened
// directly.
func Lint(pkgSpec string, cwd string) ([]Failure, error) {
	pkgs, err := List(pkgSpec, cwd)
	if err != nil {
		return nil, errors.Wrap(err, "List")
	}
	mNoEscape, err := Parse(pkgs...)
	if err != nil {
		return nil, errors.Wrap(err, "Parse")
	}

	mEscape, err := Build(pkgs...)
	if err != nil {
		return nil, errors.Wrap(err, "Build")
	}

	return LintMap(mNoEscape, mEscape), nil
}

// LintMap checks from its input whether any variables escaped to the heap but
// were forbidden from doing so.
func LintMap(m NoEscapeMap, esc EscapeMap) []Failure {
	var fs []Failure
	for filename, ls := range esc {
		for line, vars := range ls {
			if _, mustNotEscape := m[filename][line]; mustNotEscape {
				fs = append(fs, Failure{File: filename, Line: line, Vars: vars})
			}
		}
	}
	// Avoid unstable result order due to map ordering above.
	sort.Slice(fs, func(i, j int) bool {
		if fs[i].File != fs[j].File {
			return fs[i].File < fs[j].File
		}
		return fs[i].Line < fs[j].Line
	})
	return fs
}
