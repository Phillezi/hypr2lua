package mapper

import (
	"strings"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapSection(s *hyprast.Section, ctx *MapperContext) (any, error) {
	// FIXME: ehhh, i tried to fix the col.<blablabla> nested issue
	parts := strings.Split(s.Name, ".")
	sect := parts[0]
	ws := s
	if len(parts) > 1 {
		ws = &hyprast.Section{
			Name: parts[len(parts)-1],
			Body: ws.Body,
		}
		for i := len(parts) - 2; i > 0; i-- {
			p := parts[i]
			ws = &hyprast.Section{
				Name: p,
				Body: []hyprast.Node{ws},
			}
		}
	}
	prev := ctx.PushSection(sect)
	defer ctx.PopSection(prev)

	return &luaast.ExprStmt{Expr: &luaast.Call{
		Callee: luaast.MemberChain("hl.config"),
		Args: []luaast.Expr{&luaast.Table{
			Fields: []luaast.Field{
				{
					Key: luaast.MemberChain(sect),
					Value: &luaast.Table{
						Fields: convertSectionBody(ws.Body),
					},
				},
			},
		}},
	}}, nil
}
