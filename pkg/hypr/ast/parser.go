package ast

import (
	"bytes"
	"fmt"
	"io"
	"strconv"
	"strings"
)

type Parser struct {
	lx *Lexer

	cur  Token
	peek Token
}

func NewParser(r io.Reader) *Parser {
	lx := NewLexer(r)

	p := &Parser{
		lx: lx,
	}

	// prime tokens
	p.next()
	p.next()

	return p
}

func NewParserString(s string) *Parser {
	return NewParser(strings.NewReader(s))
}

func NewParserBytes(b []byte) *Parser {
	return NewParser(bytes.NewReader(b))
}

func (p *Parser) next() {
	p.cur = p.peek
	p.peek = p.lx.NextToken()
}

func (p *Parser) Parse() (*File, error) {
	file := &File{}

	for p.cur.Type != EOF {

		if p.cur.Type == NEWLINE {
			p.next()
			continue
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		if node != nil {
			file.Nodes = append(file.Nodes, node)
		}
	}

	return file, nil
}

func (p *Parser) parseNode() (Node, error) {
	switch p.cur.Type {

	case COMMENT:
		node := &Comment{
			Text: p.cur.Value,
		}
		p.next()
		return node, nil

	case VARIABLE:
		return p.parseVariable()

	case IDENT:
		return p.parseIdentNode()
	}

	return nil, nil
}

func (p *Parser) parseVariable() (Node, error) {
	name := p.cur.Value
	p.next()

	if p.cur.Type != EQUALS {
		return nil, p.error("expected = after variable")
	}

	p.next()

	value, err := p.parseExpr()
	if err != nil {
		return nil, err
	}

	return &Variable{
		Name:  name,
		Value: value,
	}, nil
}

func (p *Parser) parseIdentNode() (Node, error) {
	name := p.cur.Value
	p.next()

	switch p.cur.Type {

	case LBRACE:
		return p.parseSection(name)

	case EQUALS:
		return p.parseAssignmentOrDirective(name)
	}

	return nil, nil
}

func (p *Parser) parseSection(name string) (Node, error) {
	section := &Section{
		Name: name,
	}

	p.next() // consume '{'

	for {

		if p.cur.Type == EOF {
			return nil, p.error("unexpected EOF in section")
		}

		if p.cur.Type == RBRACE {
			p.next() // consume '}'
			break
		}

		if p.cur.Type == NEWLINE {
			p.next()
			continue
		}

		node, err := p.parseNode()
		if err != nil {
			return nil, err
		}

		if node != nil {
			section.Body = append(section.Body, node)
		}
	}

	return section, nil
}

func (p *Parser) parseAssignmentOrDirective(name string) (Node, error) {
	p.next() // consume '='

	args, err := p.parseCSV()
	if err != nil {
		return nil, err
	}

	if isBind(name) {
		return p.parseBind(name, args)
	}

	if name == "source" && len(args) == 1 {
		return &Source{
			Path: args[0],
		}, nil
	}

	if name == "env" && len(args) == 2 {
		return &Env{
			Key:   args[0],
			Value: args[1],
		}, nil
	}

	if name == "exec-once" {
		return &Exec{
			Once: true,
			Command: Concat{
				Parts: args,
			},
		}, nil
	}

	if name == "exec" {
		return &Exec{
			Once: false,
			Command: Concat{
				Parts: args,
			},
		}, nil
	}

	if name == "monitor" && len(args) >= 4 {
		return &Monitor{
			Name:     args[0],
			Mode:     args[1],
			Position: args[2],
			Scale:    args[3],
		}, nil
	}

	if name == "windowrule" {
		rules, action, err := parseMatcher(args)
		if err != nil {
			return nil, err
		}
		return &WindowRule{Action: action, Matches: rules}, nil
	}

	if name == "layerrule" {
		rules, action, err := parseMatcher(args)
		if err != nil {
			return nil, err
		}
		return &LayerRule{Action: action, Matches: rules}, nil
	}

	if isDirective(name) {
		return &Directive{
			Name: name,
			Args: args,
		}, nil
	}

	if len(args) == 1 {
		return &Assignment{
			Key:   name,
			Value: args[0],
		}, nil
	}

	return &Directive{
		Name: name,
		Args: args,
	}, nil
}

func (p *Parser) parseBind(kind string, args []Expr) (Node, error) {
	if len(args) < 3 {
		return nil, p.error("invalid bind (need mods, key, dispatcher)")
	}

	return &Bind{
		Kind: kind,

		Mods: splitMods(args[0]),

		Key:        args[1],
		Dispatcher: args[2],

		Args: args[3:],
	}, nil
}

func (p *Parser) parseCSV() ([]Expr, error) {
	var values []Expr
	var parts []Expr

	flush := func() {
		if len(parts) == 0 {
			values = append(values, Nil{})
			return
		}

		if len(parts) == 1 {
			values = append(values, parts[0])
		} else {
			values = append(values, &Concat{
				Parts: parts,
			})
		}

		parts = nil
	}

	for {
		if p.cur.Type == NEWLINE || p.cur.Type == EOF || p.cur.Type == COMMENT {
			break
		}

		if p.cur.Type == COMMA {
			flush()
			p.next()
			continue
		}

		switch p.cur.Type {

		case VARIABLE:
			parts = append(parts, &VariableRef{
				Name: p.cur.Value,
			})

		case STRING:
			parts = append(parts, &String{
				Value: p.cur.Value,
			})

		case INTEGER:
			v, err := strconv.Atoi(p.cur.Value)
			if err == nil {
				parts = append(parts, &Integer{
					Value: int64(v),
				})
			} else {
				parts = append(parts, &String{
					Value: p.cur.Value,
				})
			}
		case FLOAT:
			v, err := strconv.ParseFloat(p.cur.Value, 64)
			if err == nil {
				parts = append(parts, &Float{
					Value: v,
				})
			} else {
				parts = append(parts, &String{
					Value: p.cur.Value,
				})
			}

		case IDENT:
			// IMPORTANT: preserve raw identifiers INCLUDING ":" and "-"
			parts = append(parts, &String{
				Value: p.cur.Value,
			})

		default:
			// fallback: keep raw token text
			parts = append(parts, &String{
				Value: p.cur.Value,
			})
		}

		p.next()
	}

	flush()

	return values, nil
}

func (p *Parser) parseExpr() (Expr, error) {
	switch p.cur.Type {

	case STRING:
		v := &String{Value: p.cur.Value}
		p.next()
		return v, nil

	case INTEGER:
		n, err := strconv.ParseInt(p.cur.Value, 10, 64)
		if err != nil {
			return nil, err
		}

		p.next()
		return &Integer{Value: n}, nil

	case FLOAT:
		n, err := strconv.ParseFloat(p.cur.Value, 64)
		if err != nil {
			return nil, err
		}

		p.next()
		return &Float{Value: n}, nil

	case BOOLEAN:
		v := &Boolean{Value: p.cur.Value == "true"}
		p.next()
		return v, nil

	case VARIABLE:
		v := &VariableRef{Name: p.cur.Value}
		p.next()
		return v, nil

	case IDENT:
		v := &String{Value: p.cur.Value}
		p.next()
		return v, nil

	case LBRACKET:
		return p.parseArray()
	}

	return nil, p.error("unexpected token in expression")
}

/*func (p *Parser) parseExpr() (Expr, error) {
	switch p.cur.Type {

	case STRING:
		return &String{Value: p.cur.Value}, nil

	case INTEGER:
		n, err := strconv.ParseInt(p.cur.Value, 10, 64)
		if err != nil {
			return nil, err
		}
		return &Integer{Value: n}, nil

	case FLOAT:
		n, err := strconv.ParseFloat(p.cur.Value, 64)
		if err != nil {
			return nil, err
		}
		return &Float{Value: n}, nil

	case BOOLEAN:
		return &Boolean{Value: p.cur.Value == "true"}, nil

	case VARIABLE:
		return &VariableRef{Name: p.cur.Value}, nil

	case IDENT:
		return &String{Value: p.cur.Value}, nil

	case LBRACKET:
		return p.parseArray()
	}

	return nil, p.error("unexpected token in expression")
}*/

func (p *Parser) parseArray() (Expr, error) {
	arr := &Array{}

	p.next() // consume '['

	for {

		if p.cur.Type == RBRACKET {
			p.next() // consume ']'
			break
		}

		if p.cur.Type == COMMA {
			p.next()
			continue
		}

		expr, err := p.parseExpr()
		if err != nil {
			return nil, err
		}

		arr.Items = append(arr.Items, expr)

		p.next()
	}

	return arr, nil
}

/*func (p *Parser) error(msg string) error {
	return fmt.Errorf("%s at %d:%d", msg, p.cur.Line, p.cur.Column)
}*/

func (p *Parser) error(msg string) error {
	return fmt.Errorf(
		"%s at %d:%d (token=%q type=%v)",
		msg,
		p.cur.Line,
		p.cur.Column,
		p.cur.Value,
		p.cur.Type,
	)
}

func isDirective(name string) bool {
	switch name {
	case "input",
		"workspace",
		"gestures",
		"debug",
		"misc":
		return true
	}

	// IMPORTANT:
	// Hypr allows unknown directives too
	// so we treat unknown-but-known-patterns as directive fallback

	return false
}

func isBind(name string) bool {
	switch name {
	case "bind", "binde", "bindm", "bindl", "bindel":
		return true
	}
	return false
}

func splitMods(expr Expr) []Expr {
	var out []Expr

	var walk func(Expr)

	walk = func(e Expr) {
		switch v := e.(type) {

		case *Concat:
			for _, part := range v.Parts {
				walk(part)
			}

		case *String:
			for field := range strings.FieldsSeq(v.Value) {
				out = append(out, &String{
					Value: field,
				})
			}

		case *VariableRef:
			out = append(out, v)

		case Nil:

		default:
			panic(
				"splitMods: unsupported expr: " +
					fmt.Sprintf("%T", v),
			)
		}
	}

	walk(expr)

	return out
}

func parseMatcher(args []Expr) (rules []MatchExpr, action Expr, err error) {
	args = flattenExpr(args)
	for i := 0; i < len(args); i++ {
		arg := args[i]
		switch v := arg.(type) {
		case *String:
			if m, ok := strings.CutPrefix(v.Value, "match:"); ok {
				if i+1 >= len(args) {
					return nil, nil, fmt.Errorf("expected value after matcher rule")
				}
				rules = append(rules, MatchExpr{
					Field: m,
					Value: args[i+1],
				})
				i++
			} else {
				if action != nil {
					if v, ok := action.(*Concat); ok {
						v.Parts = append(v.Parts, arg)
					} else {
						// return nil, fmt.Errorf("multiple actions for rule is not supported, found %v and %v", action, arg)
						tmp := Concat{
							Parts: []Expr{
								action,
								arg,
							},
						}
						action = &tmp
					}
				} else {
					action = arg
				}
			}
		default:
			return nil, nil, fmt.Errorf("unsupported matcher: %T", v)
		}
	}
	return rules, action, nil
}

// flattenExpr takes in a slice of expressions and expands all concats to their types (flattened)
func flattenExpr(args []Expr) []Expr {
	out := make([]Expr, 0, len(args)*2)
	for _, arg := range args {
		switch v := arg.(type) {
		case *Concat:
			for _, p := range v.Parts {
				out = append(out, p)
			}
		case *Array:
			for _, it := range v.Items {
				out = append(out, it)
			}
		default:
			out = append(out, arg)
		}
	}
	return out
}
