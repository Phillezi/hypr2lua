package mapper

import (
	"fmt"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func convertExpr(e hyprast.Expr) luaast.Expr {
	switch v := e.(type) {

	case *hyprast.String:
		return luaast.Str(v.Value)

	case *hyprast.VariableRef:
		return luaast.MemberChain(v.Name)

	case *hyprast.Integer:
		return luaast.Int(v.Value)

	case *hyprast.Float:
		return &luaast.Float{Value: v.Value}

	case *hyprast.Boolean:
		return luaast.Bool(v.Value)

	case *hyprast.Array:
		expr := luaast.Array{
			Values: make([]luaast.Expr, 0, len(v.Items)),
		}
		for _, p := range v.Items {
			expr.Values = append(expr.Values, convertExpr(p))
		}
		return &expr

	case *hyprast.Concat:
		parts := make([]luaast.Expr, 0, len(v.Parts))
		for _, p := range v.Parts {
			parts = append(parts, convertExpr(p))
		}
		return &luaast.Concat{Parts: parts}

	}

	return luaast.Str(fmt.Sprintf("%v", e))
}

func stringify(expr hyprast.Expr) *luaast.String {
	switch v := expr.(type) {

	case *hyprast.String:
		return luaast.Str(v.Value)

	case *hyprast.VariableRef:
		return luaast.Str(v.Name)

	case *hyprast.Integer:
		return luaast.Str(fmt.Sprintf("%d", v.Value))

	case *hyprast.Float:
		return luaast.Str(fmt.Sprintf("%f", v.Value))

	case *hyprast.Boolean:
		return luaast.Str(func(b bool) string {
			if b {
				return "true"
			}
			return "false"
		}(v.Value))

	case *hyprast.Array:
		str := luaast.String{}
		for _, p := range v.Items {
			str.Value += stringify(p).Value
		}
		return &str

	case *hyprast.Concat:
		var str luaast.String
		for _, p := range v.Parts {
			str.Value += stringify(p).Value
		}
		return &str
	case hyprast.Nil:
		return luaast.Str("<nil>")
	default:
		panic("unimplemented stringify for type: " + fmt.Sprintf("%T", v))
	}
}

func convertNode(node hyprast.Node) luaast.Node {
	switch v := node.(type) {
	case *hyprast.Assignment:
		return &luaast.Assign{
			Left:  []luaast.Expr{luaast.MemberChain(v.Key)},
			Right: []luaast.Expr{convertExpr(v.Value)},
		}
	case *hyprast.Comment:
		return &luaast.Comment{
			Text: v.Text,
		}
	case *hyprast.Directive:
		return &luaast.Field{
			Key: luaast.MemberChain(v.Name),
			Value: &luaast.Concat{
				Parts: all(v.Args, convertExpr),
			},
		}
	case *hyprast.Section:
		return &luaast.Field{
			Key: luaast.MemberChain(v.Name),
			Value: &luaast.Table{
				Fields: convertSectionBody(v.Body),
			},
		}
	default:
		panic(fmt.Errorf("unsupported node conv %T", v))
	}
}

func convertSectionBody(nodes []hyprast.Node) []luaast.Field {
	var out []luaast.Field

	for _, n := range nodes {
		o := convertNode(n)
		switch v := o.(type) {
		case *luaast.Field:
			if v != nil {
				out = append(out, *v)
			}
		case *luaast.Assign:
			out = append(out, luaast.Field{
				Key: &luaast.Concat{
					Parts: v.Left,
				},
				Value: &luaast.Concat{
					Parts: v.Right,
				},
			})
		case *luaast.Comment:
			// comments will be ignored
		default:
			panic(fmt.Errorf("unsupported section body node %T", v))
		}
	}

	return out
}

func all[IN any, OUT any](in []IN, c func(in IN) OUT) []OUT {
	var out []OUT = make([]OUT, len(in))
	for i, input := range in {
		out[i] = c(input)
	}
	return out
}

func convertMatchers(matchers []hyprast.MatchExpr) []luaast.Field {
	out := make([]luaast.Field, len(matchers))
	for i, matcher := range matchers {
		out[i] = luaast.Field{
			Key:   luaast.MemberChain(matcher.Field),
			Value: convertExpr(matcher.Value),
		}
	}
	return out
}

func convertAction(action hyprast.Expr) luaast.Field {
	var key string
	var value luaast.Expr
	switch v := action.(type) {
	case *hyprast.Concat:
		key = stringify(v.Parts[0]).Value
		value = convertExpr(v.Parts[1])
	case *hyprast.String:
		key = v.Value
		value = luaast.Bool(true)
	}
	return luaast.Field{
		Key:   luaast.MemberChain(key),
		Value: value,
	}
}
