package ast

type TokenType int

const (
	EOF TokenType = iota

	ILLEGAL

	IDENT
	STRING
	INTEGER
	FLOAT
	BOOLEAN

	VARIABLE

	LBRACE
	RBRACE
	LBRACKET
	RBRACKET

	LPAREN
	RPAREN

	COMMA
	COLON
	EQUALS

	NEWLINE
	COMMENT
)

type Token struct {
	Type  TokenType
	Value string

	Line   int
	Column int
}
