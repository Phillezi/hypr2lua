package mapper

import (
	"log"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapBind(v *hyprast.Bind, ctx *MapperContext) (any, error) {
	log.Println("BIND DEBUG___________________")
	log.Println("Key:\t", stringify(v.Key).Value)
	log.Println("Mods:\t", all(all(v.Mods, stringify), func(in *luaast.String) string {
		return in.Value
	}))
	log.Println("Args:\t", all(all(v.Args, stringify), func(in *luaast.String) string {
		return in.Value
	}))
	log.Println("Dispatcher:\t", stringify(v.Dispatcher).Value)
	log.Println("Kind:\t", v.Kind)

	mods := make([]luaast.Expr, len(v.Mods))
	for i, m0 := range v.Mods {
		mods[i] = convertExpr(m0)
	}

	key := convertExpr(v.Key)
	dispatcher := convertExpr(v.Dispatcher)

	args := []luaast.Expr{&luaast.Concat{
		Parts: append(mods, key),
	}, dispatcher}

	for _, a := range v.Args {
		args = append(args, convertExpr(a))
	}

	return luaast.ExprStmt{
		Expr: &luaast.Call{
			Callee: luaast.MemberChain("hl.bind"),
			Args:   args,
		},
	}, nil
}
