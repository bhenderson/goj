package goj

// pathSel is an interface for each path component
type pathSel interface {
	Eval(v interface{}) bool
}

type pathRec struct{}

func (p pathRec) Eval(v interface{}) bool {
	return true
}

type pathParent struct{}

func (p pathParent) Eval(v interface{}) bool {
	return true
}

type pathKey struct {
	key string
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
