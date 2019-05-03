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

package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/tbg/goescape"
)

func main() {
	if err := mainE(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func mainE() error {
	var pkgSpec string
	var cwd string

	switch len(os.Args) {
	case 2:
		pkgSpec = os.Args[1]
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return err
		}
	case 3:
		pkgSpec = os.Args[1]
		cwd = os.Args[2]
	default:
		return fmt.Errorf("invalid args: %v", os.Args)
	}

	var err error
	cwd, err = filepath.Abs(cwd)
	if err != nil {
		return errors.Wrap(err, "computing absolute path")
	}
	if err := os.Chdir(cwd); err != nil {
		return errors.Wrap(err, "Chdir")
	}

	fs, err := goescape.Lint(pkgSpec, cwd)
	if err != nil {
		return errors.Wrap(err, "Lint")
	}
	var buf strings.Builder
	for _, f := range fs {
		buf.WriteString(f.Error())
		buf.WriteRune('\n')
	}
	if len(fs) > 0 {
		return errors.New(buf.String())
	}
	return nil
}
