package compiler

import (
	"fmt"
	"os"

	"github.com/phillezi/hypr2lua/pkg/hypr/ast"
	"github.com/phillezi/hypr2lua/pkg/lua/emitter"
	"github.com/phillezi/hypr2lua/pkg/lua/stubs"
	"github.com/phillezi/hypr2lua/pkg/mapper"
)

type Compiler struct {
	mapper *mapper.Mapper
	schema *stubs.Schema
}

func New(schema *stubs.Schema) *Compiler {
	m := mapper.New()
	return &Compiler{
		mapper: m,
		schema: schema,
	}
}

func (c *Compiler) Compile(file *ast.File) error {
	o, err := c.mapper.MapFile(file, c.schema)
	if err != nil {
		return err
	}

	if o == nil {
		return fmt.Errorf("compiled to nil")
	}

	if err := emitter.Emit(os.Stdout, o); err != nil {
		return err
	}
	return err
}
