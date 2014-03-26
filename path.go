package goj

import "fmt"

type Path struct {
	p   string
	v   interface{}
	r   interface{}
	err error
	sel []pathSel
	res []pathSel
}

func NewPath(s string, v interface{}) (*Path, error) {
	p := &Path{p: s, v: v}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Path) String() string {
	return fmt.Sprintf("%V", p.sel)
}
