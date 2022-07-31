package compile

import (
	"errors"
	"fmt"
	"iracema/ast"
	"iracema/bytecode"
	"iracema/lang"
	"iracema/token"
	"strconv"
)

var unaryOps = map[token.Type]string{
	token.Plus:  "uadd",
	token.Minus: "usub",
	token.Not:   "unot",
}

var binaryOps = map[token.Type]string{
	token.Plus:       "+",
	token.Minus:      "-",
	token.Star:       "*",
	token.Slash:      "/",
	token.Equal:      "==",
	token.NotEqual:   "!=",
	token.Less:       "<",
	token.LessEqual:  "<=",
	token.Great:      ">",
	token.GreatEqual: ">=",
	token.Or:         "or",
	token.And:        "and",
}

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
	name         string
	argc         byte
	consts       []lang.IrObject
	paramIndices []byte
	locals       []*local
	instrs       []*instr
	jumps        []int
	catchOffset  int

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
	if err := c.compileStmt(file); err != nil {
		return nil, err
	}

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

func (c *compiler) compileStmt(stmt ast.Stmt) error {
	switch node := stmt.(type) {
	case *ast.File:
		if len(node.Stmts) == 0 {
			c.add(bytecode.PushNone, 0)
			c.add(bytecode.Return, 0)
			return nil
		}

		for _, stmt := range node.Stmts {
			if err := c.compileStmt(stmt); err != nil {
				return err
			}
		}

		if !c.block.hasReturn {
			c.add(bytecode.PushNone, 0)
			c.add(bytecode.Return, 0)
		}

	case *ast.ExprStmt:
		return c.compileExpr(node.Expr, false)

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
				if err := c.compileExpr(value, true); err != nil {
					return err
				}
				local.initialized = true
				c.add(bytecode.SetLocal, local.index)

			case *ast.IndexExpr:
				if err := c.compileExpr(lhs.Expr, true); err != nil {
					return err
				}

				if err := c.compileExpr(lhs.Index, true); err != nil {
					return err
				}

				if err := c.compileExpr(value, true); err != nil {
					return err
				}

				ci := lang.NewCallInfo("insert", 2)
				c.add(bytecode.CallMethod, c.addConstant(ci))
				c.add(bytecode.Pop, 0)
			}
		}

	case *ast.WhileStmt:
		return c.compileWhileStmt(node)

	case *ast.ForStmt:
		return c.compileForStmt(node)

	case *ast.SwitchStmt:
		return c.compileSwitchStmt(node)

	case *ast.ObjectDecl:
		return c.compileObjectDecl(node)

	case *ast.FunDecl:
		return c.compileFunDecl(node)

	case *ast.BlockStmt:
		return c.compileBlock(node, false)

	case *ast.ReturnStmt:
		return c.compileReturnStmt(node)

	case *ast.StopStmt:
		return c.compileStopStmt(node)

	case *ast.NextStmt:
		return c.compileNextStmt(node)

	case *ast.IfStmt:
		return c.compileIfStmt(node)
	}

	return nil
}

func (c *compiler) compileExpr(expr ast.Expr, isEvaluated bool) error {
	switch node := expr.(type) {
	case *ast.Ident:
		if node.IsAttr() {
			c.add(bytecode.GetAttr, c.addConstant(node.Value))
			return nil
		}

		if node.IsConstant() {
			c.add(bytecode.GetConstant, c.addConstant(node.Value))
			return nil
		}

		if local := c.resolve(node.Value); local != nil {
			if !local.initialized {
				return errors.New("underfined " + local.name)
			}

			c.add(bytecode.GetLocal, local.index)
			return nil
		}

		c.add(bytecode.PushSelf, 0)
		ci := lang.NewCallInfo(node.Value, 0)
		c.add(bytecode.CallMethod, c.addConstant(ci))

	case *ast.BasicLit:
		if isEvaluated {
			if err := c.compileLiteral(node); err != nil {
				return err
			}
		}

	case *ast.UnaryExpr:
		if err := c.compileExpr(node.Expr, true); err != nil {
			return err
		}
		c.addUnary(node.Operator)

	case *ast.BinaryExpr:
		if err := c.compileBinaryExpr(node); err != nil {
			return err
		}

		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	case *ast.ArrayLit:
		return c.compileArrayLireral(node)

	case *ast.GroupExpr:
		return c.compileExpr(node.Expr, isEvaluated)

	case *ast.BlockExpr:
		return errors.New("not implemented")

	case *ast.IndexExpr:
		if err := c.compileExpr(node.Expr, true); err != nil {
			return err
		}

		if err := c.compileExpr(node.Index, true); err != nil {
			return err
		}

		ci := lang.NewCallInfo("get", 1)
		c.add(bytecode.CallMethod, c.addConstant(ci))
		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	case *ast.CallExpr:
		if node.Receiver != nil {
			if err := c.compileExpr(node.Receiver, true); err != nil {
				return err
			}
		} else {
			c.add(bytecode.PushSelf, 0)
		}

		for _, arg := range node.Arguments {
			if err := c.compileExpr(arg, true); err != nil {
				return err
			}
		}

		ci := lang.NewCallInfo(node.Method.Value, byte(len(node.Arguments)))
		c.add(bytecode.CallMethod, c.addConstant(ci))
		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	case *ast.SuperExpr:
		c.add(bytecode.PushSelf, 0)

		var ci lang.IrObject
		if node.ExplicitArgs {
			ci = lang.NewCallInfo(c.name, byte(len(node.Arguments)))

			for _, arg := range node.Arguments {
				if err := c.compileExpr(arg, true); err != nil {
					return err
				}
			}
		} else {
			ci = lang.NewCallInfo(c.name, byte(len(c.paramIndices)))

			for _, index := range c.paramIndices {
				c.add(bytecode.GetLocal, index)
			}
		}

		c.add(bytecode.CallSuper, c.addConstant(ci))
		if !isEvaluated {
			c.add(bytecode.Pop, 0)
		}

	default:
		return errors.New("unknown expr")
	}

	return nil
}

