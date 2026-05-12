package mapper

import (
	"testing"

	hyprast "github.com/phillezi/hypr2lua/pkg/hypr/ast"
	luaast "github.com/phillezi/hypr2lua/pkg/lua/ast"
)

func TestEnv(t *testing.T) {
	m := New()
	r, _ := m.mapEnv(&hyprast.Env{
		Key:   &hyprast.String{Value: "foo"},
		Value: &hyprast.String{Value: "bar"},
	}, nil)
	if v, ok := r.(*luaast.Call); ok {
		t.Logf("%v", *v)
	} else {
		t.Fatalf("this should not happen, invalid test")
	}
}

func TestExec(t *testing.T) {
	m := New()
	r, _ := m.mapExec(&hyprast.Exec{
		Once: true,
		Command: hyprast.Concat{
			Parts: []hyprast.Expr{
				&hyprast.String{Value: "go"},
				&hyprast.String{Value: "test"},
				&hyprast.String{Value: "./..."},
			},
		},
	}, nil)
	if v, ok := r.(*luaast.Call); ok {
		t.Logf("%v", *v)
	} else {
		t.Fatalf("this should not happen, invalid test")
	}
}

func TestMonitor(t *testing.T) {
	m := New()
	r, _ := m.mapEnv(&hyprast.Env{
		Key:   &hyprast.String{Value: "foo"},
		Value: &hyprast.String{Value: "bar"},
	}, nil)
	if v, ok := r.(*luaast.Call); ok {
		t.Logf("%v", *v)
	} else {
		t.Fatalf("this should not happen, invalid test")
	}
}
