package goj

import (
	// "log"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var EOF = "EOF"
var UEOF = "unexpected EOF"

// pathSel is an interface for each path component
type pathSel interface {
	Equal(v pathSel) bool
}

type pathRec struct{}

func (p pathRec) Equal(v pathSel) bool {
	return true
}

func (p pathRec) String() string {
	return "**"
}

type pathParent struct{}

func (p pathParent) Equal(v pathSel) bool {
	return true
}

func (p pathParent) String() string {
	return ".."
}

type pathKey struct {
	val string
}

func (p pathKey) Equal(v pathSel) bool {
	x, ok := v.(pathKey)
	if !ok {
		return false
	}
	if p.val == x.val {
		return true
	}
	ok, err := filepath.Match(p.val, x.val)
	return ok && err == nil
}

type pathVal struct {
	val interface{}
}

func (p pathVal) Equal(v pathSel) bool {
	x, ok := v.(pathVal)
	if !ok {
		return false
	}
	// type assertion?
	rhs := fmt.Sprint(p.val)
	lhs := fmt.Sprint(x.val)
	if rhs == lhs {
		return true
	}
	ok, err := filepath.Match(rhs, lhs)
	return ok && err == nil
}

// original index value of []interface{}
type pathIdx struct {
	val int
	len int
}

func (p pathIdx) Equal(v pathSel) bool {
	return true
}

// jsonpath index value
type pathIndex struct {
	val []int
}

func (p pathIndex) Equal(v pathSel) bool {
	x, ok := v.(pathIdx)
	if !ok {
		return false
	}

	for _, i := range p.val {
		if i == x.val {
			return true
		}
	}

	return false
}

type pathSlice struct {
	b, e, s int
}

func (p pathSlice) Equal(v pathSel) bool {
	x, ok := v.(pathIdx)
	if !ok {
		return false
	}
	lhs := x.val
	t := x.len
	b, e, s := p.b, p.e, p.s
	// reverse index
	if b < 0 {
		b = b + t
	}
	if e < 0 {
		e = e + t
	}
	if b > lhs || lhs > e {
		return false
	}
	if s > 1 && (lhs-b)%s != 0 {
		return false
	}
	return true
}

// common case of [*]
var pathStar = pathSlice{0, -1, 1}

func newPathIndex(s string) (sel pathSel, err error) {
	if len(s) == 0 {
		return nil, errors.New("array index cannot be empty")
	}

	if s == "*" || s == ":" || s == "::" {
		return pathStar, nil
	}

	r := strings.NewReader(s)

	var ok bool
	sel, ok = newPathSlice(r)

	if ok {
		return
	}

	// rewind
	r.Seek(0, 0)
	sel, ok = newPathSet(r)

	if ok {
		return
	}

	return nil, errors.New("invalid array index")
}

// scan reader for new slice
// does not handle whitespace
func newPathSlice(r *strings.Reader) (sel pathSlice, ok bool) {
	var b, e, s int
	var c byte

	fmt.Fscanf(r, "%d", &b)
	fmt.Fscanf(r, "%c", &c)
	if c != ':' {
		return sel, false
	}
	fmt.Fscanf(r, "%d", &e)
	fmt.Fscanf(r, "%c", &c)
	if c != ':' {
		return sel, false
	}
	fmt.Fscanf(r, "%d", &s)
	if r.Len() != 0 {
		return sel, false
	}

	// defaults
	if e == 0 {
		e = -1
	}
	if s == 0 {
		s = 1
	}

	return pathSlice{b, e, s}, true
}

func newPathSet(r *strings.Reader) (sel pathIndex, ok bool) {
	var n, i int
	var e error

	d := "%d"
	c := ","
	ok = true

	for r.Len() > 0 {
		n, e = fmt.Fscanf(r, d, &i)
		if n == 0 || e != nil {
			return pathIndex{}, false
		}
		_, e = fmt.Fscanf(r, c)
		if e != nil && e.Error() != UEOF {
			return pathIndex{}, false
		}
		if n > 0 {
			sel.val = append(sel.val, i)
		}
	}

	return
}
