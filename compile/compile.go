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

	TOP_SCOPE = 1 << iota
	CLASS_SCOPE
	METHOD_SCOPE
)

type controlflow struct {
	loop  int
	start *basicblock
	exit  *basicblock

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

type instr struct {
	opcode  bytecode.Opcode
	operand byte
	target  *basicblock
}

func (i *instr) hasTarget() bool {
	return i.target != nil
}

type fragment struct {
	name         string
	scope        int
	argc         byte
	consts       []lang.IrObject
	paramIndices []byte
	locals       []*local
	catchOffset  int

	control    *controlflow
	entrypoint *basicblock // ref to first block
	block      *basicblock
	previous   *fragment
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
	blk := new(basicblock)
	c.fragment = &fragment{
		name:       "main",
		scope:      TOP_SCOPE,
		block:      blk,
		entrypoint: blk,
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
	markReachable(c.entrypoint)
	c.patchJumps()

	for block := c.entrypoint; block != nil; block = block.next {
		for _, instr := range block.instrs {
			code := uint16(instr.opcode)<<8 | uint16(instr.operand)
			bytecode = append(bytecode, code)
		}
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

type stack []*basicblock

func (s *stack) Push(block *basicblock) {
	*s = append(*s, block)
}

func (s *stack) Pop() *basicblock {
	i := len(*s) - 1
	elem := (*s)[i]
	*s = (*s)[:i]
	return elem
}

func (s *stack) Empty() bool { return len(*s) == 0 }

func skipEmptyBlocks(entrypoint *basicblock) {
	// eliminate empty blocks
	for block := entrypoint; block != nil; block = block.next {
		next := block.next
		if next == nil {
			break
		}

		for len(next.instrs) == 0 && next.next != nil {
			next = next.next
		}

		block.next = next
	}

	for block := entrypoint; block != nil; block = block.next {
		if len(block.instrs) == 0 {
			continue
		}

		for _, ins := range block.instrs {
			if ins.opcode == bytecode.Jump || ins.opcode == bytecode.JumpIfTrue || ins.opcode == bytecode.JumpIfFalse {
				for len(ins.target.instrs) == 0 {
					ins.target = ins.target.next
				}
			}
		}
	}
}

func markReachable(entrypoint *basicblock) {
	if entrypoint == nil {
		return
	}

	skipEmptyBlocks(entrypoint)
	defer skipEmptyBlocks(entrypoint)

	//mark rachable
	s := new(stack)

	entrypoint.reachable = true
	s.Push(entrypoint)

	for !s.Empty() {
		block := s.Pop()
		block.visited = true

		if block.next != nil && block.hasFallthrough() {
			if next := block.next; !next.visited {
				next.reachable = true
				s.Push(next)
			}
		}

		for _, ins := range block.instrs {
			if ins.hasTarget() {
				if target := ins.target; !target.visited {
					target.reachable = true
					s.Push(target)
				}
			}
		}
	}

	for block := entrypoint; block != nil; block = block.next {
		if block.reachable {
			continue
		}

		block.instrs = nil
	}
}

func (c *compiler) patchJumps() {
	var start int
	for block := c.entrypoint; block != nil; block = block.next {
		block.offset = start
		start += len(block.instrs)
	}

	for block := c.entrypoint; block != nil; block = block.next {
		for _, ins := range block.instrs {
			if ins.hasTarget() {
				if ins.opcode == bytecode.WithCatch {
					ins.opcode = bytecode.Nop
					c.catchOffset = ins.target.offset
					continue
				}

				ins.operand = byte(ins.target.offset)
			}
		}
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

		c.add(bytecode.PushNone, 0)
		c.add(bytecode.Return, 0)

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

/*
* CFG for the following snippet
*
* a = 100
* while a > 0 {
*   puts(a)
*
*   a = a - 1
* }
*
*      ┌───────────────────────────────────────────────┐
* ┌────►  0000 PUSH                100                 │
* │    │  0002 SET_LOCAL           a@0                 │
* │    │  0004 GET_LOCAL           a@0                 │
* │    │  0006 PUSH                0                   │
* │    │  0008 CALL_METHOD         name: > argc: 1     │
* │    │  0010 JUMP_IF_FALSE       30                  ├────┐
* │    └──────────────────────┬────────────────────────┘    │
* │                           │                             │
* │                          next                           │
* │                           │                             │
* │    ┌──────────────────────▼────────────────────────┐    │
* │    │  0012 PUSH_SELF                               │    │
* │    │  0014 GET_LOCAL           a@0                 │    │
* │    │  0016 CALL_METHOD         name: puts argc: 1  │    │
* │    │  0018 POP                                     │    │
* │    │  0020 GET_LOCAL           a@0                 │    │
* │    │  0022 PUSH                1                   │    │
* │    │  0024 CALL_METHOD         name: - argc: 1     │    │
* │    │  0026 SET_LOCAL           a@0                 │    │
* └────┤  0028 JUMP                4                   │    │
*      └──────────────────────┬────────────────────────┘    │
*                             │                             │
*                            next                           │
*                             │                             │
*      ┌──────────────────────▼────────────────────────┐    │
*      │  0030 PUSH_NONE                               ◄────┘
*      │  0032 RETURN                                  │
*      └───────────────────────────────────────────────┘
* */

func (c *compiler) compileWhileStmt(node *ast.WhileStmt) error {
	cond := new(basicblock)
	loop := new(basicblock)
	exit := new(basicblock)

	c.pushControlFlow(WHILE_LOOP, cond, exit)
	defer c.popControlFlow()

	c.useBlock(cond)
	if err := c.compileConditional(node.Cond, exit); err != nil {
		return err
	}

	c.useBlock(loop)
	if err := c.compileBlock(node.Body, false); err != nil {
		return err
	}

	c.addJump(bytecode.Jump, cond)
	c.useBlock(exit)

	return nil
}

/*
* CFG for the following snippet
*
* for el in [1,2,3] {
*   puts(el)
* }
*
*
*       ┌──────────────────────────────────────────────┐
* ┌─────►  0000 PUSH               1                   │
* │     │  0002 PUSH               2                   │
* │     │  0004 PUSH               3                   │
* │     │  0006 BUILD_ARRAY        size: 3             │
* │     │  0012 NEWITERATOR                            │
* │     │  0014 ITERATE                                │
* │     │  0016 JUMP_IF_FALSE      30                  ├────┐
* │     └─────────────────────┬────────────────────────┘    │
* │                           │                             │
* │                          next                           │
* │                           │                             │
* │     ┌─────────────────────▼────────────────────────┐    │
* │     │  0018 SET_LOCAL          el@1                │    │
* │     │  0020 PUSH_SELF                              │    │
* │     │  0022 GET_LOCAL          el@1                │    │
* │     │  0024 CALL_METHOD        name: puts argc: 1  │    │
* │     │  0026 POP                                    │    │
* └─────┤  0028 JUMP               14                  │    │
*       └─────────────────────┬────────────────────────┘    │
*                             │                             │
*                            next                           │
*                             │                             │
*       ┌─────────────────────▼────────────────────────┐    │
*       │  0030 PUSH_NONE                              ◄────┘
*       │  0032 RETURN                                 │
*       └──────────────────────────────────────────────┘
**/

func (c *compiler) compileForStmt(node *ast.ForStmt) error {
	exit := new(basicblock)
	setup := new(basicblock)
	loop := new(basicblock)

	c.pushControlFlow(FOR_LOOP, loop, exit)
	defer c.popControlFlow()

	c.useBlock(setup)
	if err := c.compileExpr(node.Iterable, true); err != nil {
		return err
	}
	c.add(bytecode.NewIterator, 0)

	c.useBlock(loop)
	c.add(bytecode.Iterate, 0)
	c.addJump(bytecode.JumpIfFalse, exit)

	local := c.defineLocal(node.Element.Value, true)
	c.add(bytecode.SetLocal, local.index)

	if err := c.compileBlock(node.Body, false); err != nil {
		return err
	}

	c.addJump(bytecode.Jump, loop)
	c.useBlock(exit)

	return nil
}

/*
* CFG for the following snippet
*
* n = 10
* switch n {
* case 10: puts(10)
* case 20: puts(20)
* default: puts("DEFAULT")
* }
*          ┌───────────────────────────────────────────────┐
*          │  0000 PUSH                10                  │
*          │  0002 SET_LOCAL           n@0                 │
*          │  0004 GET_LOCAL           n@0                 │
*          │  0006 PUSH                10                  │
*          │  0008 CALL_METHOD         name: == argc: 1    │
*          │  0010 JUMP_IF_FALSE       22                  ├───┐
*          └─────────────────────┬─────────────────────────┘   │
*                                │                             │
*                               next                           │
*                                │                             │
*          ┌─────────────────────▼─────────────────────────┐   │
*          │  0012 PUSH_SELF                               │   │
*          │  0014 PUSH                10                  │   │
*          │  0016 CALL_METHOD         name: puts argc: 1  │   │
*          │  0018 POP                                     │   │
* ┌────────┤  0020 JUMP                48                  │   │
* │        └─────────────────────┬─────────────────────────┘   │
* │                              │                             │
* │                             next                           │
* │                              │                             │
* │        ┌─────────────────────▼─────────────────────────┐   │
* │        │  0022 GET_LOCAL           n@0                 ◄───┘
* │        │  0024 PUSH                20                  │
* │        │  0026 CALL_METHOD         name: == argc: 1    │
* │        │  0028 JUMP_IF_FALSE       40                  ├───┐
* │        └─────────────────────┬─────────────────────────┘   │
* │                              │                             │
* │                             next                           │
* │                              │                             │
* │        ┌─────────────────────▼─────────────────────────┐   │
* │        │  0030 PUSH_SELF                               │   │
* │        │  0032 PUSH                20                  │   │
* │        │  0034 CALL_METHOD         name: puts argc: 1  │   │
* │        │  0036 POP                                     │   │
* │  ┌─────┤  0038 JUMP                48                  │   │
* │  │     └─────────────────────┬─────────────────────────┘   │
* │  │                           │                             │
* │  │                          next                           │
* │  │                           │                             │
* │  │     ┌─────────────────────▼─────────────────────────┐   │
* │  │     │  0040 PUSH_SELF                               ◄───┘
* │  │     │  0042 PUSH                "DEFAULT"           │
* │  │     │  0044 CALL_METHOD         name: puts argc: 1  │
* │  │     │  0046 POP                                     │
* │  │     └─────────────────────┬─────────────────────────┘
* │  │                           │
* │  │                          next
* │  │                           │
* │  │     ┌─────────────────────▼─────────────────────────┐
* └──┴─────►  0048 PUSH_NONE                               │
*          │  0050 RETURN                                  │
*          └───────────────────────────────────────────────┘
**/

func (c *compiler) compileSwitchStmt(node *ast.SwitchStmt) error {
	endBlock := new(basicblock)

	lenCases := len(node.Cases) - 1
	for i, caseClause := range node.Cases {
		nextCase := new(basicblock)

		if err := c.compileExpr(node.Key, true); err != nil {
			return err
		}

		if err := c.compileExpr(caseClause.Value, true); err != nil {
			return err
		}

		callInfo := lang.NewCallInfo("==", 1)
		c.add(bytecode.CallMethod, c.addConstant(callInfo))

		c.addJump(bytecode.JumpIfFalse, nextCase)
		if err := c.compileBlock(caseClause.Body, false); err != nil {
			return err
		}

		if i != lenCases || node.Default != nil {
			c.addJump(bytecode.Jump, endBlock)
		}

		c.useBlock(nextCase)
	}

	if node.Default != nil {
		if err := c.compileBlock(node.Default.Body, false); err != nil {
			return err
		}
	}

	c.useBlock(endBlock)
	return nil
}

/*
* CFG for the following snippet
*
* a = 100
* if a > 10 {
* 	puts("BIGGER")
* } else {
* 	puts("SMALLER")
* }
*
*         ┌───────────────────────────────────────────────┐
*         │ 0000  PUSH                100                 │
*         │ 0002  SET_LOCAL           a@1                 │
*         │ 0004  GET_LOCAL           a@1                 │
*         │ 0006  PUSH                10                  │
*         │ 0008  CALL_METHOD         name: > argc: 1     │
*         │ 0010  JUMP_IF_FALSE       22                  ├─────┐
*         └──────────────────────┬────────────────────────┘     │
*                                │                              │
*                              next                             │
*                                │                              │
*         ┌──────────────────────▼────────────────────────┐     │
*         │ 0012  PUSH_SELF                               │     │
*         │ 0014  PUSH                "BIGGER"            │     │
*         │ 0016  CALL_METHOD         name: puts argc: 1  │     │
*         │ 0018  POP                                     │     │
*    ┌────┤ 0020  JUMP                30                  │     │
*    │    └──────────────────────┬────────────────────────┘     │
*    │                           │                              │
*    │                         next                             │
*    │                           │                              │
*    │    ┌──────────────────────▼────────────────────────┐     │
*    │    │ 0022  PUSH_SELF                               ◄─────┘
*    │    │ 0024  PUSH                "SMALLER"           │
*    │    │ 0026  CALL_METHOD         name: puts argc: 1  │
*    │    │ 0028  POP                                     │
*    │    └──────────────────────┬────────────────────────┘
*    │                           │
*    │                         next
*    │                           │
*    │    ┌──────────────────────▼────────────────────────┐
*    └────► 0030 PUSH_NONE                                │
*         │ 0032 RETURN                                   │
*         └───────────────────────────────────────────────┘
* */

func (c *compiler) compileIfStmt(node *ast.IfStmt) error {
	var elseBlock, endBlock *basicblock

	if node.Else == nil {
		endBlock = new(basicblock)
		elseBlock = endBlock
	} else {
		elseBlock = new(basicblock)
		endBlock = new(basicblock)
	}

	if err := c.compileConditional(node.Cond, elseBlock); err != nil {
		return err
	}

	if err := c.compileBlock(node.Then, false); err != nil {
		return err
	}

	if node.Else != nil {
		c.addJump(bytecode.Jump, endBlock)
		c.useBlock(elseBlock)
		if err := c.compileStmt(node.Else); err != nil {
			return err
		}
	}

	c.useBlock(endBlock)
	return nil
}

func (c *compiler) addJump(op bytecode.Opcode, target *basicblock) {
	if !c.block.hasFallthrough() {
		c.useBlock(new(basicblock))
	}

	c.block.instrs = append(c.block.instrs, &instr{opcode: op, target: target})
}

func (c *compiler) openScope(name string, scope int) {
	blk := new(basicblock)
	c.fragment = &fragment{
		name:       name,
		scope:      scope,
		previous:   c.fragment,
		block:      blk,
		entrypoint: blk,
	}

	c.fragments = append(c.fragments, c.fragment)
}

func (c *compiler) closeScope() {
	c.fragment = c.fragment.previous
}

func (c *compiler) compileObjectDecl(obj *ast.ObjectDecl) error {
	if c.scope != TOP_SCOPE {
		return errors.New("can only there declare object in the top most scope")
	}

	if obj.Parent != nil {
		c.add(bytecode.GetConstant, c.addConstant(obj.Parent.Value))
	} else {
		c.add(bytecode.PushNone, 0)
	}

	c.openScope(obj.Name.Value, CLASS_SCOPE)
	if err := c.compileBlock(obj.Body, true); err != nil {
		return err
	}

	objBody := c.assemble()
	c.closeScope()

	c.add(bytecode.DefineObject, c.addConstant(objBody))
	c.add(bytecode.Pop, 0)

	return nil
}

func (c *compiler) compileFunDecl(fun *ast.FunDecl) error {
	if c.scope == METHOD_SCOPE {
		return errors.New("can not declare a method inside of a method")
	}

	c.openScope(fun.Name.Value, METHOD_SCOPE)

	var catch *basicblock
	if len(fun.Catches) != 0 {
		catch = new(basicblock)
		c.addJump(bytecode.WithCatch, catch)
	}

	c.argc = byte(len(fun.Parameters))
	for _, param := range fun.Parameters {
		p := c.defineLocal(param.Value, true)
		c.paramIndices = append(c.paramIndices, p.index)
	}

	if err := c.compileBlock(fun.Body, true); err != nil {
		return err
	}

	if catch != nil {
		c.useBlock(catch)
		for _, ch := range fun.Catches {
			catch = new(basicblock)

			c.add(bytecode.MatchType, c.addConstant(ch.Type.Value))
			c.addJump(bytecode.JumpIfFalse, catch)

			if ch.Ref != nil {
				local := c.defineLocal(ch.Ref.Value, true)
				c.add(bytecode.SetLocal, local.index)
			}

			if err := c.compileBlock(ch.Body, true); err != nil {
				return err
			}

			c.useBlock(catch)
		}

		c.add(bytecode.Throw, 0)
		c.block.hasReturn = true
	}

	method := c.assemble()
	c.closeScope()
	c.add(bytecode.DefineFunction, c.addConstant(method))
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

func (c *compiler) compileConditional(expr ast.Expr, next *basicblock) error {
	switch x := expr.(type) {
	case *ast.BinaryExpr:
		switch x.Operator.Type {
		case token.And:
			if err := c.compileExpr(x.Left, true); err != nil {
				return err
			}
			c.addJump(bytecode.JumpIfFalse, next)

			if err := c.compileExpr(x.Right, true); err != nil {
				return err
			}
			c.addJump(bytecode.JumpIfFalse, next)

		case token.Or:
			body := new(basicblock)
			if err := c.compileExpr(x.Left, true); err != nil {
				return err
			}
			c.addJump(bytecode.JumpIfTrue, body)

			if err := c.compileExpr(x.Right, true); err != nil {
				return err
			}

			c.addJump(bytecode.JumpIfFalse, next)
			c.useBlock(body)

		default:
			if err := c.compileBinaryExpr(x); err != nil {
				return err
			}
			c.addJump(bytecode.JumpIfFalse, next)
		}

	default:
		if err := c.compileExpr(x, true); err != nil {
			return err
		}
		c.addJump(bytecode.JumpIfFalse, next)
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
	c.addJump(bytecode.Jump, c.control.exit)

	return nil
}

func (c *compiler) compileNextStmt(_next *ast.NextStmt) error {
	c.rewindControlFlow(true)

	if c.control == nil {
		return errors.New("NEXT OUTSIDE LOOP")
	}

	c.addJump(bytecode.Jump, c.control.start)

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

func (c *compiler) useBlock(next *basicblock) {
	c.block.next = next
	c.block = next
}

func (c *compiler) add(opcode bytecode.Opcode, operand byte) {
	if !c.block.hasFallthrough() {
		c.useBlock(new(basicblock))
	}

	ins := &instr{opcode: opcode, operand: operand}
	c.block.instrs = append(c.block.instrs, ins)
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

func (c *compiler) pushControlFlow(loop int, start, exit *basicblock) {
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
