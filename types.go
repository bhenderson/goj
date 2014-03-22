package goj

import (
	// "log"
	"fmt"
	"path/filepath"
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
