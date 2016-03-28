package ql

// Token represents a lexical token.
type Token int

const (
	// Special tokens
	ILLEGAL Token = iota
	EOF           // end of file
	WS            // white space

	// Literals
	IDENT // main
	DIGIT

	// Misc characters
	ASTERISK   // *
	COMMA      // ,
	COMPARITOR // =,<.>,!=

	// Keywords
	SELECT
	FROM
	WHERE
	AND
)
