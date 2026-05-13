package ast

import (
	"bufio"
	"bytes"
	"io"
	"slices"
	"strings"
	"unicode"
)

type Lexer struct {
	r *bufio.Reader

	line int
	col  int
}

func NewLexer(r io.Reader) *Lexer {
	return &Lexer{
		r:    bufio.NewReader(r),
		line: 1,
		col:  0,
	}
}

func NewLexerString(s string) *Lexer {
	return NewLexer(strings.NewReader(s))
}

func NewLexerBytes(b []byte) *Lexer {
	return NewLexer(bytes.NewReader(b))
}

func (l *Lexer) read() rune {
	ch, _, err := l.r.ReadRune()
	if err != nil {
		return 0
	}

	if ch == '\n' {
		l.line++
		l.col = 0
	} else {
		l.col++
	}

	return ch
}

func (l *Lexer) unread() {
	_ = l.r.UnreadRune()
}

func (l *Lexer) peek() rune {
	ch := l.read()
	if ch != 0 {
		l.unread()
	}
	return ch
}

func (l *Lexer) token(t TokenType, v string) Token {
	return Token{
		Type:   t,
		Value:  v,
		Line:   l.line,
		Column: l.col,
	}
}

func (l *Lexer) skipWhitespace() {
	for {
		ch := l.peek()

		if ch == ' ' || ch == '\t' || ch == '\r' {
			l.read()
			continue
		}

		break
	}
}

func (l *Lexer) NextToken() Token {
	l.skipWhitespace()

	ch := l.peek()

	switch ch {

	case 0:
		return l.token(EOF, "")

	case '\n':
		l.read()
		return l.token(NEWLINE, "\n")

	case '#':
		return l.lexComment()

	case '{':
		l.read()
		return l.token(LBRACE, "{")

	case '}':
		l.read()
		return l.token(RBRACE, "}")

	case '[':
		l.read()
		return l.token(LBRACKET, "[")

	case ']':
		l.read()
		return l.token(RBRACKET, "]")

	case '(':
		l.read()
		return l.token(LPAREN, "(")

	case ')':
		l.read()
		return l.token(RPAREN, ")")

	case ',':
		l.read()
		return l.token(COMMA, ",")

	case '=':
		l.read()
		return l.token(EQUALS, "=")

	case '"':
		return l.lexString()

	case '$':
		return l.lexVariable()

	case '~', '/':
		return l.lexPath()
	}

	/*if unicode.IsDigit(ch) {
		// This is super inefficient
		if !l.isBareStartingWithNumber() {
			return l.lexNumber()
		} else {
			return l.lexBare()
		}
	}*/

	if isIdentStart(ch) {
		return l.lexIdent()
	}

	// l.read()
	// return l.token(ILLEGAL, string(ch))
	return l.lexBare()
}

func (l *Lexer) isBareStartingWithNumber() bool {
	count := 0

	// consume leading digits
	var ch rune
	for {
		ch, _, _ = l.r.ReadRune()
		count++
		if !unicode.IsDigit(ch) && ch != '.' {
			break
		}
	}
	// rewind everything we consumed
	for ; count > 0; count-- {
		_ = l.r.UnreadRune()
	}

	// if the next character after digits is not a separator,
	// then this token should be treated as a bare string
	switch ch {
	case 0, ' ', '\t', '\n', '#', ',', '}', ']', ')':
		return false
	}

	return true
}

