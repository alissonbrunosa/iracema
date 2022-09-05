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

	tc.sig = tc.this.Method(fun.Name.Value)
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
	case *ast.VarDecl, *ast.AssignStmt, *ast.ExprStmt: // ignore
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
			tc.assign(decl.Value, tc.checkExpr(decl.Value), typ, "in declaration")
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

		tc.assign(rhs, tc.checkExpr(rhs), tc.checkExpr(lhs), "in assignment")
	}
}

func (tc *typechecker) checkExpr(expr ast.Expr) Type {
	switch node := expr.(type) {
	case *ast.BasicLit:
		return tc.checkBasicLit(node)

	case *ast.Ident:
		return tc.checkIdent(node)

	case *ast.MethodCallExpr:
		return tc.checkMethodCallExpr(node)

	case *ast.FunctionCallExpr:
		return tc.checkFunctionCallExpr(node)

	case *ast.SuperExpr:
		return tc.checkSuperExpr(node)

	case *ast.UnaryExpr:
		return tc.checkUnary(node)

	case *ast.BinaryExpr:
		return tc.checkBinary(node)

	case *ast.NewExpr:
		return tc.checkNewExpr(node)

	case *ast.MemberSelector:
		return tc.checkMemberSelector(node)

	default:
		fmt.Printf("%+v\n", node)
		panic("unreacheble")
	}
}

func (tc *typechecker) checkBasicLit(node *ast.BasicLit) Type {
	if litType := LIT_TYPES[node.T.String()]; litType != nil {
		return litType
	}

	if node.T == token.This {
		return tc.this
	}

	return INVALID
}

func (tc *typechecker) checkIdent(node *ast.Ident) (t Type) {
	if t = tc.lookupType(node.Value); t == nil {
		tc.errorf(node, "undefined: %s", node.Value)
		return INVALID
	}

	return
}

func (tc *typechecker) lookupType(name string) Type {
	if t, ok := LIT_TYPES[name]; ok {
		return t
	}

	if t := tc.lookup(name); t != nil {
		return t
	}

	return nil
}

func (tc *typechecker) checkMethodCallExpr(mCall *ast.MethodCallExpr) Type {
	path := mCall.Selector
	baseType := tc.checkExpr(path.Base)

	sig := baseType.Method(path.Member.Value)
	if sig == nil {
		tc.errorf(path.Member, "object '%s' has no method '%s'", baseType, path.Member.Value)
		return INVALID
	}

	argc := len(mCall.Arguments)
	if len(sig.params) != argc {
		tc.errorf(mCall, "wrong number of arguments (given %d, expected %d)", argc, len(sig.params))
		return sig.ret
	}

	tc.checkArguments(sig, mCall.Arguments...)
	return sig.ret
}

func (tc *typechecker) checkFunctionCallExpr(fCall *ast.FunctionCallExpr) Type {
	baseType := tc.this

	sig := baseType.Method(fCall.Name.Value)
	if sig == nil {
		tc.errorf(fCall.Name, "object '%s' has no method '%s'", baseType, fCall.Name.Value)
		return INVALID
	}

	argc := len(fCall.Arguments)
	if len(sig.params) != argc {
		tc.errorf(fCall, "wrong number of arguments (given %d, expected %d)", argc, len(sig.params))
		return sig.ret
	}

	tc.checkArguments(sig, fCall.Arguments...)
	return sig.ret
}

func (tc *typechecker) checkSuperExpr(node *ast.SuperExpr) Type {
	if tc.sig == nil {
		tc.errorf(node, "super called outside of method")
		return INVALID
	}

	parent := tc.this.Parent()
	if parent == nil {
		tc.errorf(node, "no superclass of '%s' has method '%s'", tc.this, tc.sig.name)
		return INVALID
	}

	sig := parent.Method(tc.sig.name)
	if sig == nil {
		tc.errorf(node, "no superclass of '%s' has method '%s'", tc.this, tc.sig.name)
		return INVALID
	}

	argc := len(node.Arguments)
	if len(sig.params) != argc {
		tc.errorf(node, "wrong number of arguments (given %d, expected %d)", argc, len(sig.params))
		return sig.ret
	}

	tc.checkArguments(sig, node.Arguments...)
	return sig.ret
}

