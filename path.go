package goj

import "fmt"

type Path struct {
	p   string
	err error
	sel []pathSel
}

func NewPath(s string) (*Path, error) {
	p := &Path{p: s}
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
		return map[string]interface{}{}
	case []interface{}:
		return []interface{}{}
	}
	return nil
}
