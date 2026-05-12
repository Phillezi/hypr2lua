package ast

type Node interface {
	node()
}

type File struct {
	Nodes []Node
}

func (*File) node() {}

type Section struct {
	Name string
	Body []Node
}

func (*Section) node() {}

type Assignment struct {
	Key   string
	Value Expr
}

func (*Assignment) node() {}

type Directive struct {
	Name string
	Args []Expr
}

func (*Directive) node() {}

type Bind struct {
	Kind       string
	Mods       []Expr
	Key        Expr
	Dispatcher Expr
	Args       []Expr
}

func (*Bind) node() {}

type Variable struct {
	Name  string
	Value Expr
}

func (*Variable) node() {}

type Source struct {
	Path Expr
}

func (*Source) node() {}

type Comment struct {
	Text string
}

func (*Comment) node() {}

type Expr interface {
	expr()
}

type String struct {
	Value string
}

func (*String) expr() {}

type Integer struct {
	Value int64
}

func (*Integer) expr() {}

type Float struct {
	Value float64
}

func (*Float) expr() {}

type Boolean struct {
	Value bool
}

func (*Boolean) expr() {}

type VariableRef struct {
	Name string
}

func (*VariableRef) expr() {}

type Array struct {
	Items []Expr
}

func (*Array) expr() {}

type Concat struct {
	Parts []Expr
}

func (*Concat) expr() {}

type Nil struct{}

func (Nil) expr() {}