func (tc *typechecker) checkUnary(node *ast.UnaryExpr) Type {
	operandType := tc.checkExpr(node.Expr)
	if t := unary(node.Operator.Type, operandType); t != nil {
		return t
	}

	fn := operandType.Method(node.Operator.String())
	if fn == nil {
		tc.errorf(node, "object '%s' do not implement '%s' unary operator", operandType, node.Operator)
		return INVALID
	}

	return INVALID
}

func (tc *typechecker) checkBinary(node *ast.BinaryExpr) Type {
	lhsType := tc.checkExpr(node.Left)
	rhsType := tc.checkExpr(node.Right)

	if typ := binary(node.Operator.Type, lhsType, rhsType); typ != nil {
		return typ
	}

	fn := lhsType.Method(node.Operator.String())
	if fn == nil {
		tc.errorf(node, "object '%s' do not implement '%s' operator", lhsType, node.Operator)
		return INVALID
	}

	tc.checkArguments(fn, node.Right)
	return fn.ret
}

func (tc *typechecker) checkNewExpr(node *ast.NewExpr) Type {
	objType := tc.lookupType(node.Type.Value)
	initFun := objType.Method("init")
	tc.checkArguments(initFun, node.Arguments...)
	return objType
}

func (tc *typechecker) checkMemberSelector(node *ast.MemberSelector) (field Type) {
	objType := tc.checkExpr(node.Base)

	if field = objType.Field(node.Member.Value); field == nil {
		tc.errorf(node.Member, "'%s' object has no field '%s'", objType, node.Member.Value)
		return INVALID
	}

	return field
}

func (tc *typechecker) checkArguments(sig *signature, args ...ast.Expr) {
	context := fmt.Sprintf("in argument to %s", sig.name)

	for i, paramType := range sig.params {
		tc.assign(args[i], tc.checkExpr(args[i]), paramType, context)
	}
}

func (tc *typechecker) errorf(node ast.Node, format string, args ...any) {
	if node != nil {
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
		tc.assign(ret.Value, tc.checkExpr(ret.Value), tc.sig.ret, "in return statement")
	}
}

func (tc *typechecker) checkIfStmt(ifStmt *ast.IfStmt) {
	tc.assign(ifStmt.Cond, tc.checkExpr(ifStmt.Cond), BOOL, "in if statement")

	tc.checkBlockStmt(ifStmt.Then)
	if ifStmt.Else != nil {
		tc.checkStmt(ifStmt.Else)
	}
}

func (tc *typechecker) checkWhileStmt(while *ast.WhileStmt) {
	defer func(f byte) { tc.flag = f }(tc.flag)

	tc.flag |= loop
	tc.assign(while.Cond, tc.checkExpr(while.Cond), BOOL, "in while statement")
	tc.checkBlockStmt(while.Body)
}

func (tc *typechecker) checkStmtSwitch(stmt *ast.SwitchStmt) {
	keyType := tc.checkExpr(stmt.Key)

	eqFun := keyType.Method("==")
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

			if name == "init" && f.Return != nil {
				tc.errorf(f.Return, "function init can not have return value")
			} else if f.Return != nil {
				retType = tc.lookupType(f.Return.Value)
			}

			sig := &signature{name: name, params: paramTypes, ret: retType}
			if m := objType.addMethod(sig); m != nil {
				tc.errorf(f, "function %s is already defined in object %s", sig.name, objType)
			}
		}

		objType.complete()
	})

	if !tc.Insert(objType.name, objType) {
		tc.errorf(decl.Name, "object %s is already declared", objType.Name())
	}

	return objType
}

func (tc *typechecker) paramTypes(fields []*ast.Field) []Type {
	list := make([]Type, len(fields))

	for i, f := range fields {
		list[i] = tc.lookupType(f.Type.Value)
	}

	return list
}

func (tc *typechecker) assign(node ast.Node, value, typ Type, context string) {
	if value.Is(typ) {
		return
	}

	tc.errorf(node, "expected '%s', found '%s' %s", typ, value, context)
}

func isNumber(t Type) bool {
	return t == INT || t == FLOAT
}

func isComparable(t Type) bool {
	return isNumber(t) || t == STRING
}

func unary(operator token.Type, operandType Type) Type {
	if isNumber(operandType) && (operator == token.Plus || operator == token.Minus) {
		return operandType
	}

	if operandType == BOOL && operator == token.Not {
		return operandType
	}

	return nil
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
