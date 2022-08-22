package typecheck

import (
	"fmt"
	"iracema/ast"
	"strings"
)

type Type interface {
	Name() string
	Is(Type) bool
}

type ErrList []error

func (el ErrList) Error() string {
	b := new(strings.Builder)

	for _, err := range el {
		// remove INVALID errors
		if strings.Index(err.Error(), "INVALID") > 0 {
			continue
		}

		b.WriteString(err.Error())
	}

	return b.String()
}

type gamma struct {
	parent *gamma
	env    map[string]Type
}

func (s *gamma) lookup(name string) Type {
	if s.parent == nil {
		return s.env[name]
	}

	if typ := s.parent.lookup(name); typ != nil {
		return typ
	}

	return s.env[name]
}

type typechecker struct {
	*gamma

	file *ast.File
	errs ErrList
}

func (tc *typechecker) check() ErrList {
	tc.checkStatement(tc.file.Stmts)

	return tc.errs
}
func (tc *typechecker) checkStatement(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		switch s := stmt.(type) {

		case *ast.VarDecl:
			tc.checkLetDecl(s)

		case *ast.AssignStmt:
			tc.checkAssignStmt(s)

		default:
			fmt.Printf("%+v\n", s)
			panic("unreacheble")
		}
	}
}

func (tc *typechecker) checkLetDecl(let *ast.VarDecl) {
	name := let.Name.Value

	if let.Type == nil && let.Value == nil {
		tc.env[name] = INVALID
		return
	}

	var typ Type
	if let.Type != nil {
		typ = tc.checkExpr(let.Type)
		tc.env[name] = typ

		if let.Value != nil {
			value := tc.checkExpr(let.Value)
			if !value.Is(typ) {
				tc.errorf("cannot use %s as %s value in assignment", value.Name(), typ.Name())
			}
		}

		return
	}

	tc.env[name] = tc.checkExpr(let.Value)
}

func (tc *typechecker) checkAssignStmt(assign *ast.AssignStmt) {
	if len(assign.Left) != len(assign.Right) {
		tc.errorf("assignment mismatch: %d variables but %d values", len(assign.Left), len(assign.Right))
		return
	}

	size := len(assign.Left)
	for i := 0; i < size; i++ {
		rhs := assign.Right[i]
		lhs := assign.Left[i]

		lhsType := tc.checkExpr(lhs)
		rhsType := tc.checkExpr(rhs)

		if !rhsType.Is(lhsType) {
			tc.errorf("cannot use %s as %s value in assignment", rhsType.Name(), lhsType.Name())
		}
	}
}

func (tc *typechecker) checkExpr(expr ast.Expr) Type {
	switch node := expr.(type) {
	case *ast.BasicLit:
		if litType := LIT_TYPES[node.Token.Type.String()]; litType != nil {
			return litType
		}

		return INVALID
	case *ast.Ident:
		if t := tc.lookupType(node.Value); t != nil {
			return t
		}

		tc.errorf("undefined: %s", node.Value)
		return INVALID

	default:
		fmt.Printf("%+v\n", node)
		panic("unreacheble")
	}
}

func (tc *typechecker) lookupType(name string) Type {
	if k, ok := LIT_TYPES[name]; ok {
		return k
	}

	if typ := tc.lookup(name); typ != nil {
		return typ
	}

	return nil
}

func (tc *typechecker) errorf(format string, args ...any) {
	err := fmt.Errorf(format, args...)
	tc.errs = append(tc.errs, err)
}

func Check(file *ast.File) ErrList {
	tc := new(typechecker)
	tc.file = file
	tc.gamma = &gamma{
		env: make(map[string]Type),
	}

	return tc.check()
}
