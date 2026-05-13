package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapExec(v *hyprast.Exec, ctx *MapperContext) (any, error) {
	exec := &luaast.Call{
		Callee: luaast.MemberChain("hl.exec_cmd"),
		Args: []luaast.Expr{
			convertExpr(&v.Command),
		},
	}
	if v.Once {
		exec = &luaast.Call{
			Callee: luaast.MemberChain("hl.on"),
			Args: []luaast.Expr{
				luaast.Str("hyprland.start"),
				&luaast.Function{
					Body: []luaast.Stmt{
						&luaast.Comment{
							Text: "Hi you, yes you! This conversion kinda sucks, please manually add all former exec-once into a function block like this",
						},
						&luaast.ExprStmt{
							Expr: exec,
						},
					},
				},
			},
		}
	}

	return exec, nil
}
