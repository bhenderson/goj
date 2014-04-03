package goj

import (
	"errors"
	"fmt"
	// "log"
	"path/filepath"
	"strings"
)

var EOF = "EOF"
var UEOF = "unexpected EOF"

// pathSel is an interface for each path component
type pathSel interface {
	Equal(l *Leaf) bool
}

type pathRec struct{}

func (p *pathRec) Equal(l *Leaf) bool {
	return true
}

func (p *pathRec) String() string {
	return "**"
}

type pathParent struct{}

func (p *pathParent) Equal(l *Leaf) bool {
	return true
}

func (p *pathParent) String() string {
	return ".."
}

type pathKey struct {
	val string
}

func (p *pathKey) Equal(l *Leaf) bool {
	var rhs string
	if l.kind != leafKey {
		if l.kind != leafIdx {
			return false
		}
		rhs = fmt.Sprint(l.val)
	} else {
		rhs = l.val.(string)
	}
	if p.val == rhs {
		return true
	}
	ok, err := filepath.Match(p.val, rhs)
	return ok && err == nil
}

type pathVal struct {
	val interface{}
}

func (p *pathVal) Equal(l *Leaf) bool {
	if l.kind != leafVal {
		return false
	}
	// type assertion?
	lhs := fmt.Sprint(p.val)
	rhs := fmt.Sprint(l.val)
	if lhs == rhs {
		return true
	}
	ok, err := filepath.Match(lhs, rhs)
	return ok && err == nil
}

// original index value of []interface{}
type pathIdx struct {
	val int
	len int
}

func (p *pathIdx) Equal(l *Leaf) bool {
	return true
}

// jsonpath index value
type pathIndex struct {
	val []int
}

func (p *pathIndex) Equal(l *Leaf) bool {
	if l.kind != leafIdx {
		return false
	}

	for _, i := range p.val {
		if i == l.val.(int) {
			return true
		}
	}

	return false
}

type pathSlice struct {
	b, e, s int
}

func (p *pathSlice) Equal(l *Leaf) bool {
	if l.kind != leafIdx {
		return false
	}
	lhs := l.val.(int)
	t := l.max
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

// common case of *
var pathStar = &pathKey{"*"}

// common case of [*]
var pathSliceAll = &pathSlice{0, -1, 1}

func newPathIndex(s string) (sel pathSel, err error) {
	if len(s) == 0 {
		return nil, errors.New("array index cannot be empty")
	}

	if s == "*" || s == ":" || s == "::" {
		return pathSliceAll, nil
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
func newPathSlice(r *strings.Reader) (sel *pathSlice, ok bool) {
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

	return &pathSlice{b, e, s}, true
}

func newPathSet(r *strings.Reader) (*pathIndex, bool) {
	sel := &pathIndex{}
	var n, i int
	var e error

	d := "%d"
	c := ","
	ok := true

	for r.Len() > 0 {
		n, e = fmt.Fscanf(r, d, &i)
		if n == 0 || e != nil {
			return &pathIndex{}, false
		}
		_, e = fmt.Fscanf(r, c)
		if e != nil && e.Error() != UEOF {
			return &pathIndex{}, false
		}
		if n > 0 {
			sel.val = append(sel.val, i)
		}
	}

	return sel, ok
}
