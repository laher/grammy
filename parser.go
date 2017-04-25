package main

import (
	"fmt"
	"io"
	"strings"
)

type StatementType int

const (
	DefinitionStatement StatementType = iota
	BlankStatement
	MessageStatement
)

// Statement represents a UML definition statement.
type Statement struct {
	Type       StatementType
	Definition *Definition
	Message    *Message
}

type Message struct {
	From        string
	To          string
	Arrow       Token
	Description string
}

type Definition struct {
	Name  string
	Type  Token
	Alias string
}

// Parser represents a parser.
type Parser struct {
	s                 *Scanner
	LastStatementType StatementType
	Definitions       []*Definition
	buf               struct {
		tok Token  // last read token
		lit string // last read literal
		n   int    // buffer size (max=1)
	}
}

// NewParser returns a new instance of Parser.
func NewParser(r io.Reader) *Parser {
	return &Parser{s: NewScanner(r)}
}

var ActorTokens = []Token{ACTOR, PARTICIPANT, DATABASE}
var ArrowTokens = []Token{RARR, LARR, RDARR, LDARR}
var EOLTokens = []Token{NL, EOF}

func isActor(t Token) bool {
	for _, at := range ActorTokens {
		if t == at {
			return true
		}
	}
	return false
}

func isEOL(t Token) bool {
	for _, at := range EOLTokens {
		if t == at {
			return true
		}
	}
	return false
}

func isArrow(t Token) bool {
	for _, at := range ArrowTokens {
		if t == at {
			return true
		}
	}
	return false
}

func (p *Parser) Parse() (*Statement, error) {
	tok, _ := p.scanIgnoreWhitespace()
	var s *Statement
	var err error
	switch tok {
	case ACTOR, PARTICIPANT, DATABASE:
		p.unscan()
		s, err = p.parseDefinition()
	case EOF:
		return nil, io.EOF
	case NL:
		return &Statement{Type: BlankStatement}, nil
	case IDENT:
		p.unscan()
		s, err = p.parseMessage()
	}
	if err != nil {
		return s, err
	}
	fmt.Printf("Parsed: t:%v, d:%#v, m:%v\n", s.Type, s.Definition, s.Message)
	return s, err
}

func (p *Parser) parseMessage() (*Statement, error) {
	// First token should be a name.
	tok, lit := p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected IDENT", lit)
	}
	stmt := &Statement{Type: MessageStatement}
	// First token should be a "ACTOR" keyword.
	m := &Message{}
	stmt.Message = m
	fmt.Printf("Parse\n")
	m.From = lit
	tok, lit = p.scanIgnoreWhitespace()
	if !isArrow(tok) {
		return nil, fmt.Errorf("found %q, expected ARROW", lit)
	}
	m.Arrow = tok

	tok, lit = p.scanIgnoreWhitespace()
	if tok != IDENT {
		return nil, fmt.Errorf("found %q, expected IDENT", lit)
	}
	m.To = lit

	// Next we should see the ":" keyword.
	tok, lit = p.scanIgnoreWhitespace()

	switch tok {
	case COLON:
		// Finally we should read the message.
		for {
			tok, lit := p.scan()
			if tok == WS || tok == IDENT {
				m.Description += lit
			} else {
				break
			}
		}
		m.Description = strings.TrimSpace(m.Description)
		tok, lit = p.scanIgnoreWhitespace()
		switch tok {
		case NL, EOF:
			return stmt, nil
		default:
			return nil, fmt.Errorf("found %q, expected NL or EOF", lit)
		}

	case NL, EOF:
		return stmt, nil
	default:
		return nil, fmt.Errorf("found %q, expected COLON, NL or EOF", lit)
	}
	return stmt, nil
}

func (p *Parser) parseDefinition() (*Statement, error) {
	tok, lit := p.scanIgnoreWhitespace()
	if !isActor(tok) {
		return nil, fmt.Errorf("found %q, expected one of %v", lit, ActorTokens)
	}
	stmt := &Statement{Type: DefinitionStatement}
	// First token should be a "ACTOR" keyword.
	d := &Definition{}
	stmt.Definition = d
	fmt.Printf("Parse\n")

	d.Type = tok

	// Read a field.
	tok, lit = p.scanIgnoreWhitespace()
	if isActor(tok) {
		return nil, fmt.Errorf("found %q, expected name", lit)
	}
	d.Name = lit
	// ok this is worth keeping now
	p.Definitions = append(p.Definitions, d)

	// Next we should see the "AS" keyword.
	tok, lit = p.scanIgnoreWhitespace()

	switch tok {
	case AS:
		// Finally we should read the alias.
		tok, lit := p.scanIgnoreWhitespace()
		if tok != IDENT {
			return nil, fmt.Errorf("found %q, expected alias name", lit)
		}
		d.Alias = lit
		tok, lit = p.scanIgnoreWhitespace()

		switch tok {
		case NL, EOF:
			return stmt, nil
		default:
			return nil, fmt.Errorf("found %q, expected NL or EOF", lit)
		}

	case NL, EOF:
		return stmt, nil
	default:
		return nil, fmt.Errorf("found %q, expected AS, NL or EOF", lit)
	}
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
