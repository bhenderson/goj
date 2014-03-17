package goj

// pathSel is an interface for each path component
type pathSel interface {
	Eval(v interface{}) bool
}

type pathRec struct{}

func (p pathRec) Eval(v interface{}) bool {
	return true
}

func (p pathRec) String() string {
	return "**"
}

type pathParent struct{}

func (p pathParent) Eval(v interface{}) bool {
	return true
}

func (p pathParent) String() string {
	return ".."
}

type pathKey struct {
	val string
}

func (p pathKey) Eval(v interface{}) bool {
	return true
}

type pathVal struct {
	val interface{}
}

func (p pathVal) Eval(v interface{}) bool {
	return true
}

type pathIndex struct {
	i int
}

func (p pathIndex) Eval(v interface{}) bool {
	return true
}
