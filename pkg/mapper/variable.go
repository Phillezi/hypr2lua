package mapper

import (
	"fmt"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapVariable(v *hyprast.Variable, ctx *MapperContext) (any, error) {
	return LuaMkVar(v.Name, v.Value)
}

func LuaMkVar(name string, value any) (luaast.Stmt, error) {
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

	case hyprast.Nil:
		return nil, fmt.Errorf("cannot make var with nil assignment")

	default:
		return nil, fmt.Errorf("unsupported value type: %T", value)
	}
	return &luaast.Local{
		Names:  []string{name},
		Values: []luaast.Expr{rhs},
	}, nil
}
