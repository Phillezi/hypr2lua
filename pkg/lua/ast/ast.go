package ast

type Node interface {
	lua()
}

type Stmt interface {
	Node
	stmt()
}

type Expr interface {
	Node
	expr()
}

type File struct {
	Body []Stmt
}

func (*File) lua() {}

type Identifier struct {
	Name string
}

func (*Identifier) lua()  {}
func (*Identifier) expr() {}

type String struct {
	Value string
}

func (*String) lua()  {}
func (*String) expr() {}

type Integer struct {
	Value int64
}

func (*Integer) lua()  {}
func (*Integer) expr() {}

type Float struct {
	Value float64
}

func (*Float) lua()  {}
func (*Float) expr() {}

type Boolean struct {
	Value bool
}

func (*Boolean) lua()  {}
func (*Boolean) expr() {}

type Nil struct{}

func (*Nil) lua()  {}
func (*Nil) expr() {}

type Table struct {
	Fields []Field
}

func (*Table) lua()  {}
func (*Table) expr() {}

type Field struct {
	Key   Expr
	Value Expr
}

func (*Field) lua() {}

type Array struct {
	Values []Expr
}

func (*Array) lua()  {}
func (*Array) expr() {}

type Member struct {
	Base Expr
	Name string
}

func (*Member) lua()  {}
func (*Member) expr() {}

type Index struct {
	Base  Expr
	Index Expr
}

func (*Index) lua()  {}
func (*Index) expr() {}

type Call struct {
	Callee Expr
	Args   []Expr
}

func (*Call) lua()  {}
func (*Call) expr() {}

type Function struct {
	Params []string
	Body   []Stmt
}

func (*Function) lua()  {}
func (*Function) expr() {}

type Binary struct {
	Left  Expr
	Op    string
	Right Expr
}

func (*Binary) lua()  {}
func (*Binary) expr() {}

type Unary struct {
	Op    string
	Value Expr
}

func (*Unary) lua()  {}
func (*Unary) expr() {}

type Assign struct {
	Left  []Expr
	Right []Expr
}

func (*Assign) lua()  {}
func (*Assign) stmt() {}

type Local struct {
	Names  []string
	Values []Expr
}

func (*Local) lua()  {}
func (*Local) stmt() {}

type Return struct {
	Values []Expr
}

func (*Return) lua()  {}
func (*Return) stmt() {}

type ExprStmt struct {
	Expr Expr
}

func (*ExprStmt) lua()  {}
func (*ExprStmt) stmt() {}

type If struct {
	Condition Expr
	Then      []Stmt
	Else      []Stmt
}

func (*If) lua()  {}
func (*If) stmt() {}

type Comment struct {
	Text string
}

func (*Comment) lua()  {}
func (*Comment) stmt() {}

type Concat struct {
	Parts []Expr
}

func (*Concat) lua()  {}
func (*Concat) expr() {}
