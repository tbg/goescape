package examples

import "github.com/tbg/goescape"

// For convenient testing, all variables starting with `h` escape to the heap.
// Those starting with `s` stay on the stack.

func f1() []byte {
	var h1 struct {
		goescape.Stack
		buf [64]byte
	}
	return h1.buf[:]
}

func f2() []byte {
	h2 := struct {
		goescape.Stack
		buf [64]byte
	}{}
	return h2.buf[:]
}

func f3() bool {
	s := struct {
		goescape.Stack
		buf [64]byte
	}{}
	s.buf[0] = 1
	return s.buf[0]+s.buf[1] == 1
}
