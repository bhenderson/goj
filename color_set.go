package goj

type badColor struct{}

func (b *badColor) Error() string {
	return "color must be one of always, auto, never"
}

// Implements flag.Var interface
type colorSet string

const (
	ColorAlways colorSet = "always"
	ColorAuto   colorSet = "auto"
	ColorNever  colorSet = "never"
)

var Colors = [...]colorSet{ColorAlways, ColorAuto, ColorNever}

func (c *colorSet) Set(s string) error {
	x := colorSet(s)
	switch x {
	case ColorAlways, ColorAuto, ColorNever:
		*c = x
		return nil
	}

	return &badColor{}
}
func (c *colorSet) Get() interface{} { return c }
func (c *colorSet) String() string   { return string(*c) }
func (c *colorSet) IsBoolFlag() bool { return true }
