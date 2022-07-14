package compile

import (
	"fmt"
	"iracema/ast"
	"iracema/bytecode"
	"iracema/lang"
	"iracema/token"
	"strconv"
)

type local struct {
	name        string
	index       byte
	depth       int
	initialized bool
}

func (l *local) String() string {
	return fmt.Sprintf("%s@%d", l.name, l.index)
}

type branchType byte

const (
	basic branchType = 1 << iota
	loop
	catch
)

type jumpLabel struct {
	index  int
	cond   bool
	target *branch
}

type branch struct {
	kind      branchType
	startAt   byte
	hasReturn bool
	next      *branch
}

type instr struct {
	opcode  bytecode.Opcode
	operand byte
}

type fragment struct {
	name        string
	argc        byte
	consts      []lang.IrObject
	locals      []*local
	instrs      []*instr
	jumpLabels  []*jumpLabel
	catchOffset int

	curBranch *branch
	previous  *fragment
}

type compiler struct {
	*fragment
	fragments []*fragment
}

func New() *compiler {
	c := new(compiler)
	c.init()

	return c
}

func (c *compiler) init() {
	f := &fragment{name: "main"}
	startBranch := &branch{kind: basic}
	c.fragments = append([]*fragment{}, f)
	c.fragment = f
	c.useBranch(startBranch)
}

func (c *compiler) Compile(file *ast.File) (*lang.Method, error) {
	c.compileStmt(file)
	return c.assemble(), nil
}

func (c *compiler) assemble() *lang.Method {
	var bytecode []uint16
	c.patchJumps()

	for _, instr := range c.instrs {
		code := uint16(instr.opcode)<<8 | uint16(instr.operand)
		bytecode = append(bytecode, code)
	}

	return lang.NewIrMethod(
		c.name,
		c.argc,
		bytecode,
		byte(len(c.locals)),
		c.consts,
		c.catchOffset,
	)
}

func (c *compiler) patchJumps() {
	if len(c.jumpLabels) == 0 {
		return
	}

	for _, label := range c.jumpLabels {
		if label.cond {
			c.instrs[label.index] = &instr{opcode: bytecode.Jump, operand: label.target.startAt}
			continue
		}

		c.instrs[label.index] = &instr{opcode: bytecode.JumpIfFalse, operand: label.target.startAt}
	}
}

func (c *compiler) compileStmt(stmt ast.Stmt) {
	switch node := stmt.(type) {
	case *ast.File:
		if len(node.Stmts) == 0 {
			c.add(bytecode.PushNone, 0)
			c.add(bytecode.Return, 0)
			return
		}

		for _, stmt := range node.Stmts {
			c.compileStmt(stmt)
		}

		if !c.curBranch.hasReturn {
			c.add(bytecode.PushNone, 0)
			c.add(bytecode.Return, 0)
		}

	case *ast.ExprStmt:
		c.compileExpr(node.Expr, false)

	case *ast.AssignStmt:
		for i, value := range node.Right {
			switch lhs := node.Left[i].(type) {
			case *ast.Ident:
				if lhs.IsAttr() {
					c.compileExpr(value, true)
					c.add(bytecode.SetAttr, c.addConstant(lhs.Value))
					continue
				}

				local := c.defineLocal(lhs.Value, false)
				c.compileExpr(value, true)
				local.initialized = true
				c.add(bytecode.SetLocal, local.index)

			case *ast.IndexExpr:
				c.compileExpr(lhs.Expr, true)
				c.compileExpr(lhs.Index, true)
				c.compileExpr(value, true)
				ci := lang.NewCallInfo("insert", 2)
				c.add(bytecode.CallMethod, c.addConstant(ci))
				c.add(bytecode.Pop, 0)
			}
		}

	case *ast.WhileStmt:
		c.compileWhileStmt(node)

	case *ast.ForStmt:
		c.compileForStmt(node)

	case *ast.SwitchStmt:
		c.compileSwitchStmt(node)

	case *ast.ObjectDecl:
		object := c.compileObjectDecl(node)

		if node.Parent != nil {
			c.add(bytecode.GetConstant, c.addConstant(node.Parent.Value))
		} else {
			c.add(bytecode.PushNone, 0)
		}

		c.add(bytecode.DefineObject, c.addConstant(object))
		c.add(bytecode.Pop, 0)

	case *ast.FunDecl:
		fun := c.compileFunDecl(node)
		c.add(bytecode.DefineFunction, c.addConstant(fun))

	case *ast.BlockStmt:
		c.compileBlock(node, false)

	case *ast.ReturnStmt:
		c.compileExpr(node.Expr, true)
		c.add(bytecode.Return, 0)
		c.curBranch.hasReturn = true

	case *ast.StopStmt:
		if c.curBranch.kind == loop {
			c.jumpToBlock(c.curBranch.next, true)
			return
		}
		panic("STOP OUTSIDE LOOP")

	case *ast.NextStmt:
		if c.curBranch.kind == loop {
			c.jumpToBlock(c.curBranch, true)
			return
		}
		panic("NEXT OUTSIDE LOOP")

	case *ast.IfStmt:
		c.compileIfStmt(node)

	default:
		return
	}
}

