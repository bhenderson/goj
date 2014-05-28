package goj

import (
	"errors"
	"fmt"
	"path/filepath"
	"regexp"
	"strings"
)

var uEOF = "unexpected EOF"

// pathSel is an interface for each path component
type pathSel interface {
	Equal(l *Leaf) bool
}

// pathRec implements pathSel and stands for **
type pathRec struct{}

func (p *pathRec) Equal(l *Leaf) bool { return true }
func (p *pathRec) String() string     { return "**" }

// pathParent implements pathSel and stands for ..
type pathParent struct{}

func (p *pathParent) Equal(l *Leaf) bool { return true }
func (p *pathParent) String() string     { return ".." }

type pathKey struct {
	val string
}

var special = map[rune]struct{}{
	'{': struct{}{},
	'}': struct{}{},
	'(': struct{}{},
	')': struct{}{},
	'|': struct{}{},
	'$': struct{}{},
	'+': struct{}{},
	'*': struct{}{},
}

func isReg(s string) bool {
	for _, c := range s {
		if _, ok := special[c]; ok {
			return true
		}
	}
	return false
}

func (p *pathKey) Equal(l *Leaf) bool {
	var rhs string
	if l.kind == leafKey {
		rhs = l.val.(string)
	} else {
		if l.kind != leafIdx {
			return false
		}
		rhs = fmt.Sprint(l.val)
	}

	lhs := p.val
	if lhs == rhs {
		return true
	}
	ok, err := filepath.Match(lhs, rhs)
	if ok && err == nil {
		return true
	}
	// TODO performance hit?
	if isReg(lhs) {
		re, err := regexp.Compile("^" + lhs + "$")
		if err != nil {
			return false
		}
		if re.MatchString(rhs) {
			return true
		}
	}
	return false
}

type pathVal struct {
	val string
}

func (p *pathVal) Equal(l *Leaf) bool {
	if l.kind != leafVal {
		return false
	}

	if l.val == nil {
		return p.val == ""
	}

	// type assertion?
	lhs := p.val
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

	rhs := l.val.(int)
	t := l.max

	for _, i := range p.val {
		if i < 0 {
			i = i + t
		}
		if i == rhs {
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
	rhs := l.val.(int)
	t := l.max
	b, e, s := p.b, p.e, p.s
	// reverse index
	if b < 0 {
		b = b + t
	}
	if e < 0 {
		e = e + t
	}
	if b > rhs || rhs > e {
		return false
	}
	// step
	if s > 1 && (rhs-b)%s != 0 {
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
		if e != nil && e.Error() != uEOF {
			return &pathIndex{}, false
		}
		if n > 0 {
			sel.val = append(sel.val, i)
		}
	}

	return sel, ok
}
