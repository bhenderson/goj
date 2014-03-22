package goj

import (
	// "log"
	"fmt"
	"path/filepath"
	"strings"
)

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
	val string
}

// single value index
func (p pathIndex) Equal(v pathSel) bool {
	x, ok := v.(pathIdx)
	if !ok {
		return false
	}
	rhs := p.val
	lhs := fmt.Sprint(x.val)
	return rhs == lhs
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

// TODO replace pathIndex
func newPathSet(s string) pathSet {
	p := pathSet{}

	var i int

	r := strings.NewReader(s)
	f := "%d,"

	for r.Len() > 0 {
		n, _ := fmt.Fscanf(r, f, &i)
		if n > 0 {
			p.val = append(p.val, i)
		}
	}

	return p
}

type pathSet struct {
	val []int
}

func (p pathSet) Equal(v pathSel) bool {
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
