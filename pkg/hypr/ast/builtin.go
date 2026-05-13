package ast

type Bind struct {
	Kind       string
	Mods       []Expr
	Key        Expr
	Dispatcher Expr
	Args       []Expr
}

func (*Bind) node() {}

type Monitor struct {
	Name     Expr
	Mode     Expr
	Position Expr
	Scale    Expr
}

func (*Monitor) node() {}

type Exec struct {
	Once    bool
	Command Concat
}

func (*Exec) node() {}

type Env struct {
	Key   Expr
	Value Expr
}

func (*Env) node() {}

type MatchExpr struct {
	Field string // class, title, xwayland, fullscreen, pin, etc.
	Value Expr   // regex / number / bool / identifier
}

func (*MatchExpr) expr() {}

type WindowRule struct {
	Action  Expr        // "no_focus on", "float", "monitor DP-2"
	Matches []MatchExpr // multiple match conditions
}

func (*WindowRule) node() {}

type LayerRule struct {
	Action  Expr
	Matches []MatchExpr
}

func (*LayerRule) node() {}
