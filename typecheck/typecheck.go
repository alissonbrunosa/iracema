package typecheck

import (
	"fmt"
	"iracema/ast"
	"iracema/token"
	"strings"
)

type ErrList []error

func (el ErrList) Error() string {
	b := new(strings.Builder)

	return b.String()
}

func (list ErrList) Clean() ErrList {
	var newList ErrList

	for _, err := range list {
		// remove INVALID errors
		if strings.Index(err.Error(), "INVALID") > 0 {
			continue
		}

		newList = append(newList, err)
	}

	return newList
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

type signature struct {
	name   string
	ret    Type
	params []Type
}

type typechecker struct {
	*gamma

	this        Type
	sig         *signature
	file        *ast.File
	errs        ErrList
	laterChecks []func()
}

func (tc *typechecker) check() ErrList {
	for _, objDecl := range tc.file.ObjectList {
		tc.defineObject(objDecl)
	}

	// All object type are defined by now
	// set fieldset and methodset
	for _, laterCheck := range tc.laterChecks {
		laterCheck()
	}

	for _, varDecl := range tc.file.VarList {
		tc.checkLetDecl(varDecl)
	}

	for _, funDecl := range tc.file.FunList {
		tc.checkFunDecl(funDecl)
	}

	for _, objDecl := range tc.file.ObjectList {
		tc.checkObjectDecl(objDecl)
	}

	return tc.errs
}

func (tc *typechecker) checkStatement(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		switch s := stmt.(type) {

		case *ast.VarDecl:
			tc.checkLetDecl(s)

		case *ast.AssignStmt:
			tc.checkAssignStmt(s)

		case *ast.ReturnStmt:
			tc.checkReturnStmt(s)

		case *ast.ExprStmt:
			tc.checkExpr(s.Expr)

		case *ast.IfStmt:
			//TODO: check this

		case *ast.ForStmt:
			// TODO: check this as well

		case *ast.WhileStmt:
			// TODO: check this as well

		case *ast.SwitchStmt:
			// TODO: check this as well

		default:
			fmt.Printf("%+v\n", s)
			panic("unreacheble")
		}
	}
}

func (tc *typechecker) checkFunDecl(fun *ast.FunDecl) {
	defer func(s *gamma) {
		tc.sig = nil
		tc.gamma = s
	}(tc.gamma)

	tc.sig = tc.this.LookupMethod(fun.Name.Value)
	tc.gamma = &gamma{
		parent: tc.gamma,
		env:    make(map[string]Type),
	}

	for i, param := range fun.Parameters {
		tc.env[param.Name.Value] = tc.sig.params[i]
	}

	tc.checkBodyStmt(fun.Body)
	if tc.sig.ret != NONE && !tc.hasReturn(fun.Body) {
		tc.errorf(fun.Body.RightBrace, "missing return")
	}
}

func (tc *typechecker) hasReturn(stmt ast.Stmt) bool {
	switch nd := stmt.(type) {
	case *ast.VarDecl, *ast.AssignStmt, *ast.ExprStmt:
		// ignore

	case *ast.BlockStmt:
		size := len(nd.Stmts) - 1
		if size < 0 {
			return false
		}

		return tc.hasReturn(nd.Stmts[size])

	case *ast.ReturnStmt:
		return true

	case *ast.ForStmt, *ast.WhileStmt:
		return false

	case *ast.IfStmt:
		return nd.Else != nil && tc.hasReturn(nd.Then) && tc.hasReturn(nd.Else)

	case *ast.SwitchStmt:
		if nd.Default == nil {
			return false
		}

		for _, c := range nd.Cases {
			if !tc.hasReturn(c.Body) {
				return false
			}
		}

		return tc.hasReturn(nd.Default.Body)

	default:
		panic("unreachable")
	}

	return false
}

func (tc *typechecker) checkObjectDecl(obj *ast.ObjectDecl) {
	defer func(this Type) { tc.this = this }(tc.this)

	tc.this = tc.checkExpr(obj.Name)
	for _, fun := range obj.FunList {
		tc.checkFunDecl(fun)
	}
}

func (tc *typechecker) checkBodyStmt(body *ast.BlockStmt) {
	tc.checkStatement(body.Stmts)
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
				tc.errorf(let.Token, "cannot use '%s' as '%s' value in declaration", value.Name(), typ.Name())
			}
		}

		return
	}

	tc.env[name] = tc.checkExpr(let.Value)
}

