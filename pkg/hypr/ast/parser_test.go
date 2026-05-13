package ast_test

import (
	"bytes"
	"strconv"
	"strings"
	"testing"

	"github.com/phillezi/hypr2lua/pkg/hypr/ast"
)

func TestEnv(t *testing.T) {
	envs := []string{
		"env= foo,bar",
		"env =HELLO,world",
		"env=bar,foo",
	}

	r := bytes.NewReader([]byte(strings.Join(envs, "\n")))
	p := ast.NewParser(r)

	f, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Nodes) != len(envs) {
		t.Fatalf("expected %d nodes, got %d", len(envs), len(f.Nodes))
	}

	for i, n := range f.Nodes {
		env, ok := n.(*ast.Env)
		if !ok {
			t.Fatalf("node %d: expected *ast.Env, got %T", i, n)
		}

		parts := strings.Split(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(envs[i], "env")), "=")), ",")
		if len(parts) != 2 {
			t.Fatalf("invalid test input: %s", envs[i])
		}
		keyIdent, ok := env.Key.(*ast.String)
		if !ok {
			t.Fatalf("node %d key: expected *ast.Ident, got %T", i, env.Key)
		}
		if keyIdent.Value != parts[0] {
			t.Errorf("node %d key mismatch: got %q want %q", i, keyIdent.Value, parts[0])
		}
		valIdent, ok := env.Value.(*ast.String)
		if !ok {
			t.Fatalf("node %d value: expected *ast.Ident, got %T", i, env.Value)
		}
		if valIdent.Value != parts[1] {
			t.Errorf("node %d value mismatch: got %q want %q", i, valIdent.Value, parts[1])
		}
	}
}

func TestExec(t *testing.T) {
	execs := []string{
		"exec-once= foot --server",
		"exec-once = tinybar",
		"exec-once=firefox &",
	}

	r := bytes.NewReader([]byte(strings.Join(execs, "\n")))
	p := ast.NewParser(r)

	f, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Nodes) != len(execs) {
		t.Fatalf("expected %d nodes, got %d", len(execs), len(f.Nodes))
	}

	for i, n := range f.Nodes {
		exec, ok := n.(*ast.Exec)
		if !ok {
			t.Fatalf("node %d: expected *ast.Exec, got %T", i, n)
		}

		cmd := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(execs[i], "exec-once")), "="))
		_ = cmd
		if exec.Once == false {
			t.Fatalf("expected exec.Once to be true but it was false")
		}
		if 1 != len(exec.Command.Parts) {
			t.Logf("FAIL, command parts output block:")
			for i, p := range exec.Command.Parts {
				t.Logf("\t%d:%v", i, p)
			}
			t.Fatalf("on %s: expected command parts to be equal to test input, test case: %d, command: %d", execs[i], 1, len(exec.Command.Parts))
		}

	}
}

func TestMonitor(t *testing.T) {
	monitors := []string{
		"monitor=,preferred,auto,auto",
		"monitor=DP-2,2560x1440@164.8,0x0,1",
		"monitor=DP-1,1920x1080@143.85,2560x0,1,",
	}

	r := bytes.NewReader([]byte(strings.Join(monitors, "\n")))
	p := ast.NewParser(r)

	f, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Nodes) != len(monitors) {
		t.Fatalf("expected %d nodes, got %d", len(monitors), len(f.Nodes))
	}

	for i, n := range f.Nodes {
		monitor, ok := n.(*ast.Monitor)
		if !ok {
			t.Fatalf("node %d: expected *ast.Monitor, got %T", i, n)
		}

		parts := strings.Split(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(monitors[i], "monitor")), "=")), ",")
		if len(parts) < 4 {
			t.Fatalf("invalid test input: %s", monitors[i])
		}

		_ = monitor

	}
}

func TestWindowRule(t *testing.T) {
	rules := []string{
		"windowrule = match:class ^(?!(foot|code)),no_blur on",
	}

	r := bytes.NewReader([]byte(strings.Join(rules, "\n")))
	p := ast.NewParser(r)

	f, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Nodes) != len(rules) {
		t.Fatalf("expected %d nodes, got %d", len(rules), len(f.Nodes))
	}

	for i, n := range f.Nodes {
		windowRule, ok := n.(*ast.WindowRule)
		if !ok {
			t.Fatalf("node %d: expected *ast.WindowRule, got %T", i, n)
		}

		parts := strings.Split(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(rules[i], "windowrule")), "=")), ",")
		if len(parts) < 1 {
			t.Fatalf("invalid test input: %s", rules[i])
		}

		_ = windowRule

	}
}

func TestLayerRule(t *testing.T) {
	rules := []string{
		"layerrule = match:namespace tinybar, blur on",
	}

	r := bytes.NewReader([]byte(strings.Join(rules, "\n")))
	p := ast.NewParser(r)

	f, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Nodes) != len(rules) {
		t.Fatalf("expected %d nodes, got %d", len(rules), len(f.Nodes))
	}

	for i, n := range f.Nodes {
		layerRule, ok := n.(*ast.LayerRule)
		if !ok {
			t.Fatalf("node %d: expected *ast.LayerRule, got %T", i, n)
		}

		parts := strings.Split(strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(strings.TrimPrefix(rules[i], "layerrule")), "=")), ",")
		if len(parts) < 1 {
			t.Fatalf("invalid test input: %s", rules[i])
		}

		_ = layerRule

	}
}

// FIXME: make sure it passes for the other var types too
func TestVariable(t *testing.T) {
	vars := []string{
		"$myvar = test123",
		//"$mystr = hello world",
		//"$myint = 67",
		//"$mybool = true",
		//"$myfloat = 6.7",
	}

	r := bytes.NewReader([]byte(strings.Join(vars, "\n")))
	p := ast.NewParser(r)

	f, err := p.Parse()
	if err != nil {
		t.Fatal(err)
	}

	if len(f.Nodes) != len(vars) {
		t.Fatalf("expected %d nodes, got %d", len(vars), len(f.Nodes))
	}

	for i, n := range f.Nodes {
		variable, ok := n.(*ast.Variable)
		if !ok {
			t.Fatalf("node %d: expected *ast.Variable, got %T", i, n)
		}

		parts := strings.Split(vars[i], "=")
		value := strings.TrimSpace(parts[len(parts)-1])

		switch v := variable.Value.(type) {
		case *ast.String:
			if v.Value != value {
				t.Fatalf("incorrect variable value, expected: %q, got: %q", value, v.Value)
			}
		case *ast.Boolean:
			b, err := strconv.ParseBool(value)
			if err != nil {
				t.Fatalf("bad test case, failed to parse bool value: %s", err.Error())
			}
			if v.Value != b {
				t.Fatalf("incorrect variable value, expected: %v, got: %v", b, v.Value)
			}
		case *ast.Integer:
			b, err := strconv.ParseInt(value, 10, 64)
			if err != nil {
				t.Fatalf("bad test case, failed to parse int value: %s", err.Error())
			}
			if v.Value != b {
				t.Fatalf("incorrect variable value, expected: %d, got: %d", b, v.Value)
			}
		case *ast.Float:
			b, err := strconv.ParseFloat(value, 64)
			if err != nil {
				t.Fatalf("bad test case, failed to parse float value: %s", err.Error())
			}
			if v.Value != b {
				t.Fatalf("incorrect variable value, expected: %f, got: %f", b, v.Value)
			}
		default:
			t.Fatalf("unhandled test case type for variable, %T is not impl", v)
		}
	}
}
