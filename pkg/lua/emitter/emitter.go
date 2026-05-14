package emitter

import (
	"fmt"
	"io"
	"strings"

	"github.com/phillezi/hypr2lua/pkg/lua/ast"
)

type Emitter struct {
	w      io.Writer
	indent int
}

func New(w io.Writer) *Emitter {
	return &Emitter{
		w: w,
	}
}

func Emit(w io.Writer, file *ast.File) error {
	e := New(w)
	return e.emitFile(file)
}

func EmitString(file *ast.File) (string, error) {
	var buf strings.Builder

	err := Emit(&buf, file)
	if err != nil {
		return "", err
	}

	return buf.String(), nil
}

func (e *Emitter) emitFile(f *ast.File) error {
	for _, stmt := range f.Body {

		if err := e.emitStmt(stmt); err != nil {
			return err
		}

		if _, err := e.w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	return nil
}

func (e *Emitter) writeIndent() error {
	for i := 0; i < e.indent; i++ {
		if _, err := e.w.Write([]byte("    ")); err != nil {
			return err
		}
	}

	return nil
}

func (e *Emitter) emitStmt(s ast.Stmt) error {
	switch n := s.(type) {

	case *ast.Assign:
		return e.emitAssign(n)

	case *ast.Local:
		return e.emitLocal(n)

	case *ast.Return:
		return e.emitReturn(n)

	case *ast.ExprStmt:
		if err := e.writeIndent(); err != nil {
			return err
		}
		return e.emitExpr(n.Expr)

	case *ast.Comment:
		if err := e.writeIndent(); err != nil {
			return err
		}

		_, err := e.w.Write(
			[]byte("-- " + n.Text),
		)

		return err
	default:
		return fmt.Errorf("unsupported stmt type %T", n)
	}
}

func (e *Emitter) emitAssign(a *ast.Assign) error {
	if err := e.writeIndent(); err != nil {
		return err
	}

	for i, left := range a.Left {

		if i > 0 {
			if _, err := e.w.Write([]byte(", ")); err != nil {
				return err
			}
		}

		if err := e.emitExpr(left); err != nil {
			return err
		}
	}

	if _, err := e.w.Write([]byte(" = ")); err != nil {
		return err
	}

	for i, right := range a.Right {

		if i > 0 {
			if _, err := e.w.Write([]byte(", ")); err != nil {
				return err
			}
		}

		if err := e.emitExpr(right); err != nil {
			return err
		}
	}

	return nil
}

func (e *Emitter) emitLocal(l *ast.Local) error {
	if err := e.writeIndent(); err != nil {
		return err
	}

	if _, err := e.w.Write([]byte("local ")); err != nil {
		return err
	}

	for i, name := range l.Names {

		if i > 0 {
			if _, err := e.w.Write([]byte(", ")); err != nil {
				return err
			}
		}

		if _, err := io.WriteString(e.w, name); err != nil {
			return err
		}
	}

	if len(l.Values) > 0 {

		if _, err := e.w.Write([]byte(" = ")); err != nil {
			return err
		}

		for i, v := range l.Values {

			if i > 0 {
				if _, err := e.w.Write([]byte(", ")); err != nil {
					return err
				}
			}

			if err := e.emitExpr(v); err != nil {
				return err
			}
		}
	}

	return nil
}

func (e *Emitter) emitReturn(r *ast.Return) error {
	if err := e.writeIndent(); err != nil {
		return err
	}

	if _, err := e.w.Write([]byte("return ")); err != nil {
		return err
	}

	for i, v := range r.Values {

		if i > 0 {
			if _, err := e.w.Write([]byte(", ")); err != nil {
				return err
			}
		}

		if err := e.emitExpr(v); err != nil {
			return err
		}
	}

	return nil
}

func (e *Emitter) emitExpr(expr ast.Expr) error {
	switch n := expr.(type) {

	case *ast.Identifier:
		_, err := io.WriteString(e.w, n.Name)
		return err

	case *ast.String:
		_, err := fmt.Fprintf(e.w, "%q", n.Value)
		return err

	case *ast.Integer:
		_, err := fmt.Fprintf(e.w, "%d", n.Value)
		return err

	case *ast.Float:
		_, err := fmt.Fprintf(e.w, "%f", n.Value)
		return err

	case *ast.Boolean:
		if n.Value {
			_, err := e.w.Write([]byte("true"))
			return err
		}
		_, err := e.w.Write([]byte("false"))
		return err
	case *ast.Binary:
		if err := e.emitExpr(n.Left); err != nil {
			return err
		}

		// surround operator with spaces for readability
		if _, err := fmt.Fprintf(e.w, " %s ", n.Op); err != nil {
			return err
		}

		if err := e.emitExpr(n.Right); err != nil {
			return err
		}

		return nil

	case *ast.Nil:
		_, err := e.w.Write([]byte("nil"))
		return err

	case *ast.Member:
		if err := e.emitExpr(n.Base); err != nil {
			return err
		}
		_, err := e.w.Write([]byte("." + n.Name))
		return err

	case *ast.Call:
		if err := e.emitExpr(n.Callee); err != nil {
			return err
		}

		if _, err := e.w.Write([]byte("(")); err != nil {
			return err
		}

		for i, a := range n.Args {

			if i > 0 {
				if _, err := e.w.Write([]byte(", ")); err != nil {
					return err
				}
			}

			if err := e.emitExpr(a); err != nil {
				return err
			}
		}

		_, err := e.w.Write([]byte(")"))
		return err

	case *ast.Table:
		return e.emitTable(n)
	case *ast.Array:
		clen := len(n.Values)
		for i, child := range n.Values {
			if err := e.emitExpr(child); err != nil {
				return err
			}
			if i < clen-1 {
				if _, err := e.w.Write([]byte(", ")); err != nil {
					return err
				}
			}
		}
		return nil
	case *ast.Concat:
		clen := len(n.Parts)
		for i, part := range n.Parts {
			if err := e.emitExpr(part); err != nil {
				return err
			}
			if i < clen-1 {
				if _, err := e.w.Write([]byte(" .. ")); err != nil {
					return err
				}
			}
		}
		return nil

	case *ast.Function:
		if _, err := e.w.Write([]byte("function(")); err != nil {
			return err
		}
		clen := len(n.Params)
		for i, param := range n.Params {
			if _, err := e.w.Write([]byte(param)); err != nil {
				return err
			}
			if i < clen-1 {
				if _, err := e.w.Write([]byte(", ")); err != nil {
					return err
				}
			}
		}
		if _, err := e.w.Write([]byte(")\n")); err != nil {
			return err
		}
		e.indent++
		for _, stmt := range n.Body {
			if err := e.emitStmt(stmt); err != nil {
				return err
			}
			if _, err := e.w.Write([]byte("\n")); err != nil {
				return err
			}
		}
		e.indent--
		if _, err := e.w.Write([]byte("end")); err != nil {
			return err
		}
		return nil

	default:
		if n == nil { // skip nil
			return nil
		}
		return fmt.Errorf("unsupported expr type %T", n)

	}
}

func (e *Emitter) emitTable(t *ast.Table) error {
	if _, err := e.w.Write([]byte("{")); err != nil {
		return err
	}

	if len(t.Fields) == 0 {
		_, err := e.w.Write([]byte("}"))
		return err
	}

	if _, err := e.w.Write([]byte("\n")); err != nil {
		return err
	}

	e.indent++

	for i, f := range t.Fields {

		if err := e.writeIndent(); err != nil {
			return err
		}

		if err := e.emitExpr(f.Key); err != nil {
			return err
		}

		if _, err := e.w.Write([]byte(" = ")); err != nil {
			return err
		}

		if err := e.emitExpr(f.Value); err != nil {
			return err
		}

		if i < len(t.Fields)-1 {
			if _, err := e.w.Write([]byte(",")); err != nil {
				return err
			}
		}

		if _, err := e.w.Write([]byte("\n")); err != nil {
			return err
		}
	}

	e.indent--

	if err := e.writeIndent(); err != nil {
		return err
	}

	_, err := e.w.Write([]byte("}"))

	return err
}
