package typecheck

import (
	"fmt"
	"iracema/ast"
	"iracema/token"
	"strings"
)

const (
	loop = 1 << iota
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

func (g *gamma) Insert(name string, t Type) bool {
	if _, ok := g.env[name]; ok {
		return false
	}

	g.env[name] = t
	return true
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

	flag        byte
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
		tc.checkVarDecl(varDecl)
	}

	for _, funDecl := range tc.file.FunList {
		tc.checkFunDecl(funDecl)
	}

	for _, objDecl := range tc.file.ObjectList {
		tc.checkObjectDecl(objDecl)
	}

	return tc.errs
}

func (tc *typechecker) checkStmt(stmt ast.Stmt) {
	switch node := stmt.(type) {
	case *ast.VarDecl:
		tc.checkVarDecl(node)

	case *ast.AssignStmt:
		tc.checkAssignStmt(node)

	case *ast.ReturnStmt:
		tc.checkReturnStmt(node)

	case *ast.ExprStmt:
		tc.checkExpr(node.Expr)

	case *ast.IfStmt:
		tc.checkIfStmt(node)

	case *ast.ForStmt:
		// TODO: check this as well

	case *ast.WhileStmt:
		tc.checkWhileStmt(node)

	case *ast.SwitchStmt:
		tc.checkStmtSwitch(node)

	case *ast.BlockStmt:
		tc.checkBlockStmt(node)

	case *ast.StopStmt:
		if tc.flag&loop == 0 {
			tc.errorf(node, "stop statement outside loop")
		}

	case *ast.NextStmt:
		if tc.flag&loop == 0 {
			tc.errorf(node, "next statement outside loop")
		}

	default:
		fmt.Printf("%+v -> %T\n", node, node)
		panic("unreacheble")
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
		if !tc.Insert(param.Name.Value, tc.sig.params[i]) {
			tc.errorf(param.Name, "variable %s is already defined in function %s", param.Name.Value, tc.sig.name)
		}
	}

	tc.checkBlockStmt(fun.Body)
	if tc.sig.ret != NONE && !tc.hasReturn(fun.Body) {
		tc.errorf(fun, "missing return for function: %s", fun.Name.Value)
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

func (tc *typechecker) checkBlockStmt(body *ast.BlockStmt) {
	for _, stmt := range body.Stmts {
		tc.checkStmt(stmt)
	}
}

func (tc *typechecker) checkVarDecl(decl *ast.VarDecl) {
	name := decl.Name.Value

	if decl.Type == nil && decl.Value == nil {
		tc.env[name] = INVALID
		return
	}

	var typ Type
	if decl.Type != nil {
		typ = tc.checkExpr(decl.Type)
		if !tc.Insert(name, typ) {
			tc.errorf(decl.Name, "variable %s is already defined in function %s", name, tc.sig.name)
		}

		if decl.Value != nil {
			value := tc.checkExpr(decl.Value)
			if !value.Is(typ) {
				tc.errorf(decl.Value, "cannot use '%s' as '%s' value in declaration", value.Name(), typ.Name())
			}
		}

		return
	}

	if !tc.Insert(name, tc.checkExpr(decl.Value)) {
		tc.errorf(decl.Name, "variable %s is already defined in function %s", name, tc.sig.name)
	}
}

func (tc *typechecker) checkAssignStmt(assign *ast.AssignStmt) {
	if len(assign.Left) != len(assign.Right) {
		tc.errorf(assign, "assignment mismatch: %d variables but %d values", len(assign.Left), len(assign.Right))
		return
	}

	size := len(assign.Left)
	for i := 0; i < size; i++ {
		rhs := assign.Right[i]
		lhs := assign.Left[i]

		lhsType := tc.checkExpr(lhs)
		rhsType := tc.checkExpr(rhs)

		if !rhsType.Is(lhsType) {
			tc.errorf(rhs, "cannot use '%s' as '%s' value in assignment", rhsType.Name(), lhsType.Name())
		}
	}
}

func (tc *typechecker) checkExpr(expr ast.Expr) Type {
	switch node := expr.(type) {
	case *ast.BasicLit:
		if litType := LIT_TYPES[node.T.String()]; litType != nil {
			return litType
		}

		return INVALID
	case *ast.Ident:
		if t := tc.lookupType(node.Value); t != nil {
			return t
		}

		tc.errorf(node, "undefined: %s", node.Value)
		return INVALID

	case *ast.CallExpr:
		return tc.checkCallExpr(node)

	case *ast.SuperExpr:
		// TODO: check this
		return INVALID

	case *ast.BinaryExpr:
		return tc.checkBinary(node)

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
		tc.errorf(call.Method, "object '%s' has no method '%s'", recvType, call.Method.Value)
		return INVALID
	}

	argc := len(call.Arguments)
	if len(sig.params) != argc {
		// TODO: fix the token position
		tc.errorf(nil, "wrong number of arguments (given %d, expected %d)", argc, len(sig.params))
		return sig.ret
	}

	tc.checkArguments(sig, call.Arguments...)
	return sig.ret
}

func (tc *typechecker) checkBinary(node *ast.BinaryExpr) Type {
	lhsType := tc.checkExpr(node.Left)
	rhsType := tc.checkExpr(node.Right)

	if typ := binary(node.Operator.Type, lhsType, rhsType); typ != nil {
		return typ
	}

	fn := lhsType.LookupMethod(node.Operator.Type.String())
	if fn == nil {
		tc.errorf(node, "object '%s' do not implement '%s' operator", lhsType, node.Operator)
		return INVALID
	}

	tc.checkArguments(fn, node.Right)
	return fn.ret
}

func (tc *typechecker) checkArguments(sig *signature, args ...ast.Expr) {
	argTypes := make([]Type, len(args))
	for i, arg := range args {
		argTypes[i] = tc.checkExpr(arg)
	}

	for i, paramType := range sig.params {
		argType := argTypes[i]
		if !argType.Is(paramType) {
			tc.errorf(args[i], "cannot use '%s' as '%s' in argument to %s", argType, paramType, sig.name)
		}
	}
}

func (tc *typechecker) errorf(node ast.Node, format string, args ...any) {
	if node != nil {
		fmt.Printf("%T\n", node)
		pos := node.Position()
		format = "[Lin: %d Col: %d] " + format
		args = append([]any{pos.Column()}, args...)
		args = append([]any{pos.Line()}, args...)
	}
	err := fmt.Errorf(format, args...)
	tc.errs = append(tc.errs, err)
}

func (tc *typechecker) checkReturnStmt(ret *ast.ReturnStmt) {
	if tc.sig == nil {
		tc.errorf(ret, "return outside function")
		return
	}

	if tc.sig.ret == NONE && ret.Value != nil {
		tc.errorf(ret.Value, "unexpected return value")
		return
	}

	if ret.Value != nil {
		valueType := tc.checkExpr(ret.Value)

		if !valueType.Is(tc.sig.ret) {
			tc.errorf(ret.Value, "cannot use '%s' as '%s' value in return statement", valueType, tc.sig.ret)
		}
	}
}

func (tc *typechecker) checkIfStmt(ifStmt *ast.IfStmt) {
	condType := tc.checkExpr(ifStmt.Cond)
	if condType != BOOL {
		tc.errorf(ifStmt.Cond, "expected 'Bool', found '%s'", condType)
	}

	tc.checkBlockStmt(ifStmt.Then)
	if ifStmt.Else != nil {
		tc.checkStmt(ifStmt.Else)
	}
}

func (tc *typechecker) checkWhileStmt(while *ast.WhileStmt) {
	defer func(f byte) { tc.flag = f }(tc.flag)

	tc.flag |= loop
	condType := tc.checkExpr(while.Cond)
	if condType != BOOL {
		tc.errorf(while.Cond, "expected 'Bool', found '%s'", condType)
	}

	tc.checkBlockStmt(while.Body)
}

func (tc *typechecker) checkStmtSwitch(stmt *ast.SwitchStmt) {
	keyType := tc.checkExpr(stmt.Key)

	eqFun := keyType.LookupMethod("==")
	if eqFun == nil {
		tc.errorf(stmt.Key, "object '%s' do not implement '==' operator", keyType)
	}

	seen := make(map[string]ast.Expr)
	for _, c := range stmt.Cases {
		if lit, ok := c.Value.(*ast.BasicLit); ok {
			if prev, ok := seen[lit.Value]; ok {
				tc.errorf(lit, "duplicate case")
				tc.errorf(prev, "previous case")
				continue
			}

			seen[lit.Value] = c.Value
		}

		if eqFun != nil {
			fmt.Println(eqFun)
			tc.checkArguments(eqFun, c.Value)
		}
		tc.checkStmt(c.Body)
	}

	if stmt.Default != nil {
		tc.checkStmt(stmt.Default.Body)
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
	objType := newObject(decl.Name.Value, OBJECT)

	tc.laterChecks = append(tc.laterChecks, func() {
		if decl.Parent != nil {
			objType.parent = tc.checkExpr(decl.Parent)
		}

		for _, f := range decl.FieldList {
			fieldType := tc.checkExpr(f.Type)
			if field := objType.addField(f.Name.Value, fieldType); field != nil {
				tc.errorf(f, "field %s is already defined in object %s", f.Name.Value, objType)
			}
		}

		for _, f := range decl.FunList {
			name := f.Name.Value
			paramTypes := tc.paramTypes(f.Parameters)
			var retType Type = NONE

			if name == "init" {
				if f.Return != nil {
					tc.errorf(f.Return, "init fun can not have return value")
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

			if m := objType.addMethod(sig); m != nil {
				tc.errorf(f, "method %s is already defined in object %s", sig.name, objType)
			}
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

func isNumber(t Type) bool {
	return t == INT || t == FLOAT
}

func isComparable(t Type) bool {
	return isNumber(t) || t == STRING
}

func binary(operator token.Type, lhsType, rhsType Type) Type {
	if lhsType == INVALID || rhsType == INVALID {
		return INVALID
	}

	number := func(lhs, rhs Type) Type {
		if !isNumber(lhs) || !isNumber(rhs) {
			return nil
		}

		if lhs == FLOAT || rhs == FLOAT {
			return FLOAT
		}

		return INT
	}

	switch operator {
	case token.Plus:
		if t := number(lhsType, rhsType); t != nil {
			return t
		}

		if lhsType == STRING && rhsType == STRING {
			return STRING
		}

	case token.Minus, token.Slash, token.Star:
		return number(lhsType, rhsType)

	case token.Great, token.GreatEqual, token.Less, token.LessEqual:
		if isComparable(lhsType) && isComparable(rhsType) {
			return BOOL
		}

	case token.Equal, token.NotEqual:
		if isComparable(lhsType) && isComparable(rhsType) {
			return BOOL
		}

		if lhsType == BOOL && rhsType == BOOL {
			return BOOL
		}
	}

	return nil
}
