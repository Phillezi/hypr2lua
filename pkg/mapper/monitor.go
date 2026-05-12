package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapMonitor(v *hyprast.Monitor, ctx *MapperContext) (any, error) {
	return &luaast.Call{
		Callee: luaast.MemberChain("hl.monitor"),
		Args: []luaast.Expr{
			&luaast.Table{
				Fields: []luaast.Field{
					{
						Key:   luaast.MemberChain("output"),
						Value: stringify(v.Name),
					},

					{
						Key:   luaast.MemberChain("mode"),
						Value: stringify(v.Mode),
					},

					{
						Key:   luaast.MemberChain("position"),
						Value: stringify(v.Position),
					},

					{
						Key:   luaast.MemberChain("scale"),
						Value: stringify(v.Scale),
					},
				},
			},
		},
	}, nil
}
