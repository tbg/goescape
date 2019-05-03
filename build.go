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
	"fmt"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/ghemawat/stream"
	"github.com/pkg/errors"
)

// An EscapeMap maps filenames into a map from line numbers to variables that
// were found to escape to the heap.
type EscapeMap map[string]map[int][]string

var reMovedToHeap = regexp.MustCompile(`^(.*):(\d+):(\d+): moved to heap: (.*)$`)

// Build builds the given packages with the goal of extracting relevant information
// from the escape analysis, reflected in the returned EscapeMap.
func Build(pkgs ...string) (EscapeMap, error) {
	m := EscapeMap{}
	for _, pkg := range pkgs {
		if err := buildInto(m, pkg); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func buildInto(m EscapeMap, pkgName string) error {
	args := []string{
		"go", "build", "-gcflags=-m", pkgName,
	}

	// NB: not worth only capturing stderr here, if the command doesn't launch
	// this is easier to report a better error.
	out, err := exec.Command(args[0], args[1:]...).CombinedOutput()
	if err != nil {
		return errors.Errorf("%v\n%s\n%s", args, out, err)
	}

	if err := stream.ForEach(stream.ReadLines(bytes.NewReader(out)), func(s string) {
		sl := reMovedToHeap.FindStringSubmatch(s)
		if len(sl) == 0 {
			// fmt.Println(s)
			return
		}
		filename, sLine, _, name := sl[1], sl[2], sl[3], sl[4]
		if _, ok := m[filename]; !ok {
			m[filename] = map[int][]string{}
		}
		line, err := strconv.Atoi(sLine)
		if err != nil {
			fmt.Println(sLine)
		}
		m[filename][line] = append(m[filename][line], name)
	}); err != nil {
		return err
	}

	return nil
}
