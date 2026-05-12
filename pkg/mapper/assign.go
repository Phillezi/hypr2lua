package mapper

import (
	"fmt"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapAssignment(a *hyprast.Assignment, ctx *MapperContext) (any, error) {
	path := a.Key
	if ctx.Section != "" {
		path = ctx.Section + "." + path
	}

	return LuaAssign(path, a.Value)
}

func LuaAssign(path string, value any) (luaast.Stmt, error) {
	var rhs luaast.Expr

	switch v := value.(type) {

	case string:
		rhs = luaast.Str(v)

	case int:
		rhs = luaast.Int(int64(v))

	case int64:
		rhs = luaast.Int(v)

	case float64:
		rhs = new(luaast.Float{Value: v})

	case bool:
		rhs = luaast.Bool(v)

	case *hyprast.String:
		rhs = luaast.Str(v.Value)

	case *hyprast.Integer:
		rhs = luaast.Int(v.Value)

	case *hyprast.Float:
		rhs = new(luaast.Float{Value: v.Value})

	case *hyprast.Boolean:
		rhs = luaast.Bool(v.Value)

	case *hyprast.Concat:
		if len(v.Parts) == 0 {
			rhs = luaast.Str("")
			break
		}

		// start from first element
		expr := convertExpr(v.Parts[0])

		for i := 1; i < len(v.Parts); i++ {
			expr = &luaast.Binary{
				Left:  expr,
				Op:    "..",
				Right: convertExpr(v.Parts[i]),
			}
		}

		rhs = expr

	case hyprast.Nil:
		rhs = nil

	case *hyprast.Section:
		return &luaast.ExprStmt{Expr: &luaast.Call{
			Callee: luaast.MemberChain("hl.config"),
			Args: []luaast.Expr{&luaast.Table{
				Fields: []luaast.Field{
					{
						Key: luaast.MemberChain(path),
						Value: &luaast.Table{
							Fields: convertSectionBody(v.Body),
						},
					},
				},
			}},
		}}, nil

	default:
		return nil, fmt.Errorf("unsupported value type: %T", value)
	}

	return &luaast.Assign{
		Left: []luaast.Expr{
			luaast.MemberChain("hl.config." + path),
		},
		Right: []luaast.Expr{rhs},
	}, nil
}
