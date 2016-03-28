package ql

import (
	"fmt"
	"io"
)

// SelectStatement represents a QL SELECT statement.
type SelectStatement struct {
	Fields []string
	Store  string
	Comps  [][]string
}

// Parser represents a parser.
type Parser struct {
	s   *Scanner
	buf struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

// Parse parses a SQL SELECT statement.
func (p *Parser) Parse() (*SelectStatement, error) {
	stmt := &SelectStatement{}

	// First token should be a "SELECT" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != SELECT {
		return nil, fmt.Errorf("SELECT: found %q, expected SELECT", lit)
	}

	// Next we should loop over all our comma-delimited fields.
	for {
		// Read a field.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT && tok != ASTERISK {
			return nil, fmt.Errorf("FIELDS: found %q, expected field", lit)
		}
		stmt.Fields = append(stmt.Fields, lit)

		// If the next token is not a comma then break the loop.
		if tok, _ := p.scanIgnoreWhitespace(); tok != COMMA {
			p.unscan()
			break
		}
	}

	// Next we should see the "FROM" keyword.
	if tok, lit := p.scanIgnoreWhitespace(); tok != FROM {
		return nil, fmt.Errorf("FROM: found %q, expected FROM", lit)
	}

	// Then we should read the table name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("Store: found %q, expected namespace", lit)
	}
	stmt.Store = lit

	// Next we should look for the possibility of a "WHERE" keyword
	if tok, _ := p.scanIgnoreWhitespace(); tok == WHERE {
		// we should loop through our where clause
		for {
			var comp []string
			// Read a field.
			tok, lit := p.scanIgnoreWhitespace()
			if tok != IDENT && tok != ASTERISK {
				return nil, fmt.Errorf("WHERE: found %q, expected (1st) field", lit)
			}
			comp = append(comp, lit)
			// If the next token is not a comparitor then break the loop.
			tok, lit = p.scanIgnoreWhitespace()
			if tok != COMPARITOR {
				p.unscan()
				break
			}
			comp = append(comp, lit)
			// Read a field.
			tok, lit = p.scanIgnoreWhitespace()
			if tok != IDENT && tok != ASTERISK {
				fmt.Errorf("%q is %+v\n", lit, tok)
				return nil, fmt.Errorf("WHERE: found %q, expected (2nd) field", lit)
			}
			stmt.Comps = append(stmt.Comps, append(comp, lit))
			// If the next token is not a comma, or an and then break the loop.
			if tok, _ = p.scanIgnoreWhitespace(); tok != COMMA && tok != AND {
				p.unscan()
				break
			}
		}
	}

	// Return the successfully parsed statement.
	return stmt, nil
}

// scan returns the next token from the underlying scanner.
// If a token has been unscanned then read that instead.
func (p *Parser) scan() (tok Token, lit string) {
	// If we have a token on the buffer, then return it.
	if p.buf.n != 0 {
		p.buf.n = 0
		return p.buf.tok, p.buf.lit
	}

	// Otherwise read the next token from the scanner.
	tok, lit = p.s.Scan()

	// Save it to the buffer in case we unscan later.
	p.buf.tok, p.buf.lit = tok, lit

	return
}

// scanIgnoreWhitespace scans the next non-whitespace token.
func (p *Parser) scanIgnoreWhitespace() (tok Token, lit string) {
	tok, lit = p.scan()
	if tok == WS {
		tok, lit = p.scan()
	}
	return
}

// unscan pushes the previously read token back onto the buffer.
func (p *Parser) unscan() { p.buf.n = 1 }
