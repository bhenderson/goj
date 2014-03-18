package goj

type Path struct {
	p   string
	sel []pathSel
	err error
	v   interface{}
}

func NewPath(s string, v interface{}) (*Path, error) {
	p := &Path{p: s, v: v}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p, nil
}

type Pair struct {
	key string
	val interface{}
}

type PairSlice struct {
	b, e, s interface{}
}
