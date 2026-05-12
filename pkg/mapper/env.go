package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapEnv(v *hyprast.Env, ctx *MapperContext) (any, error) {
	return &luaast.Call{
		Callee: luaast.MemberChain("hl.env"),
		Args: []luaast.Expr{
			stringify(v.Key),
			stringify(v.Value),
		},
	}, nil
}
