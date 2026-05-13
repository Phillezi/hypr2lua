package mapper

import (
	"fmt"
	"log"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func (m *Mapper) mapBind(v *hyprast.Bind, ctx *MapperContext) (any, error) {
	log.Println("BIND DEBUG___________________")
	log.Println("Mods:\t", all(all(v.Mods, stringify), func(in *luaast.String) string {
		return in.Value
	}))
	log.Println("Key:\t", stringify(v.Key).Value)
	log.Println("Dispatcher:\t", stringify(v.Dispatcher).Value)
	log.Println("Args:\t", all(all(v.Args, stringify), func(in *luaast.String) string {
		return in.Value
	}))
	log.Println("Kind:\t", v.Kind)

	mods := make([]luaast.Expr, len(v.Mods))
	for i, m0 := range v.Mods {
		mods[i] = convertExpr(m0)
	}

	key := convertExpr(v.Key)
	// dispatcher := convertExpr(v.Dispatcher)

	args := []luaast.Expr{&luaast.Concat{
		Parts: append(mods, key),
	}}

	dsp, ok := v.Dispatcher.(*hyprast.String)
	if !ok {
		return nil, fmt.Errorf("invalid dispatcher, expected *hyprast.String, got %T", v.Dispatcher)
	}
	switch dsp.Value {
	case "exec":
		eArgs := make([]luaast.Expr, 0, len(v.Args))
		for _, a := range v.Args {
			eArgs = append(eArgs, convertExpr(a))
		}
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.exec_cmd"),
			Args:   eArgs,
		})
	case "killactive":
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.window.close"),
		})
	case "exit":
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.exit"),
		})
	case "fullscreen":
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.window.fullscreen"),
		})
	case "togglefloating":
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.window.float"),
			Args: []luaast.Expr{&luaast.Table{
				Fields: []luaast.Field{
					{
						Key:   luaast.MemberChain("action"),
						Value: luaast.Str("toggle"),
					},
				},
			}},
		})
	case "movefocus":
		eArgs := make([]luaast.Expr, 0, len(v.Args))
		for _, a := range v.Args {
			eArgs = append(eArgs, &luaast.Table{Fields: []luaast.Field{{Key: luaast.MemberChain("direction"), Value: convertExpr(a)}}})
		}
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.focus"),
			Args:   eArgs,
		})
	case "workspace":
		eArgs := make([]luaast.Expr, 0, len(v.Args))
		for _, a := range v.Args {
			eArgs = append(eArgs, &luaast.Table{Fields: []luaast.Field{{Key: luaast.MemberChain("workspace"), Value: convertExpr(a)}}})
		}
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.focus"),
			Args:   eArgs,
		})
	case "movetoworkspace":
		eArgs := make([]luaast.Expr, 0, len(v.Args))
		for _, a := range v.Args {
			eArgs = append(eArgs, &luaast.Table{Fields: []luaast.Field{{Key: luaast.MemberChain("workspace"), Value: convertExpr(a)}}})
		}
		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.window.move"),
			Args:   eArgs,
		})
	case "layoutmsg":
		fields := make([]luaast.Field, 0, len(v.Args))
		for _, a := range v.Args {
			switch vv := a.(type) {
			case *hyprast.String:
				switch vv.Value {
				case "swapwithmaster":
					fields = append(fields, luaast.Field{
						Key:   luaast.MemberChain("master"),
						Value: luaast.Bool(true),
					})
				default:
					return nil, fmt.Errorf("unimplemented layoutmsg %q", vv.Value)
				}
			}
		}

		args = append(args, &luaast.Call{
			Callee: luaast.MemberChain("hl.dsp.window.swap"),
			Args: []luaast.Expr{&luaast.Table{
				Fields: fields,
			}},
		})

	case "movewindow":

	case "resizewindow":

	}

	return luaast.ExprStmt{
		Expr: &luaast.Call{
			Callee: luaast.MemberChain("hl.bind"),
			Args:   args,
		},
	}, nil
}
