package stubs

import (
	"bufio"
	"io"
	"strings"
)

type Parser struct {
	schema *Schema
}

func NewParser() *Parser {
	return &Parser{
		schema: &Schema{
			Globals: map[string]Field{},
			Modules: map[string]*Module{},
		},
	}
}

func (p *Parser) Parse(r io.Reader) (*Schema, error) {
	scanner := bufio.NewScanner(r)

	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "--") {
			continue
		}

		p.parseLine(line)
	}

	return p.schema, scanner.Err()
}

func (p *Parser) parseLine(line string) {
	// ignore type hints for now
	if strings.HasPrefix(line, "---@type") {
		return
	}

	if !strings.Contains(line, "=") {
		return
	}

	parts := strings.SplitN(line, "=", 2)

	left := strings.TrimSpace(parts[0])

	path := strings.Split(left, ".")

	p.registerPath(path)
}

func (p *Parser) registerPath(path []string) {
	if len(path) < 2 {
		return
	}

	if path[0] != "hl" {
		return
	}

	module := path[1]

	if _, ok := p.schema.Modules[module]; !ok {
		p.schema.Modules[module] = &Module{
			Name:   module,
			Fields: map[string]Field{},
		}
	}

	mod := p.schema.Modules[module]

	full := strings.Join(path[2:], ".")

	mod.Fields[full] = Field{
		Name: full,
		Path: strings.Join(path, "."),
		Type: "unknown", // enriched later by type hints
	}
}
