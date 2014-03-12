package goj

type Path struct {
	p     string
	sel   []interface{}
	depth int
	err   error
}

func NewPath(s string) (*Path, error) {
	p := &Path{p: s}
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

func (d *Decoder) FilterOn(s string) error {
	var arr []interface{}
	p, err := NewPath(s)

	if err != nil {
		return err
	}

	filterPath(d.v, arr, p)

	return nil
}