func (tc *typechecker) checkAssignStmt(assign *ast.AssignStmt) {
	if len(assign.Left) != len(assign.Right) {
		tc.errorf(assign.Token, "assignment mismatch: %d variables but %d values", len(assign.Left), len(assign.Right))
		return
	}

	size := len(assign.Left)
	for i := 0; i < size; i++ {
		rhs := assign.Right[i]
		lhs := assign.Left[i]

		lhsType := tc.checkExpr(lhs)
		rhsType := tc.checkExpr(rhs)

		if !rhsType.Is(lhsType) {
			tc.errorf(assign.Token, "cannot use '%s' as '%s' value in assignment", rhsType.Name(), lhsType.Name())
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

		tc.errorf(node.Token, "undefined: %s", node.Value)
		return INVALID

	case *ast.CallExpr:
		return tc.checkCallExpr(node)

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

func (tc *typechecker) checkCallExpr(call *ast.CallExpr) Type {
	var recvType Type
	if call.Receiver != nil {
		recvType = tc.checkExpr(call.Receiver)
	} else {
		recvType = tc.this
	}

	sig := recvType.LookupMethod(call.Method.Value)
	if sig == nil {
		// TODO: fix the token position

		tc.errorf(nil, "object '%s' has no method '%s'", recvType, call.Method.Value)
		return INVALID
	}

	argc := len(call.Arguments)
	if len(sig.params) != argc {
		// TODO: fix the token position
		tc.errorf(nil, "wrong number of arguments (given %d, expected %d)", argc, len(sig.params))
		return sig.ret
	}

	var argTypes = make([]Type, argc)
	for i, arg := range call.Arguments {
		argTypes[i] = tc.checkExpr(arg)
	}

	tc.checkArguments(sig, argTypes)
	return sig.ret
}

func (tc *typechecker) checkArguments(sig *signature, argTypes []Type) {
	for i, paramType := range sig.params {
		argType := argTypes[i]
		if !argType.Is(paramType) {
			// TODO: fix this
			tc.errorf(nil, "cannot use '%s' as '%s' in argument to %s", argType, paramType, sig.name)
		}
	}
}

func (tc *typechecker) errorf(tok *token.Token, format string, args ...any) {
	if tok != nil {
		format = "[Lin: %d Col: %d] " + format
		args = append([]any{tok.Column()}, args...)
		args = append([]any{tok.Line()}, args...)
	}
	err := fmt.Errorf(format, args...)
	tc.errs = append(tc.errs, err)
}

func (tc *typechecker) checkReturnStmt(ret *ast.ReturnStmt) {
	if tc.sig == nil {
		tc.errorf(ret.Token, "return outside function")
		return
	}

	if tc.sig.ret == NONE && ret.Value != nil {
		tc.errorf(ret.Token, "unexpected return value")
		return
	}

	valueType := tc.checkExpr(ret.Value)

	if !valueType.Is(tc.sig.ret) {
		tc.errorf(ret.Token, "cannot use '%s' as '%s' value in return statement", valueType, tc.sig.ret)
	}
}

func Check(file *ast.File) ErrList {
	tc := new(typechecker)
	tc.file = file
	tc.this = SCRIPT
	tc.gamma = &gamma{
		env: make(map[string]Type),
	}

	return tc.check()
}

func (tc *typechecker) defineObject(decl *ast.ObjectDecl) Type {
	objType := newObject(decl.Name.Value)

	tc.laterChecks = append(tc.laterChecks, func() {
		if decl.Parent != nil {
			objType.parent = tc.checkExpr(decl.Parent)
		}

		for _, f := range decl.FieldList {
			fieldType := tc.checkExpr(f.Type)
			if field := objType.addField(f.Name.Value, fieldType); field != nil {
				tc.errorf(f.Token, "field %s is already defined in object %s", f.Name.Value, objType)
			}
		}

		for _, f := range decl.FunList {
			name := f.Name.Value
			paramTypes := tc.paramTypes(f.Parameters)
			var retType Type = NONE

			if name == "init" {
				if f.Return != nil {
					tc.errorf(f.Return.Token, "init fun can not have return value")
				}

				name = "new"
				retType = objType
			}

			if f.Return != nil {
				retType = tc.lookupType(f.Return.Value)
			}

			sig := &signature{
				name:   name,
				params: paramTypes,
				ret:    retType,
			}

			objType.addMethod(sig)
		}

		objType.complete()
	})

	tc.env[objType.name] = objType
	return objType
}

func (tc *typechecker) paramTypes(fields []*ast.Field) []Type {
	list := make([]Type, len(fields))

	for i, f := range fields {
		list[i] = tc.lookupType(f.Type.Value)
	}

	return list
}