func (c *compiler) compileExpr(expr ast.Expr, isEvaluated bool) {
	switch node := expr.(type) {
	case *ast.Ident:
		if node.IsAttr() {
			c.add(bytecode.GetAttr, c.addConstant(node.Value))
			return
		}

		if node.IsConstant() {
			c.add(bytecode.GetConstant, c.addConstant(node.Value))
			return
		}

		if local := c.resolve(node.Value); local != nil {
			if local.initialized == false {
				panic("underfined " + local.name)
			}

			c.add(bytecode.GetLocal, local.index)
			return
		}

		c.add(bytecode.PushSelf, 0)
		ci := lang.NewCallInfo(node.Value, 0)
		c.add(bytecode.CallMethod, c.addConstant(ci))

	case *ast.BasicLit:
		if isEvaluated {
			c.compileLiteral(node)
		}

	case *ast.UnaryExpr:
		c.compileExpr(node.Expr, true)
		c.add(bytecode.Not, 0)

	case *ast.BinaryExpr:
		c.compileExpr(node.Left, true)
		c.compileExpr(node.Right, true)
		operator := binary(node.Operator.String())
		c.add(bytecode.Binary, operator)

		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	case *ast.ArrayLit:
		c.compileArrayLireral(node)

	case *ast.GroupExpr:
		c.compileExpr(node.Expr, isEvaluated)

	case *ast.BlockExpr:
		panic("TODO: ainda falta escrever a solucao para clojures")

	case *ast.IndexExpr:
		c.compileExpr(node.Expr, true)
		c.compileExpr(node.Index, true)
		ci := lang.NewCallInfo("get", 1)
		c.add(bytecode.CallMethod, c.addConstant(ci))
		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	case *ast.CallExpr:
		if node.Receiver != nil {
			c.compileExpr(node.Receiver, true)
		} else {
			c.add(bytecode.PushSelf, 0)
		}

		for _, arg := range node.Arguments {
			c.compileExpr(arg, true)
		}

		ci := lang.NewCallInfo(node.Method.Value, byte(len(node.Arguments)))
		c.add(bytecode.CallMethod, c.addConstant(ci))
		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	default:
		return
	}
}

func binary(operator string) byte {
	switch operator {
	case "+":
		return ADD
	case "-":
		return SUB
	case "*":
		return MUL
	case "/":
		return DIV
	case "==":
		return EQ
	case "!=":
		return NE
	case ">":
		return GT
	case ">=":
		return GE
	case "<":
		return LT
	case "<=":
		return LE
	default:
		panic("not a binary operator")
	}
}

func (c *compiler) compileBlock(block *ast.BlockStmt, addReturn bool) {
	for _, stmt := range block.Stmts {
		c.compileStmt(stmt)
	}

	if addReturn && !c.curBranch.hasReturn {
		c.add(bytecode.PushNone, 0)
		c.add(bytecode.Return, 0)
		c.curBranch.hasReturn = true
	}
}

