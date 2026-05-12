package mapper

import (
	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapComment(s *hyprast.Comment, ctx *MapperContext) (any, error) {
	return []luaast.Stmt{&luaast.Comment{Text: s.Text}}, nil
}
