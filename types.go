package goj

import (
	// "log"
	"errors"
	"fmt"
	"path/filepath"
	"strings"
)

var EOF = "EOF"

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
	if b < 0 {
		b = b + t
	}
	if e < 0 {
		e = e + t
	}
	if b > lhs || lhs > e {
		return false
	}
	if s > 1 {
		if (lhs-b)%s != 0 {
			return false
		}
	}
	return true
}

// common case of [*]
var pathStar = pathSlice{0, -1, 1}

func newPathIndex(s string) (pathSel, error) {
	if len(s) == 0 {
		return nil, errors.New("array index cannot be empty")
	}

	if s == "*" || s == ":" {
		return pathStar, nil
	}

	sel, err := newPathSlice(s)

	if err != nil {
		return nil, err
	}

	if sel != (pathSlice{}) {
		return sel, nil
	}

	return nil, nil
	// return newPathRange(s)
}

func newPathSlice(str string) (sel pathSlice, err error) {
	var n, b, e, s int

	r := strings.NewReader(str)

	n, _ = fmt.Fscanf(r, "%d", &b)
	if r.Len() == 0 {
		return
	}
	r.Seek(0, 0)

	n, er := fmt.Fscanf(r, "::%d", &s)
	if er != nil && er.Error() == EOF {
		n++
	}
	if n < 1 {
		r.Seek(0, 0)
		n, _ = fmt.Fscanf(r, ":%d:%d", &e, &s)
	}
	if n < 1 {
		r.Seek(0, 0)
		n, _ = fmt.Fscanf(r, "%d:%d:%d", &b, &e, &s)
	}
	if n > 0 {
		if e == 0 {
			e = -1
		}
		if s == 0 {
			s = 1
		}
		sel = pathSlice{b, e, s}
	}

	return
}

func newPathRange(s string) (sel pathIndex, err error) {
	var i int

	r := strings.NewReader(s)
	f := "%d,"

	for r.Len() > 0 {
		n, _ := fmt.Fscanf(r, f, &i)
		if n > 0 {
			sel.val = append(sel.val, i)
		}
	}

	return
}
