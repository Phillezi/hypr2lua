package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapLayerRule(v *hyprast.LayerRule, ctx *MapperContext) (any, error) {
	return &luaast.Call{
		Callee: luaast.MemberChain("hl.layer_rule"),
		Args: []luaast.Expr{
			&luaast.Table{
				Fields: []luaast.Field{
					{
						Key: luaast.MemberChain("match"),
						Value: &luaast.Table{
							Fields: convertMatchers(v.Matches),
						},
					},
					convertAction(v.Action),
				},
			},
		},
	}, nil
}
