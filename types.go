package goj

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
	if x, ok := v.(pathKey); ok {
		return p.val == x.val
	}
	return false
}

type pathVal struct {
	val interface{}
}

func (p pathVal) Equal(v pathSel) bool {
	if x, ok := v.(pathVal); ok {
		return p.val == x.val
	}
	return false
}

type pathIndex struct {
	i int
}

func (p pathIndex) Equal(v pathSel) bool {
	return true
}
