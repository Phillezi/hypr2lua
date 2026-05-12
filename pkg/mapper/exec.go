package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapExec(v *hyprast.Exec, ctx *MapperContext) (any, error) {
	return &luaast.Call{
		Callee: luaast.MemberChain("hl.exec_cmd"),
		Args: []luaast.Expr{
			convertExpr(&v.Command),
		},
	}, nil
}
