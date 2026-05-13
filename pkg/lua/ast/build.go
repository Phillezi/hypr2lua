package ast

import (
	"strconv"
	"strings"
)

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

func Float64(v float64) *Float {
	return &Float{
		Value: v,
	}
}

func Bool(v bool) *Boolean {
	return &Boolean{
		Value: v,
	}
}

func Auto[T any](v T) Expr {
	switch v := any(v).(type) {

	case string:
		s := strings.TrimSpace(v)

		// bool
		if b, err := strconv.ParseBool(s); err == nil {
			return Bool(b)
		}

		// int
		if i, err := strconv.ParseInt(s, 10, 64); err == nil {
			return Int(i)
		}

		// float
		if f, err := strconv.ParseFloat(s, 64); err == nil {
			return Float64(f)
		}

		// fallback string
		return Str(s)

	case bool:
		return Bool(v)

	case int:
		return Int(int64(v))

	case int8:
		return Int(int64(v))

	case int16:
		return Int(int64(v))

	case int32:
		return Int(int64(v))

	case int64:
		return Int(v)

	case uint:
		return Int(int64(v))

	case uint8:
		return Int(int64(v))

	case uint16:
		return Int(int64(v))

	case uint32:
		return Int(int64(v))

	case uint64:
		return Int(int64(v))

	case float32:
		return Float64(float64(v))

	case float64:
		return Float64(v)

	default:
		panic("unsupported type")
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
