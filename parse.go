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
	"go/ast"
	"go/parser"
	"go/token"

	"golang.org/x/tools/go/ast/astutil"
)

// A NoEscapeMap maps file names into maps of protected lines. A protected line
// is a line in the associated file on which no heap allocated variables should
// be tolerated.
type NoEscapeMap map[string]map[int]struct{}

// Parse parses the given pkgs into a NoEscapeMap. Only single packages are
// admissible, i.e. "./..." isn't, due to technical limitations of when the
// Go escape analysis is emitted; use List() to work around this.
func Parse(pkgs ...string) (NoEscapeMap, error) {
	m := NoEscapeMap{}
	for _, pkgPath := range pkgs {
		if err := parseInto(m, pkgPath); err != nil {
			return nil, err
		}
	}
	return m, nil
}

func parseInto(noescapeLines NoEscapeMap, pkgPath string) error {
	fset := token.NewFileSet()

	pkgs, err := parser.ParseDir(fset, pkgPath, nil, 0)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for filename, f := range pkg.Files {
			for _, decl := range f.Decls {
				var pre astutil.ApplyFunc = func(c *astutil.Cursor) bool {
					n := c.Node()
					if n == nil {
						return true
					}
					first, last := n.Pos()-1, n.End()-1
					_, _ = first, last
					switch t := n.(type) {
					case *ast.StructType:
						fields := t.Fields
						for _, l := range fields.List {
							if sel, ok := l.Type.(*ast.SelectorExpr); ok &&
								sel.Sel.Name == "Stack" &&
								fmt.Sprint(sel.X) == "goescape" {

								line := fset.File(n.Pos()).Line(n.Pos())
								if len(noescapeLines[filename]) == 0 {
									noescapeLines[filename] = map[int]struct{}{}
								}
								noescapeLines[filename][line] = struct{}{}
							}
						}
					}
					return true // descend
				}
				astutil.Apply(decl, pre, nil)
			}
		}
	}
	return nil
}
