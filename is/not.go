package is

// Not represents a negation of the Value. This should not consume data.
// e.g. Not{'a'} should check if the first rune is not an 'a'.
type Not struct {
	Value interface{}
}