func (c *compiler) compileWhileStmt(node *ast.WhileStmt) {
	exit := &branch{kind: basic}
	loop := &branch{kind: loop, next: exit}

	c.useBranch(loop)
	c.compileExpr(node.Cond, true)
	c.jumpToBlock(exit, false)
	c.compileBlock(node.Body, false)
	c.jumpToBlock(loop, true)
	c.useBranch(exit)
}

func (c *compiler) compileForStmt(node *ast.ForStmt) {
	exit := &branch{kind: basic}
	setup := &branch{kind: basic}
	loop := &branch{kind: loop, next: exit}

	c.useBranch(setup)
	c.compileExpr(node.Iterable, true)
	c.add(bytecode.NewIterator, 0)

	c.useBranch(loop)
	c.add(bytecode.Iterate, 0)
	c.jumpToBlock(exit, false)

	local := c.defineLocal(node.Element.Value, true)
	c.add(bytecode.SetLocal, local.index)

	c.compileBlock(node.Body, false)
	c.jumpToBlock(loop, true)
	c.useBranch(exit)
}

func (c *compiler) compileSwitchStmt(node *ast.SwitchStmt) {
	endBranch := &branch{kind: basic}

	lenCases := len(node.Cases) - 1
	for i, caseClause := range node.Cases {
		nextCaseBranch := &branch{kind: basic}

		c.compileExpr(node.Key, true)
		c.compileExpr(caseClause.Value, true)

		callInfo := lang.NewCallInfo("==", 1)
		c.add(bytecode.CallMethod, c.addConstant(callInfo))

		c.jumpToBlock(nextCaseBranch, false)
		c.compileBlock(caseClause.Body, false)

		if i != lenCases || node.Default != nil {
			c.jumpToBlock(endBranch, true)
		}

		c.useBranch(nextCaseBranch)
	}

	if node.Default != nil {
		c.compileBlock(node.Default.Body, false)
	}

	c.useBranch(endBranch)
}

func (c *compiler) compileIfStmt(node *ast.IfStmt) {
	var elseBlock, endBlock *branch

	if node.Else == nil {
		endBlock = &branch{kind: basic}
		elseBlock = endBlock
	} else {
		elseBlock = &branch{kind: basic}
		endBlock = &branch{kind: basic}
	}

	c.compileExpr(node.Cond, true)
	c.jumpToBlock(elseBlock, false)
	c.compileBlock(node.Then, false)

	if node.Else != nil {
		c.jumpToBlock(endBlock, true)
		c.useBranch(elseBlock)
		c.compileStmt(node.Else)
	}

	c.useBranch(endBlock)
}

func (c *compiler) jumpToBlock(target *branch, cond bool) {
	label := &jumpLabel{index: len(c.instrs), target: target, cond: cond}
	c.jumpLabels = append(c.jumpLabels, label)
	c.instrs = append(c.instrs, nil) // take position
}

func (c *compiler) openScope(name string) {
	c.fragment = &fragment{
		name:      name,
		previous:  c.fragment,
		curBranch: &branch{kind: basic},
	}

	c.fragments = append(c.fragments, c.fragment)
}

func (c *compiler) closeScope() {
	c.fragment = c.fragment.previous
}

func (c *compiler) compileObjectDecl(node *ast.ObjectDecl) *lang.Method {
	c.openScope(node.Name.Value)
	defer c.closeScope()

	c.compileBlock(node.Body, true)
	return c.assemble()
}

