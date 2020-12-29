package parser

// AnonymousClass represents an anonymous Class.Check function.
type AnonymousClass func(p *Parser) (*Cursor, bool)

// CheckRune returns an AnonymousClass that checks whether the current rune of
// the parser matches the given rune. The same result can be achieved by using
// p.Expect(r). Where 'p' is a reference to the parser an 'r' a rune value.
func CheckRune(r rune) AnonymousClass {
	return func(p *Parser) (*Cursor, bool) {
		return p.Mark(), p.Current() == r
	}
}

// CheckRuneFunc returns an AnonymousClass that checks whether the current rune of
// the parser matches the given validator.
func CheckRuneFunc(f func(r rune) bool) AnonymousClass {
	return func(p *Parser) (*Cursor, bool) {
		return p.Mark(), f(p.Current())
	}
}

// CheckString returns an AnonymousClass that checks whether the current
// sequence runes of the parser matches the given string. The same result can be
// achieved by using p.Expect(s). Where 'p' is a reference to the parser an 's'
// a string value.
func CheckString(s string) AnonymousClass {
	return func(p *Parser) (*Cursor, bool) {
		var last *Cursor
		for _, r := range []rune(s) {
			if p.Current() != r {
				return nil, false
			}
			last = p.Mark()
			p.Next()
		}
		return last, true
	}
}

// Class provides an interface for checking classes.
type Class interface {
	// Check should return the last p.Mark() that matches the class. It should
	// also return whether it was able to check the whole class.
	//
	// e.g. if the class is defined as follows: '<=' / '=>'. Then a parser that
	// only contains '=' will not match this class and return 'nil, false'.
	Check(p *Parser) (*Cursor, bool)
}