func (l *Lexer) lexComment() Token {
	startLine := l.line
	startCol := l.col

	l.read() // consume '#'

	var b strings.Builder

	for {
		ch := l.peek()

		if ch == '\n' || ch == 0 {
			break
		}

		b.WriteRune(l.read())
	}

	return Token{
		Type:   COMMENT,
		Value:  b.String(),
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) lexString() Token {
	startLine := l.line
	startCol := l.col

	l.read() // "

	var b strings.Builder

	for {
		ch := l.peek()

		if ch == '"' || ch == 0 { // end on "
			break
		}

		if ch == '\\' {
			l.read()        // skip \
			esc := l.read() // skip the rune after \
			b.WriteRune(esc)
			continue
		}

		b.WriteRune(l.read())
	}

	if l.peek() == '"' { // remove final "
		l.read()
	}

	return Token{
		Type:   STRING,
		Value:  b.String(),
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) lexVariable() Token {
	startLine := l.line
	startCol := l.col

	l.read() // $

	var b strings.Builder

	for {
		ch := l.peek()

		if !(unicode.IsLetter(ch) ||
			unicode.IsDigit(ch) ||
			ch == '_' ||
			ch == '-') {
			break
		}

		b.WriteRune(l.read())
	}

	return Token{
		Type:   VARIABLE,
		Value:  b.String(),
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) lexPath() Token {
	startLine := l.line
	startCol := l.col

	var b strings.Builder

	for {
		ch := l.peek()

		if ch == 0 || ch == '\n' {
			break
		}

		// stop only on structural separators
		if ch == ' ' || ch == '\t' || ch == ',' || ch == '#' {
			break
		}

		b.WriteRune(l.read())
	}

	return Token{
		Type:   STRING,
		Value:  b.String(),
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) lexNumber() Token {
	startLine := l.line
	startCol := l.col

	var b strings.Builder

	hasDot := false

	for {
		ch := l.peek()

		if unicode.IsDigit(ch) {
			b.WriteRune(l.read())
			continue
		}

		if ch == '.' && !hasDot {
			hasDot = true
			b.WriteRune(l.read())
			continue
		}

		break
	}

	typ := INTEGER
	if hasDot {
		typ = FLOAT
	}

	return Token{
		Type:   typ,
		Value:  b.String(),
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) lexIdent() Token {
	startLine := l.line
	startCol := l.col

	var b strings.Builder

	for {
		ch := l.peek()

		if !isIdentPart(ch) {
			break
		}

		b.WriteRune(l.read())
	}

	value := b.String()

	if value == "yes" || value == "no" {
		return Token{
			Type: BOOLEAN,
			Value: func(value string) string {
				if value == "yes" {
					return "true"
				}
				return "false"
			}(value),
			Line:   startLine,
			Column: startCol,
		}
	}

	if value == "true" || value == "false" {
		return Token{
			Type:   BOOLEAN,
			Value:  value,
			Line:   startLine,
			Column: startCol,
		}
	}

	return Token{
		Type:   IDENT,
		Value:  value,
		Line:   startLine,
		Column: startCol,
	}
}

func (l *Lexer) ReadRawUntil(delims ...rune) string {
	var b strings.Builder

	for {
		ch := l.peek()

		if ch == 0 {
			break
		}

		stop := slices.Contains(delims, ch)

		if stop {
			break
		}

		b.WriteRune(l.read())
	}

	return strings.TrimSpace(b.String())
}

func (l *Lexer) lexBare() Token {
	startLine := l.line
	startCol := l.col

	var b strings.Builder

	for {
		ch := l.peek()

		// stop only at structural boundaries
		if ch == 0 ||
			ch == '\n' ||
			ch == ',' ||
			ch == '=' ||
			ch == '{' ||
			ch == '}' ||
			ch == '[' ||
			ch == ']' ||
			ch == '"' ||
			ch == '#' {

			break
		}

		b.WriteRune(l.read())
	}

	str := strings.TrimSpace(b.String())
	tType := INTEGER
	for _, c := range str {
		if !unicode.IsDigit(c) {
			if c == '.' {
				tType = FLOAT
			} else {
				tType = IDENT
				break
			}
		}
	}

	return Token{
		Type:   tType, // or STRING (better later)
		Value:  str,
		Line:   startLine,
		Column: startCol,
	}
}

func isIdentStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_'
}

func isIdentPart(ch rune) bool {
	return unicode.IsLetter(ch) ||
		unicode.IsDigit(ch) ||
		ch == '_' ||
		ch == '-' || ch == '+' || ch == ':' || ch == '@' || ch == '%' || ch == '.'
}
