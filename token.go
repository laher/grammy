package main

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF
	WS
	NL

	// Literals
	IDENT // main

	// Misc characters
	ASTERISK  // *
	COMMA     // ,
	SEMICOLON // ;
	COLON     // :

	RARR  // ->
	RDARR // -->
	LARR  // <-
	LDARR // <--

	AS

	// Keywords
	ACTOR
	PARTICIPANT
	DATABASE
)
