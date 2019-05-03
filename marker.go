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

// Stack is a marker that can be embedded in an anonymous struct to indicate
// that the assignee of the declaration containing the anonymous struct should
// not escape to the heap.
//
// Allocation on the stack is often preferrable from a performance point of
// view, and once a variable has been optimized to live on the stack, it becomes
// desirable to prevent a later regression.
//
// This marker, together with a linter built around the Lint() method in this
// same package can achieve that: A variable declaration
//
//     var doNotEscape SomeType = f()
//
// can be annotated via
//
//     var doNotEscape struct{
//         goescape.Stack
//         SomeType
//     }
//     doNotEscape.SomeType = f()
//
// A similar technique applies to assignments:
//
//     doNotEscape := struct{
//         goescape.Stack
//         SomeType
//     }{
//         SomeType: f(),
//     }
//
// Some caveats apply due to the linter implementation, which operates purely
// on the AST (i.e. it has no type information):
//
// - The effect of the anonymous marked struct is to protect the line
//   on which it begins from declaring any heap allocated values. It is not
//   possible to protect a stack-allocated var when heap allocations on the
//   same line occur.
// - In particular, the struct must be defined *anonymously and inline*;
//   otherwise the declaration won't be protected at all.
// - The marker type is detected via string matching on the identifier
//   "goescape"."Stack". No dot imports, imports under a different name, or type
//   aliases are allowed.
type Stack struct{}
