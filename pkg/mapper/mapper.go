package mapper

import (
	"fmt"
	"log"
	"strings"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
	"github.com/phillezi/hypr2lua/pkg/lua/stubs"
)

type MapperContext struct {
	Section string
	Vars    map[string]string
	Meta    map[string]any
	File    *hyprast.File
	Schema  *stubs.Schema
}

func NewContext(f *hyprast.File, schema *stubs.Schema) *MapperContext {
	return &MapperContext{
		File:   f,
		Vars:   make(map[string]string),
		Meta:   make(map[string]any),
		Schema: schema,
	}
}

func (c *MapperContext) PushSection(name string) string {
	prev := c.Section
	c.Section = name
	return prev
}

func (c *MapperContext) PopSection(prev string) {
	c.Section = prev
}

type Mapper struct{}

func New() *Mapper {
	return &Mapper{}
}

func (m *Mapper) MapFile(f *hyprast.File, schema *stubs.Schema) (*luaast.File, error) {
	ctx := NewContext(f, schema)

	out := &luaast.File{}

	for _, n := range f.Nodes {
		res, err := m.MapNode(n, ctx)
		if err != nil {
			return nil, err
		}

		if res == nil {
			continue
		}

		if err := appendLua(out, res); err != nil {
			return nil, err
		}
	}

	return out, nil
}

func (m *Mapper) MapNode(n hyprast.Node, ctx *MapperContext) (any, error) {
	switch v := n.(type) {
	case *hyprast.Comment:
		return m.mapComment(v, ctx)
	case *hyprast.Variable:
		return m.mapVariable(v, ctx)
	case *hyprast.Assignment:
		return m.mapAssignment(v, ctx)
	case *hyprast.Section:
		return m.mapSection(v, ctx)
	case *hyprast.Directive:
		return m.mapDirective(v, ctx)
	case *hyprast.Bind:
		return m.mapBind(v, ctx)
	case *hyprast.Source:
		return m.mapSource(v, ctx)
	case *hyprast.Env:
		return m.mapEnv(v, ctx)
	case *hyprast.Exec:
		return m.mapExec(v, ctx)
	case *hyprast.Monitor:
		return m.mapMonitor(v, ctx)
	case *hyprast.WindowRule:
		return m.mapWindowRule(v, ctx)
	case *hyprast.LayerRule:
		return m.mapLayerRule(v, ctx)
	default:
		return nil, fmt.Errorf("unsupported node: %T", n)
	}
}

func (m *Mapper) mapDirective(d *hyprast.Directive, ctx *MapperContext) (any, error) {
	switch d.Name {
	case "workspace":
		wArgs := workspaceArgs(d.Args)
		return &luaast.Call{
			Callee: luaast.MemberChain("hl.workspace_rule"),
			Args: []luaast.Expr{
				&luaast.Table{
					Fields: wArgs,
				},
			},
		}, nil
	default:
		log.Printf("Unimplemented mapping %q", d.Name)
		return &luaast.Nil{}, nil
	}
}

func appendLua(f *luaast.File, n any) error {
	switch v := n.(type) {

	case luaast.Stmt:
		f.Body = append(f.Body, v)
		return nil

	case []luaast.Stmt:
		f.Body = append(f.Body, v...)
		return nil

	case luaast.Expr:
		f.Body = append(f.Body, &luaast.ExprStmt{Expr: v})
		return nil

	case luaast.ExprStmt:
		f.Body = append(f.Body, &v)
		return nil

	default:
		return fmt.Errorf("cannot append lua node: %T", n)
	}
}

func directiveArgs(args []hyprast.Expr) []luaast.Expr {
	out := make([]luaast.Expr, 0, len(args))

	for _, a := range args {
		out = append(out, convertExpr(a))
	}

	return out
}

func workspaceArgs(args []hyprast.Expr) []luaast.Field {
	fields := make([]luaast.Field, 0, len(args))

	for i, a := range args {
		switch i {
		case 0:
			fields = append(fields, luaast.Field{
				Key:   luaast.MemberChain("workspace"),
				Value: stringify(a),
			})
		default:
			switch v := a.(type) {
			case *hyprast.String:
				parts := strings.Split(v.Value, ":")
				if len(parts) > 1 {
					if len(parts) == 2 {
						var value luaast.Expr = luaast.Auto(parts[1])
						fields = append(fields, luaast.Field{
							Key:   luaast.MemberChain(parts[0]),
							Value: value,
						})
					} else {
						panic(fmt.Errorf("unexpected workspace rule arg len, expected len to be 2, got %d", len(parts)))
					}
				} else {
					panic(fmt.Errorf("unsupported workspace rule arg, %q", v.Value))
				}
			default:
				panic(fmt.Errorf("unsupported workspace arg type %T", v))
			}
		}
	}
	return fields
}
