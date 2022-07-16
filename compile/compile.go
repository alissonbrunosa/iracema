package compile

import (
	"fmt"
	"iracema/ast"
	"iracema/bytecode"
	"iracema/lang"
	"iracema/token"
	"strconv"
)

const (
	FOR_LOOP   = 1
	WHILE_LOOP = 2
)

type controlflow struct {
	loop  int
	start *codeblock
	exit  *codeblock

	next *controlflow
}

type local struct {
	name        string
	index       byte
	initialized bool
}

func (l *local) String() string {
	return fmt.Sprintf("%s@%d", l.name, l.index)
}

type codeblock struct {
	startAt   byte
	hasReturn bool
}

type instr struct {
	opcode  bytecode.Opcode
	operand byte
	target  *codeblock
}

type fragment struct {
	name        string
	argc        byte
	consts      []lang.IrObject
	locals      []*local
	instrs      []*instr
	jumps       []int
	catchOffset int

	control  *controlflow
	block    *codeblock
	previous *fragment
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
	c.fragment = &fragment{
		name:  "main",
		block: new(codeblock),
	}

	c.fragments = append(c.fragments, c.fragment)
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
	if len(c.jumps) == 0 {
		return
	}

	for _, index := range c.jumps {
		ins := c.instrs[index]
		ins.operand = ins.target.startAt
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

		if !c.block.hasReturn {
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
		c.rewindControlFlow(false)
		c.compileExpr(node.Expr, true)
		c.add(bytecode.Return, 0)
		c.block.hasReturn = true

	case *ast.StopStmt:
		c.compileStopStmt(node)

	case *ast.NextStmt:
		c.compileNextStmt(node)

	case *ast.IfStmt:
		c.compileIfStmt(node)
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

	if addReturn && !c.block.hasReturn {
		c.add(bytecode.PushNone, 0)
		c.add(bytecode.Return, 0)
		c.block.hasReturn = true
	}
}

func (c *compiler) compileWhileStmt(node *ast.WhileStmt) {
	exit := new(codeblock)
	loop := new(codeblock)

	c.pushControlFlow(WHILE_LOOP, loop, exit)
	defer c.popControlFlow()

	c.useBlock(loop)
	c.compileExpr(node.Cond, true)
	c.jumpToBlock(bytecode.JumpIfFalse, exit)
	c.compileBlock(node.Body, false)
	c.jumpToBlock(bytecode.Jump, loop)
	c.useBlock(exit)
}

func (c *compiler) compileForStmt(node *ast.ForStmt) {
	exit := new(codeblock)
	setup := new(codeblock)
	loop := new(codeblock)

	c.pushControlFlow(FOR_LOOP, loop, exit)
	defer c.popControlFlow()

	c.useBlock(setup)
	c.compileExpr(node.Iterable, true)
	c.add(bytecode.NewIterator, 0)

	c.useBlock(loop)
	c.add(bytecode.Iterate, 0)
	c.jumpToBlock(bytecode.JumpIfFalse, exit)

	local := c.defineLocal(node.Element.Value, true)
	c.add(bytecode.SetLocal, local.index)

	c.compileBlock(node.Body, false)
	c.jumpToBlock(bytecode.Jump, loop)

	c.useBlock(exit)
}

func (c *compiler) compileSwitchStmt(node *ast.SwitchStmt) {
	endBlock := new(codeblock)

	lenCases := len(node.Cases) - 1
	for i, caseClause := range node.Cases {
		nextCaseBlock := new(codeblock)

		c.compileExpr(node.Key, true)
		c.compileExpr(caseClause.Value, true)

		callInfo := lang.NewCallInfo("==", 1)
		c.add(bytecode.CallMethod, c.addConstant(callInfo))

		c.jumpToBlock(bytecode.JumpIfFalse, nextCaseBlock)
		c.compileBlock(caseClause.Body, false)

		if i != lenCases || node.Default != nil {
			c.jumpToBlock(bytecode.Jump, endBlock)
		}

		c.useBlock(nextCaseBlock)
	}

	if node.Default != nil {
		c.compileBlock(node.Default.Body, false)
	}

	c.useBlock(endBlock)
}

func (c *compiler) compileIfStmt(node *ast.IfStmt) {
	var elseBlock, endBlock *codeblock

	if node.Else == nil {
		endBlock = new(codeblock)
		elseBlock = endBlock
	} else {
		elseBlock = new(codeblock)
		endBlock = new(codeblock)
	}

	c.compileExpr(node.Cond, true)
	c.jumpToBlock(bytecode.JumpIfFalse, elseBlock)
	c.compileBlock(node.Then, false)

	if node.Else != nil {
		c.jumpToBlock(bytecode.Jump, endBlock)
		c.useBlock(elseBlock)
		c.compileStmt(node.Else)
	}

	c.useBlock(endBlock)
}

func (c *compiler) jumpToBlock(op bytecode.Opcode, target *codeblock) {
	if c.block.hasReturn {
		return
	}

	c.jumps = append(c.jumps, len(c.instrs))
	c.instrs = append(c.instrs, &instr{opcode: op, target: target})
}

func (c *compiler) openScope(name string) {
	c.fragment = &fragment{
		name:     name,
		previous: c.fragment,
		block:    new(codeblock),
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

	end := new(codeblock)
	if len(node.Catches) != 0 {
		c.catchOffset = len(c.instrs)

		for _, ch := range node.Catches {
			catch := end
			end = new(codeblock)

			c.useBlock(catch)
			c.add(bytecode.MatchType, c.addConstant(ch.Type.Value))
			c.jumpToBlock(bytecode.JumpIfFalse, end)

			if ch.Ref != nil {
				local := c.defineLocal(ch.Ref.Value, true)
				c.add(bytecode.SetLocal, local.index)
			}

			c.compileBlock(ch.Body, true)
		}

		c.useBlock(end)
		c.add(bytecode.Throw, 0)
		c.block.hasReturn = true
	}

	return c.assemble()
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

func (c *compiler) cleanControlFlow() {
	if c.control == nil {
		return
	}

	if c.control.loop == FOR_LOOP {
		c.add(bytecode.Pop, 0)
	}
}

func (c *compiler) rewindControlFlow(inLoop bool) {
	for c.control != nil {
		switch c.control.loop {
		case FOR_LOOP, WHILE_LOOP:
			if inLoop {
				return
			}
		}

		c.cleanControlFlow()
		c.popControlFlow()
	}
}

func (c *compiler) compileReturnStmt(ret *ast.ReturnStmt) {
	c.rewindControlFlow(false)

	c.compileExpr(ret.Expr, true)
	c.add(bytecode.Return, 0)
	c.block.hasReturn = true
}

func (c *compiler) compileStopStmt(_stop *ast.StopStmt) {
	c.rewindControlFlow(true)

	if c.control == nil {
		panic("STOP OUTSIDE LOOP")
	}

	c.cleanControlFlow()
	c.jumpToBlock(bytecode.Jump, c.control.exit)
}

func (c *compiler) compileNextStmt(_next *ast.NextStmt) {
	c.rewindControlFlow(true)

	if c.control == nil {
		panic("NEXT OUTSIDE LOOP")
	}

	c.jumpToBlock(bytecode.Jump, c.control.start)
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

func (c *compiler) useBlock(b *codeblock) {
	b.startAt = byte(len(c.instrs))
	c.block = b
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

func (c *compiler) pushControlFlow(loop int, start, exit *codeblock) {
	if c.control == nil {
		c.control = &controlflow{
			loop:  loop,
			start: start,
			exit:  exit,
		}
		return
	}

	c.control = &controlflow{
		loop:  loop,
		start: start,
		exit:  exit,
		next:  c.control,
	}
}

func (c *compiler) popControlFlow() {
	if c.control == nil {
		return
	}

	c.control = c.control.next
}
