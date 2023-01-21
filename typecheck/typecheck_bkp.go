package typecheck

// type signature struct {
// 	name   string
// 	ret    Type
// 	params []Type
// }
//

//

//
// type earlyReturn struct{}
//

//
// func (tc *typechecker) checkExpr(expr ast.Expr) Type {
// 	switch node := expr.(type) {
// 	case *ast.BasicLit:
// 		if k := LIT_TYPES[node.Token.Type.String()]; k != nil {
// 			return k
// 		}
//
// 		if node.Token.Type == token.This {
// 			return tc.this
// 		}
//
// 		fmt.Println(">>>>", node.Token.Type.String(), LIT_TYPES)
// 		tc.errorf("Invalid literal")
// 		return INVALID
// 	case *ast.Ident:
// 		return tc.lookupType(node.Value)
//
// 	case *ast.UnaryExpr:
// 		// TODO: check if type implements unary operator
// 		return tc.checkExpr(node.Expr)
//
// 	case *ast.BinaryExpr:
// 		lhsType := tc.checkExpr(node.Left)
// 		rhsType := tc.checkExpr(node.Right)
//
// 		if kind, ok := tc.numberType(lhsType, rhsType); ok {
// 			return kind
// 		}
//
// 		var sig *signature
// 		if sig = lhsType.LookupFun(node.Operator.String()); sig == nil {
// 			tc.errorf("object %s does not implement operator '%s'", lhsType.Name(), node.Operator)
// 			return INVALID
// 		}
//
// 		tc.checkArguments(sig, []ast.Expr{node.Right})
// 		return sig.ret
//
// 	case *ast.CallExpr:
// 		return tc.checkCallExpr(node)
//
// 	case *ast.GroupExpr:
// 		return tc.checkExpr(node.Expr)
//
// 	case *ast.FieldSel:
// 		t := tc.this
// 		if ft := t.LookupField(node.Name.Value); ft != nil {
// 			return ft
// 		}
//
// 		tc.errorf("object %s has no field %s", t.Name(), node.Name.Value)
// 		return INVALID
//
// 	case *ast.HashLit:
// 		return tc.checkHashLit(node)
//
// 	case *ast.ArrayLit:
// 		return tc.checkArrayLit(node)
//
// 	case *ast.IndexExpr:
// 		exprType := tc.checkExpr(node.Expr)
// 		fun := exprType.LookupFun("get")
// 		if fun == nil {
// 			tc.errorf("type %s does not implement get method", exprType.Name())
// 			return exprType
// 		}
//
// 		tc.checkArguments(fun, []ast.Expr{node.Index})
// 		return fun.ret
//
// 	default:
// 		fmt.Printf("%+v\n", node)
// 		panic("unreacheble")
// 	}
// }
//
// func (tc *typechecker) checkHashLit(node *ast.HashLit) Type {
// 	if len(node.Entries) == 0 {
// 		return untypedHash()
// 	}
//
// 	var keyType Type
// 	var valueType Type
// 	for i, entry := range node.Entries {
// 		if i == 0 {
// 			if keyType = tc.checkExpr(entry.Key); keyType == nil {
// 				return nil
// 			}
// 			if valueType = tc.checkExpr(entry.Value); valueType == nil {
// 				return nil
// 			}
//
// 			continue
// 		}
//
// 		if kt := tc.checkExpr(entry.Key); kt == nil || kt != keyType {
// 			tc.errorf("mixed key types in Hash literal")
// 			return nil
// 		}
//
// 		if vt := tc.checkExpr(entry.Value); vt == nil || vt != valueType {
// 			tc.errorf("mixed value types in Hash literal")
// 			return nil
// 		}
// 	}
//
// 	return newHash(keyType, valueType)
// }
//
// func (tc *typechecker) checkArrayLit(node *ast.ArrayLit) Type {
// 	if len(node.Elements) == 0 {
// 		return untypedArray()
// 	}
//
// 	var valueType Type
// 	for i, element := range node.Elements {
// 		if i == 0 {
// 			if valueType = tc.checkExpr(element); valueType == INVALID {
// 				return INVALID
// 			}
//
// 			continue
// 		}
//
// 		if kt := tc.checkExpr(element); kt != valueType {
// 			tc.errorf("mixed values types in Array literal")
// 			return INVALID
// 		}
// 	}
//
// 	return newArray(valueType)
// }
//
// func (tc *typechecker) checkObjectDecl(obj *ast.ObjectDecl) {
// 	defer func() { tc.this = nil }()
//
// 	tc.this = tc.checkExpr(obj.Name)
// 	for _, fun := range obj.FunList {
// 		tc.checkFunDecl(fun)
// 	}
// }
//
// func (tc *typechecker) checkCallExpr(call *ast.CallExpr) Type {
// 	var recvType Type
// 	if call.Receiver != nil {
// 		recvType = tc.checkExpr(call.Receiver)
// 	} else {
// 		recvType = tc.this
// 	}
//
// 	var sig *signature
// 	if sig = recvType.LookupFun(call.Method.Value); sig == nil {
// 		tc.errorf("object %s has no method %s", recvType.Name(), call.Method.Value)
// 		return INVALID
// 	}
//
// 	tc.checkArguments(sig, call.Arguments)
//
// 	return sig.ret
// }
//
// func (tc *typechecker) checkFunDecl(fun *ast.FunDecl) {
// 	defer func(s *gamma) {
// 		tc.fun = nil
// 		tc.gamma = s
// 	}(tc.gamma)
//
// 	tc.fun = tc.this.LookupFun(fun.Name.Value)
// 	tc.gamma = &gamma{
// 		parent: tc.gamma,
// 		env:    make(map[string]Type),
// 	}
//
// 	for i, param := range fun.Parameters {
// 		tc.env[param.Name.Value] = tc.fun.params[i]
// 	}
//
// 	tc.checkBodyStmt(fun.Body)
// }
//
// func (tc *typechecker) checkAssignStmt(assign *ast.AssignStmt) {
// 	if len(assign.Left) != len(assign.Right) {
// 		tc.errorf("assignment mismatch: %d variables but %d values", len(assign.Left), len(assign.Right))
// 		return
// 	}
//
// 	size := len(assign.Left)
// 	for i := 0; i < size; i++ {
// 		rhs := assign.Right[i]
//
// 		switch lhs := assign.Left[i].(type) {
// 		case *ast.IndexExpr:
// 			tc.assignToIdenxExpr(lhs, rhs)
// 		default:
// 			tc.checkAssignToKeyword(lhs)
//
// 			lhsType := tc.checkExpr(lhs)
// 			rhsType := tc.checkExpr(rhs)
//
// 			if !rhsType.Is(lhsType) {
// 				tc.errorf("cannot use %s as %s value in assignment", rhsType.Name(), lhsType.Name())
// 			}
// 		}
// 	}
// }
//
// // IndexExpr -> names[0]
// // receiver[argument] == receiver.get(0)
// func (tc *typechecker) assignToIdenxExpr(lhs *ast.IndexExpr, rhs ast.Expr) {
// 	recvType := tc.checkExpr(lhs.Expr)
// 	if recvType == nil {
// 		return
// 	}
//
// 	switch t := recvType.(type) {
// 	case *hash:
// 		if t.isUntyped() {
// 			key := tc.checkExpr(lhs.Index)
// 			value := tc.checkExpr(rhs)
// 			t.setType(key, value)
// 		}
//
// 		fn := t.LookupFun("insert")
// 		tc.checkArguments(fn, []ast.Expr{lhs.Index, rhs})
//
// 	case *array:
// 		if t.isUntyped() {
// 			value := tc.checkExpr(rhs)
// 			t.setType(value)
// 		}
//
// 		fn := t.LookupFun("insert")
// 		tc.checkArguments(fn, []ast.Expr{lhs.Index, rhs})
//
// 	default:
// 		if fn := t.LookupFun("insert"); fn != nil {
// 			tc.checkArguments(fn, []ast.Expr{lhs.Index, rhs})
// 		}
//
// 		tc.errorf("type %s does not implement insert method", t.Name())
// 	}
// }
//
// func (tc *typechecker) checkArguments(fn *signature, args []ast.Expr) {
// 	var argTypes []Type
// 	for _, arg := range args {
// 		argTypes = append(argTypes, tc.checkExpr(arg))
// 	}
//
// 	if len(fn.params) != len(argTypes) {
// 		var b = new(strings.Builder)
// 		fmt.Fprintf(b, "wrong number of arguments (given %d, expected %d)\n", len(args), len(fn.params))
// 		argumentDetails(b, fn, argTypes)
// 		tc.errorf(b.String())
// 		return
// 	}
//
// 	for i, argType := range argTypes {
// 		if !argType.Is(fn.params[i]) {
// 			var b = new(strings.Builder)
// 			fmt.Fprintf(b, "incompatible types in call for %s\n", fn.name)
// 			argumentDetails(b, fn, argTypes)
// 			tc.errorf(b.String())
// 			return
// 		}
// 	}
// }
//
// func argumentDetails(b *strings.Builder, sig *signature, args []Type) {
// 	fmt.Fprintf(b, "\t    give: %s(", sig.name)
// 	if len(args) > 0 {
// 		fmt.Fprintf(b, "%s", args[0].Name())
// 		for i := 1; i < len(args); i++ {
// 			fmt.Fprint(b, ", ")
// 			fmt.Fprintf(b, "%s", args[i].Name())
// 		}
// 	}
// 	fmt.Fprintln(b, ")")
//
// 	fmt.Fprintf(b, "\texpected: %s(", sig.name)
// 	if len(sig.params) > 0 {
// 		fmt.Fprintf(b, "%s", sig.params[0].Name())
// 		for i := 1; i < len(sig.params); i++ {
// 			fmt.Fprint(b, ", ")
// 			fmt.Fprintf(b, "%s", sig.params[i].Name())
// 		}
// 	}
// 	fmt.Fprintln(b, ")")
// }
//
// func (tc *typechecker) checkReturnStmt(ret *ast.ReturnStmt) {
// 	if ret.Value != nil && tc.fun.ret == NONE {
// 		tc.errorf("unexpected return value, method %s has default return", tc.fun.name)
// 		return
// 	}
//
// 	valueType := tc.checkExpr(ret.Value)
// 	if !valueType.Is(tc.fun.ret) {
// 		tc.errorf("cannot use %s as %s value in return statement", valueType.Name(), tc.fun.ret.Name())
// 	}
// }
//
// func (tc *typechecker) checkAssignToKeyword(expr ast.Expr) {
// 	switch node := expr.(type) {
// 	case *ast.BasicLit:
// 		switch node.Token.Type {
// 		case token.This:
// 			tc.errorf("Can not assign value to this")
//
// 		case token.None:
// 			tc.errorf("Can not assign value to none")
//
// 		case token.Bool:
// 			tc.errorf("Can not assign value to bool(%s)", node.Value)
// 		}
// 	default:
// 		return
// 	}
// }
//
// func (tc *typechecker) checkBodyStmt(body *ast.BlockStmt) {
// 	tc.checkStatement(body.Stmts)
// }
//
// // func (tc *typechecker) checkExpr(expr ast.Expr) Type {
// // 	switch node := expr.(type) {
// // 	case *ast.Ident:
// // 		return tc.env[node.Value]
// //
// // 	case *ast.BasicLit:
// // 		if k := litTypes[node.Token.Type.String()]; k != nil {
// // 			return k
// // 		}
// //
// // 		if node.Token.Type == token.This {
// // 			return tc.env["this"]
// // 		}
// //
// // 		fmt.Println("Invalid literal")
// // 		return nil
// //
// // 	case *ast.UnaryExpr:
// // 		return tc.checkExpr(node.Expr)
// //
// // 	case *ast.BinaryExpr:
// // 		if kind, ok := tc.numberKind(node); ok {
// // 			return kind
// // 		}
// //
// // 		panic("OH BOY!")
// //
// // 	case *ast.CallExpr:
// // 		tc.checkExpr(node)
// // 		var recvKind Type
// // 		if node.Receiver != nil {
// // 			recvKind = tc.checkExpr(node.Receiver)
// // 		} else {
// // 			recvKind = tc.env["this"]
// // 		}
// //
// // 		f := recvKind.LookupFun(node.Method.Value)
// // 		return f.ret
// //
// // 	case *ast.FieldSel:
// // 		t := tc.lookupType("this")
// // 		if ft, ok := t.fields[node.Name.Value]; ok {
// // 			return ft
// // 		}
// //
// // 		fmt.Println("object %s has no field %s", t.Name(), node.Name.Value)
// // 		return nil
// //
// // 	default:
// // 		panic(earlyReturn{})
// // 	}
// // }
//
// func (tc *typechecker) numberType(lhsType, rhsType Type) (Type, bool) {
//
// 	if lhsType == FLOAT && (rhsType == INT || rhsType == FLOAT) {
// 		return FLOAT, true
// 	}
//
// 	if lhsType == INT && rhsType == FLOAT {
// 		return FLOAT, true
// 	}
//
// 	if lhsType == INT && rhsType == INT {
// 		return INT, true
// 	}
//
// 	return nil, false
// }
//
// func (tc *typechecker) lookupType(name string) Type {
// 	if k, ok := LIT_TYPES[name]; ok {
// 		return k
// 	}
//
// 	if typ := tc.lookup(name); typ != nil {
// 		return typ
// 	}
//
// 	tc.errorf("Type not defined: %s", name)
// 	return INVALID
// }
//

//
// func (tc *typechecker) typeParams(fields []*ast.Field) []Type {
// 	list := make([]Type, len(fields))
//
// 	for i, f := range fields {
// 		list[i] = tc.lookupType(f.Type.Value)
// 	}
//
// 	return list
// }
//
// func (tc *typechecker) errorf(format string, args ...any) {
// 	err := fmt.Errorf(format, args...)
// 	tc.errs = append(tc.errs, err)
// }