func (c *compiler) addUnary(t *token.Token) {
	ci := lang.NewCallInfo(unaryOps[t.Type], 0)
	c.add(bytecode.CallMethod, c.addConstant(ci))
}

func (c *compiler) compileBlock(block *ast.BlockStmt, addReturn bool) error {
	for _, stmt := range block.Stmts {
		if err := c.compileStmt(stmt); err != nil {
			return err
		}
	}

	if addReturn && !c.block.hasReturn {
		c.add(bytecode.PushNone, 0)
		c.add(bytecode.Return, 0)
		c.block.hasReturn = true
	}

	return nil
}

func (c *compiler) compileWhileStmt(node *ast.WhileStmt) error {
	exit := new(codeblock)
	loop := new(codeblock)

	c.pushControlFlow(WHILE_LOOP, loop, exit)
	defer c.popControlFlow()

	c.useBlock(loop)
	if err := c.compileConditional(node.Cond, exit); err != nil {
		return err
	}

	if err := c.compileBlock(node.Body, false); err != nil {
		return err
	}

	c.jumpToBlock(bytecode.Jump, loop)
	c.useBlock(exit)

	return nil
}

func (c *compiler) compileForStmt(node *ast.ForStmt) error {
	exit := new(codeblock)
	setup := new(codeblock)
	loop := new(codeblock)

	c.pushControlFlow(FOR_LOOP, loop, exit)
	defer c.popControlFlow()

	c.useBlock(setup)
	if err := c.compileExpr(node.Iterable, true); err != nil {
		return err
	}
	c.add(bytecode.NewIterator, 0)

	c.useBlock(loop)
	c.add(bytecode.Iterate, 0)
	c.jumpToBlock(bytecode.JumpIfFalse, exit)

	local := c.defineLocal(node.Element.Value, true)
	c.add(bytecode.SetLocal, local.index)

	if err := c.compileBlock(node.Body, false); err != nil {
		return err
	}

	c.jumpToBlock(bytecode.Jump, loop)
	c.useBlock(exit)

	return nil
}

func (c *compiler) compileSwitchStmt(node *ast.SwitchStmt) error {
	endBlock := new(codeblock)

	lenCases := len(node.Cases) - 1
	for i, caseClause := range node.Cases {
		nextCaseBlock := new(codeblock)

		if err := c.compileExpr(node.Key, true); err != nil {
			return err
		}

		if err := c.compileExpr(caseClause.Value, true); err != nil {
			return err
		}

		callInfo := lang.NewCallInfo("==", 1)
		c.add(bytecode.CallMethod, c.addConstant(callInfo))

		c.jumpToBlock(bytecode.JumpIfFalse, nextCaseBlock)
		if err := c.compileBlock(caseClause.Body, false); err != nil {
			return err
		}

		if i != lenCases || node.Default != nil {
			c.jumpToBlock(bytecode.Jump, endBlock)
		}

		c.useBlock(nextCaseBlock)
	}

	if node.Default != nil {
		if err := c.compileBlock(node.Default.Body, false); err != nil {
			return err
		}
	}

	c.useBlock(endBlock)
	return nil
}