func (c *compiler) compileFunDecl(node *ast.FunDecl) *lang.Method {
	c.openScope(node.Name.Value)
	defer c.closeScope()

	c.argc = byte(len(node.Parameters))
	for _, param := range node.Parameters {
		c.defineLocal(param.Value, true)
	}

	c.compileBlock(node.Body, true)

	end := &branch{kind: basic}
	if len(node.Catches) != 0 {
		c.catchOffset = len(c.instrs)

		for _, ch := range node.Catches {
			catch := end
			end = &branch{kind: basic}

			c.useBranch(catch)
			c.add(bytecode.MatchType, c.addConstant(ch.Type.Value))
			c.jumpToBlock(end, false)

			if ch.Ref != nil {
				local := c.defineLocal(ch.Ref.Value, true)
				c.add(bytecode.SetLocal, local.index)
			}

			c.compileBlock(ch.Body, true)
		}

		c.useBranch(end)
		c.add(bytecode.Throw, 0)
		c.curBranch.hasReturn = true
	}

	return c.assemble()
}

func isLiteral(node ast.Expr) bool {
	_, ok := node.(*ast.BasicLit)
	return ok
}

func (c *compiler) compileArrayLireral(node *ast.ArrayLit) {
	size := len(node.Elements)
	for _, el := range node.Elements {
		c.compileExpr(el, true)
	}

	c.add(bytecode.BuildArray, byte(size))
}

func (c *compiler) compileLiteral(lit *ast.BasicLit) {
	var val lang.IrObject

	switch lit.Type() {
	case token.String:
		val = lang.NewString(lit.Value)

	case token.Bool:
		value, err := strconv.ParseBool(lit.Value)
		if err != nil {
			panic("invalid boolean literal")
		}
		val = lang.Bool(value)

	case token.Int:
		value, err := strconv.Atoi(lit.Value)
		if err != nil {
			panic("invalid integer literal")
		}
		val = lang.Int(value)

	case token.Float:
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			panic("invalid float literal")
		}
		val = lang.Float(value)

	case token.Nil:
		c.add(bytecode.PushNone, 0)
		return

	default:
		fmt.Printf("DEFAULT %T, %q, %q\n", lit, lit.Value, lit.Token.Type)
		panic("invalid literal")
	}

	c.add(bytecode.Push, c.addConstant(val))
}

func literalValue(expr ast.Expr) lang.IrObject {
	lit := expr.(*ast.BasicLit)

	switch lit.Type() {
	case token.String:
		return lang.NewString(lit.Value)

	case token.Bool:
		value, err := strconv.ParseBool(lit.Value)
		if err != nil {
			panic("invalid boolean literal")
		}

		return lang.Bool(value)

	case token.Int:
		value, err := strconv.Atoi(lit.Value)
		if err != nil {
			panic("invalid integer literal")
		}
		return lang.Int(value)

	case token.Float:
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			panic("invalid float literal")
		}
		return lang.Float(value)

	case token.Nil:
		return lang.None

	default:
		fmt.Printf("DEFAULT %T, %q, %q\n", lit, lit.Value, lit.Token.Type)
		panic("invalid literal")
	}
}

func (c *compiler) defineLocal(name string, initialized bool) *local {
	if l := c.resolve(name); l != nil {
		return l
	}

	index := byte(len(c.locals))
	l := &local{name: name, index: index, initialized: initialized}
	c.locals = append(c.locals, l)
	return l
}

func (c *compiler) resolve(name string) *local {
	for _, l := range c.locals {
		if l.name == name {
			return l
		}
	}

	return nil
}

func (c *compiler) useBranch(b *branch) {
	b.startAt = byte(len(c.instrs))
	if c.curBranch == nil {
		c.curBranch = b
		return
	}

	c.curBranch.next = b
	c.curBranch = b
}

func (c *compiler) add(opcode bytecode.Opcode, operand byte) {
	c.instrs = append(c.instrs, &instr{opcode: opcode, operand: operand})
}

func (c *compiler) addConstant(arg interface{}) byte {
	switch val := arg.(type) {
	case int:
		c.consts = append(c.consts, lang.Int(val))
	case string:
		c.consts = append(c.consts, lang.NewString(val))
	case lang.IrObject:
		c.consts = append(c.consts, val)
	}

	return byte(len(c.consts) - 1)
}
