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
	p := &Path{p: s, v: v, r: copyZero(v)}
	if err := p.parse(); err != nil {
		return nil, err
	}
	return p, nil
}

func (p *Path) String() string {
	return fmt.Sprintf("%V", p.sel)
}

func copyZero(v interface{}) interface{} {
	switch v.(type) {
	case map[string]interface{}:
		return make(map[string]interface{})
	case []interface{}:
		return make([]interface{}, 0)
	}
	return nil
}
