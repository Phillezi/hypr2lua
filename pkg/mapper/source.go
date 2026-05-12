package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapSource(s *hyprast.Source, ctx *MapperContext) (any, error) {
	return luaast.ExprStmt{
		Expr: &luaast.Call{
			Callee: luaast.MemberChain("require"),
			Args: []luaast.Expr{
				// TODO: remove trailing .conf
				convertExpr(s.Path),
			},
		},
	}, nil
}
