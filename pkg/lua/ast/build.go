package ast

import "strings"

func Ident(name string) *Identifier {
	return &Identifier{
		Name: name,
	}
}

func Str(v string) *String {
	return &String{
		Value: v,
	}
}

func Int(v int64) *Integer {
	return &Integer{
		Value: v,
	}
}

func Bool(v bool) *Boolean {
	return &Boolean{
		Value: v,
	}
}

func MemberChain(path string) Expr {
	parts := strings.Split(path, ".")

	var expr Expr = Ident(parts[0])

	for _, p := range parts[1:] {
		expr = &Member{
			Base: expr,
			Name: p,
		}
	}

	return expr
}

func CallExpr(
	name string,
	args ...Expr,
) *Call {
	return &Call{
		Callee: MemberChain(name),
		Args:   args,
	}
}
