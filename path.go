package goj

import "fmt"

type Path struct {
	// the original string, useful for constructing errors.
	p string
	// an error, if any. Useful in the state function.
	err error
	// the array of path elements.
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
