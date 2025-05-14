package types

import "fmt"

type parser struct {
	s string
	i int
}

func (p *parser) get(n int) *string {
	if p.i+n > len(p.s) {
		return nil
	}
	s := p.s[p.i : p.i+n]
	p.i += n
	return &s
}

func (p *parser) isAtEnd() bool {
	return p.i >= len(p.s)
}

func (p *parser) readUntil(c byte) string {
	start := p.i
	for !p.isAtEnd() && p.s[p.i] != c {
		p.i++
	}
	result := p.s[start:p.i]
	p.i++
	return result
}

func (p *parser) expect(c byte) error {
	if p.isAtEnd() || p.s[p.i] != c {
		return fmt.Errorf("expected '%c' but got '%c'", c, p.s[p.i])
	}
	p.i++
	return nil
}
