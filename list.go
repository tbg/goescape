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
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/ghemawat/stream"
	"github.com/pkg/errors"
)

// List invokes go list and massages the output so that we get relative paths
// that begin with a "." (which can be used as go pkg specs).
func List(pkgSpec string, cwd string) ([]string, error) {
	args := []string{"go", "list", "-f", "{{ .Dir }}", pkgSpec}
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return nil, errors.Errorf("%v\n%s\n%s", args, out, err)
	}

	var pkgs []string
	if err := stream.ForEach(stream.ReadLines(bytes.NewReader(out)), func(pkg string) {
		pkgs = append(pkgs, pkg)
	}); err != nil {
		return nil, err
	}

	for i := range pkgs {
		var err error
		pkgs[i], err = filepath.Rel(cwd, pkgs[i])
		if len(pkgs[i]) == 0 || pkgs[i][0] != '.' {
			pkgs[i] = "." + string(filepath.Separator) + pkgs[i]
		}
		if err != nil {
			return nil, err
		}
	}

	return pkgs, nil
}