func (c *compiler) compileIfStmt(node *ast.IfStmt) error {
	var elseBlock, endBlock *codeblock

	if node.Else == nil {
		endBlock = new(codeblock)
		elseBlock = endBlock
	} else {
		elseBlock = new(codeblock)
		endBlock = new(codeblock)
	}

	if err := c.compileConditional(node.Cond, elseBlock); err != nil {
		return err
	}

	if err := c.compileBlock(node.Then, false); err != nil {
		return err
	}

	if node.Else != nil {
		c.jumpToBlock(bytecode.Jump, endBlock)
		c.useBlock(elseBlock)
		if err := c.compileStmt(node.Else); err != nil {
			return err
		}
	}

	c.useBlock(endBlock)
	return nil
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

func (c *compiler) compileObjectDecl(obj *ast.ObjectDecl) error {
	if obj.Parent != nil {
		c.add(bytecode.GetConstant, c.addConstant(obj.Parent.Value))
	} else {
		c.add(bytecode.PushNone, 0)
	}

	c.openScope(obj.Name.Value)
	if err := c.compileBlock(obj.Body, true); err != nil {
		return err
	}

	objBody := c.assemble()
	c.closeScope()

	c.add(bytecode.DefineObject, c.addConstant(objBody))
	c.add(bytecode.Pop, 0)

	return nil
}

func (c *compiler) compileFunDecl(node *ast.FunDecl) error {
	c.openScope(node.Name.Value)

	c.argc = byte(len(node.Parameters))
	for _, param := range node.Parameters {
		p := c.defineLocal(param.Value, true)
		c.paramIndices = append(c.paramIndices, p.index)
	}

	if err := c.compileBlock(node.Body, true); err != nil {
		return err
	}

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

			if err := c.compileBlock(ch.Body, true); err != nil {
				return err
			}
		}

		c.useBlock(end)
		c.add(bytecode.Throw, 0)
		c.block.hasReturn = true
	}

	fun := c.assemble()
	c.closeScope()
	c.add(bytecode.DefineFunction, c.addConstant(fun))
	return nil
}

func (c *compiler) compileArrayLireral(node *ast.ArrayLit) error {
	size := len(node.Elements)
	for _, el := range node.Elements {
		if err := c.compileExpr(el, true); err != nil {
			return err
		}
	}

	c.add(bytecode.BuildArray, byte(size))
	return nil
}

func (c *compiler) compileLiteral(lit *ast.BasicLit) error {
	var val lang.IrObject

	switch lit.Type() {
	case token.String:
		val = lang.NewString(lit.Value)

	case token.Bool:
		value, err := strconv.ParseBool(lit.Value)
		if err != nil {
			return err
		}
		val = lang.Bool(value)

	case token.Int:
		value, err := strconv.Atoi(lit.Value)
		if err != nil {
			return err
		}
		val = lang.Int(value)

	case token.Float:
		value, err := strconv.ParseFloat(lit.Value, 64)
		if err != nil {
			return err
		}
		val = lang.Float(value)

	case token.None:
		c.add(bytecode.PushNone, 0)
		return nil

	default:
		return errors.New("invalid literal")
	}

	c.add(bytecode.Push, c.addConstant(val))
	return nil
}

func (c *compiler) compileConditional(expr ast.Expr, next *codeblock) error {
	switch x := expr.(type) {
	case *ast.BinaryExpr:
		switch x.Operator.Type {
		case token.And:
			if err := c.compileExpr(x.Left, true); err != nil {
				return err
			}
			c.jumpToBlock(bytecode.JumpIfFalse, next)

			if err := c.compileExpr(x.Right, true); err != nil {
				return err
			}
			c.jumpToBlock(bytecode.JumpIfFalse, next)

		case token.Or:
			body := new(codeblock)
			if err := c.compileExpr(x.Left, true); err != nil {
				return err
			}
			c.jumpToBlock(bytecode.JumpIfTrue, body)

			if err := c.compileExpr(x.Right, true); err != nil {
				return err
			}

			c.jumpToBlock(bytecode.JumpIfFalse, next)
			c.useBlock(body)

		default:
			if err := c.compileBinaryExpr(x); err != nil {
				return err
			}
			c.jumpToBlock(bytecode.JumpIfFalse, next)
		}

	default:
		if err := c.compileExpr(x, true); err != nil {
			return err
		}
		c.jumpToBlock(bytecode.JumpIfFalse, next)
	}

	return nil
}

func (c *compiler) compileBinaryExpr(expr *ast.BinaryExpr) error {
	if err := c.compileExpr(expr.Left, true); err != nil {
		return err
	}

	if err := c.compileExpr(expr.Right, true); err != nil {
		return err
	}

	ci := lang.NewCallInfo(binaryOps[expr.Operator.Type], 1)
	c.add(bytecode.CallMethod, c.addConstant(ci))
	return nil
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

func (c *compiler) compileReturnStmt(ret *ast.ReturnStmt) error {
	c.rewindControlFlow(false)

	if ret.Value != nil {
		if err := c.compileExpr(ret.Value, true); err != nil {
			return err
		}
	} else {
		c.add(bytecode.PushNone, 0)
	}

	c.add(bytecode.Return, 0)
	c.block.hasReturn = true

	return nil
}

func (c *compiler) compileStopStmt(_stop *ast.StopStmt) error {
	c.rewindControlFlow(true)

	if c.control == nil {
		return errors.New("STOP OUTSIDE LOOP")
	}

	c.cleanControlFlow()
	c.jumpToBlock(bytecode.Jump, c.control.exit)

	return nil
}

func (c *compiler) compileNextStmt(_next *ast.NextStmt) error {
	c.rewindControlFlow(true)

	if c.control == nil {
		return errors.New("NEXT OUTSIDE LOOP")
	}

	c.jumpToBlock(bytecode.Jump, c.control.start)

	return nil
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
