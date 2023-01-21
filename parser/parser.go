package parser

import (
	"errors"
	"fmt"
	"io"
	"iracema/ast"
	"iracema/lexer"
	"iracema/token"
	"os"
)

var startDecl = map[token.Type]bool{
	token.Var:    true,
	token.Fun:    true,
	token.Object: true,
}

var startStmt = map[token.Type]bool{
	token.Object: true,
	token.Fun:    true,
	token.If:     true,
	token.For:    true,
	token.While:  true,
	token.Switch: true,
	token.Stop:   true,
	token.Return: true,
}

var closingToken = map[token.Type]bool{
	token.RightParen:   true,
	token.RightBrace:   true,
	token.RightBracket: true,
}

var switchStartStmt = map[token.Type]bool{
	token.Case:       true,
	token.Colon:      true,
	token.Default:    true,
	token.RightBrace: true,
}

type parser struct {
	lexer lexer.Lexer
	tok   *token.Token
	err   error
}

func ParseFile(filename string) (*ast.File, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return Parse(f)
}

func Parse(input io.Reader) (*ast.File, error) {
	p := new(parser)
	p.init(input)

	return p.parse(), p.err
}

func (p *parser) init(source io.Reader) {
	p.lexer = lexer.New(source, p.setError)
	p.advance()
}

func isDone(tok *token.Token) bool {
	return tok.Type == token.EOF ||
		tok.Type == token.RightBrace ||
		tok.Type == token.Case ||
		tok.Type == token.Default
}

func (p *parser) parse() *ast.File {
	file := new(ast.File)

	for p.tok.Type != token.EOF {
		switch p.tok.Type {
		case token.Use:
			for p.consume(token.Use) {
				name := p.expect(token.String)
				file.Imports = append(file.Imports, name.Literal)

				p.consume(token.NewLine)
			}
		case token.Var:
			file.VarList = append(file.VarList, p.parseVarDecl())

		case token.Fun:
			file.FunList = append(file.FunList, p.parseFunDecl())

		case token.Object:
			file.ObjectList = append(file.ObjectList, p.parseObjectDecl())

		default:
			mesg := fmt.Sprintf("unexpected %s, expecting VarDecl, FunDecl or ObjectDecl", p.tok)
			p.setError(p.tok.Position, mesg)
			p.sync(startDecl)
			continue
		}

		if !p.consume(token.NewLine) && p.tok.Type != token.EOF {
			err := fmt.Sprintf("unexpected %s, expecting EOF or new line", p.tok)
			p.setError(p.tok.Position, err)
			p.sync(startDecl)
		}
	}

	p.expect(token.EOF)

	return file
}

func (p *parser) parseStmtList() (list []ast.Stmt) {
	for !isDone(p.tok) {
		stmt := p.parseStmt()
		if stmt == nil {
			return
		}

		list = append(list, stmt)
		if !p.consume(token.NewLine) && p.tok.Type != token.RightBrace {
			err := fmt.Sprintf("unexpected %s, expecting } or new line", p.tok)
			p.setError(p.tok.Position, err)
			p.sync(startStmt)
			continue
		}
	}

	return
}

func (p *parser) sync(to map[token.Type]bool) {
	for ; !p.at(token.EOF); p.advance() {
		if to[p.tok.Type] {
			p.advance()
			break
		}
	}
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.tok.Type {
	case token.Object:
		return p.parseObjectDecl()

	case token.Var:
		return p.parseVarDecl()

	case token.Fun:
		return p.parseFunDecl()

	case token.If:
		return p.parseIfStmt()

	case token.While:
		return p.parseWhileStmt()

	case token.For:
		return p.parseForStmt()

	case token.Switch:
		return p.parseSwitchStmt()

	case token.Stop:
		stop := new(ast.StopStmt)
		stop.Pos = p.tok.Position
		p.advance()
		return stop

	case token.Next:
		next := new(ast.NextStmt)
		next.Pos = p.tok.Position
		p.advance()
		return next

	case token.Return:
		return p.parseReturnStmt()

	case
		token.Ident, token.String, token.Bool, token.Int, token.Float,
		token.LeftParen, token.Not, token.Plus, token.Minus,
		token.LeftBracket, token.LeftBrace, token.None, token.Block,
		token.Super, token.This:
		return p.parseSimpleStmt()

	default:
		return nil
	}
}

func (p *parser) parseObjectDecl() *ast.ObjectDecl {
	p.expect(token.Object)

	obj := new(ast.ObjectDecl)
	obj.Name = p.parseIdent(true)

	if p.consume(token.Is) {
		obj.Parent = p.parseIdent(true)
	}

	p.expect(token.LeftBrace)
	for p.tok.Type != token.RightBrace && p.tok.Type != token.EOF {
		switch p.tok.Type {
		case token.Var:
			obj.FieldList = append(obj.FieldList, p.parseVarDecl())

		case token.Fun:
			obj.FunList = append(obj.FunList, p.parseFunDecl())
		case token.NewLine:
			p.advance()
			continue

		default:
			mesg := fmt.Sprintf("unexpected %s, expecting FunDecl or Field", p.tok)
			p.setError(p.tok.Position, mesg)
			return obj
		}
	}

	p.expect(token.RightBrace)

	return obj
}

func (p *parser) parseVarDecl() *ast.VarDecl {
	tok := p.expect(token.Var)

	decl := new(ast.VarDecl)
	decl.Pos = tok.Position
	decl.Name = p.parseIdent(false)
	if p.consume(token.Assign) {
		decl.Value = p.parseExpr()
	} else {
		decl.Type = p.parseType()
		if p.consume(token.Assign) {
			decl.Value = p.parseExpr()
		}
	}

	return decl
}

func (p *parser) parseType() ast.Expr {
	name := p.parseIdent(true)

	switch p.tok.Type {
	case token.Less:
		p.advance()
		paramType := p.parseIdent(true)
		p.expect(token.Great)
		return &ast.Type{BaseType: name, ParamType: paramType}

	default:
		return name
	}
}

func (p *parser) parseFunDecl() *ast.FunDecl {
	funToken := p.expect(token.Fun)

	fun := new(ast.FunDecl)
	fun.Pos = funToken.Position
	fun.Name = p.parseIdent(false)
	fun.Parameters = p.parseParameterList()

	if p.consume(token.Arrow) {
		fun.Return = p.parseIdent(true)
	}

	fun.Body = p.parseBlockStmt()
	fun.Catches = p.parseCatchList()
	return fun
}

func (p *parser) parseCatchList() (list []*ast.CatchDecl) {
	for p.consume(token.Catch) {
		p.expect(token.LeftParen)
		ref := p.parseIdent(false)
		p.expect(token.Colon)
		typ := p.parseIdent(true)
		p.expect(token.RightParen)
		catch := &ast.CatchDecl{
			Ref:  ref,
			Type: typ,
			Body: p.parseBlockStmt(),
		}

		list = append(list, catch)
	}

	return
}

func (p *parser) parseBlockStmt() *ast.BlockStmt {
	lBrace := p.expect(token.LeftBrace)
	stmts := p.parseStmtList()
	rBrace := p.expect(token.RightBrace)

	block := new(ast.BlockStmt)
	block.Pos = lBrace.Position
	block.Stmts = stmts
	block.RightBrace = rBrace.Position
	return block
}

func (p *parser) parseIfStmt() ast.Stmt {
	p.expect(token.If)
	predicate := p.parseExpr()
	consequent := p.parseBlockStmt()

	var alternative ast.Stmt
	if p.at(token.Else) {
		p.expect(token.Else)

		switch p.tok.Type {
		case token.If:
			alternative = p.parseIfStmt()
		case token.LeftBrace:
			alternative = p.parseBlockStmt()
		default:
			p.setError(p.tok.Position, "expected left brace or if statement")
			p.advance()
		}
	}

	return &ast.IfStmt{Cond: predicate, Then: consequent, Else: alternative}
}

func (p *parser) parseWhileStmt() ast.Stmt {
	p.expect(token.While)

	return &ast.WhileStmt{Cond: p.parseExpr(), Body: p.parseBlockStmt()}
}

func (p *parser) parseForStmt() ast.Stmt {
	p.expect(token.For)
	element := p.parseIdent(false)
	p.expect(token.In)
	iterator := p.parseExpr()
	body := p.parseBlockStmt()

	return &ast.ForStmt{Element: element, Iterable: iterator, Body: body}
}

func (p *parser) parseSwitchStmt() ast.Stmt {
	p.expect(token.Switch)

	s := new(ast.SwitchStmt)
	s.Key = p.parseExpr()
	p.expect(token.LeftBrace)

	for p.tok.Type != token.EOF && p.tok.Type != token.RightBrace {
		if caseClause, isDefault := p.parseCase(); isDefault {
			s.Default = caseClause
		} else {
			s.Cases = append(s.Cases, caseClause)
		}
	}

	p.expect(token.RightBrace)

	return s
}

func (p *parser) parseCase() (*ast.CaseClause, bool) {
	c := new(ast.CaseClause)

	if p.consume(token.Case) {
		c.Value = p.parseExpr()
		p.expect(token.Colon)
		c.Body = &ast.BlockStmt{Stmts: p.parseStmtList()}
		return c, false
	}

	if p.consume(token.Default) {
		p.expect(token.Colon)
		c.Body = &ast.BlockStmt{Stmts: p.parseStmtList()}
		return c, true
	}

	p.setError(p.tok.Position, "expected case, default or }")
	p.sync(switchStartStmt)
	return nil, false
}

func (p *parser) parseReturnStmt() ast.Stmt {
	retToken := p.expect(token.Return)

	var value ast.Expr
	if p.tok.Type != token.NewLine && p.tok.Type != token.RightBrace {
		value = p.parseExpr()
	}

	return &ast.ReturnStmt{Token: retToken, Value: value}
}

func (p *parser) parseParameterList() (list []*ast.Field) {
	if !p.at(token.LeftParen) {
		return
	}

	p.expect(token.LeftParen)
	for p.tok.Type != token.RightParen && p.tok.Type != token.EOF {
		list = append(list, p.parseField())

		if !p.consumeCommaOrExpect(token.RightParen) {
			return
		}
	}

	p.expect(token.RightParen)

	return
}

func (p *parser) parseSimpleStmt() ast.Stmt {
	lhsExpr := p.parseExprList()

	switch p.tok.Type {
	case token.Assign:
		assign := new(ast.AssignStmt)
		assign.Pos = p.tok.Position
		assign.Left = lhsExpr
		p.expect(token.Assign)
		assign.Right = p.parseExprList()
		return assign
	}

	return &ast.ExprStmt{Expr: lhsExpr[0]}
}

func (p *parser) parseExprList() (list []ast.Expr) {
	list = append(list, p.parseExpr())

	for p.tok.Type == token.Comma {
		p.advance()
		list = append(list, p.parseExpr())
	}

	return
}

func (p *parser) parseExpr() ast.Expr {
	return p.parseBinaryExpr(token.LowestPrecedence)
}

func (p *parser) parseBinaryExpr(precedence int) ast.Expr {
	expr := p.parseUnaryExpr()

	for p.tok.Precedence() > precedence {
		tok := p.expect(p.tok.Type)

		bExpr := new(ast.BinaryExpr)
		bExpr.Left = expr
		bExpr.Operator = tok
		bExpr.Pos = tok.Position
		bExpr.Right = p.parseBinaryExpr(tok.Precedence())
		expr = bExpr
	}

	return expr
}

func (p *parser) parseUnaryExpr() ast.Expr {
	switch p.tok.Type {
	case token.Not, token.Plus, token.Minus:
		uExpr := new(ast.UnaryExpr)
		uExpr.Pos = p.tok.Position
		uExpr.Operator = p.tok
		p.advance()
		uExpr.Expr = p.parseUnaryExpr()
		return uExpr

	default:
		return p.parsePrimaryExpr()
	}
}

func (p *parser) parsePrimaryExpr() (expr ast.Expr) {
	expr = p.parseOperand()

	for {
		switch p.tok.Type {
		case token.Dot:
			expr = p.parseMemberSelector(expr)

		case token.LeftParen:
			expr = p.parseCallExpr(expr)

		case token.LeftBracket:
			expr = p.parseIndexExpr(expr)
		default:
			return
		}
	}
}

func (p *parser) parseOperand() ast.Expr {
	switch p.tok.Type {
	case
		token.Int, token.Float, token.String, token.Bool,
		token.None:
		return p.parseBasicLit()

	case token.Block:
		return p.parseBlockExpr()

	case token.Ident:
		return p.parseIdent(false)

	case token.LeftParen:
		return p.parseGroupExpr()

	case token.LeftBracket:
		return p.parseArrayLit()

	case token.LeftBrace:
		return p.parseHashLit()

	case token.This:
		tok := p.tok
		p.advance()
		this := new(ast.BasicLit)
		this.T = tok.Type
		this.Pos = tok.Position
		return this

	case token.Super:
		return p.parseSuperExpr()

	case token.New:
		return p.parseNewExpr()

	default:
		mesg := fmt.Sprintf("unexpected %s, expecting expression", p.tok)
		p.setError(p.tok.Position, mesg)
		p.sync(closingToken)
		return new(ast.BadExpr)
	}
}

func (p *parser) parseBasicLit() *ast.BasicLit {
	defer p.advance()

	lit := new(ast.BasicLit)
	lit.T = p.tok.Type
	lit.Pos = p.tok.Position

	switch lit.T {
	case token.String:
		lit.Value = readEscape(p.tok.Literal)
	default:
		lit.Value = p.tok.Literal
	}

	return lit
}

func (p *parser) parseBlockExpr() *ast.BlockExpr {
	p.expect(token.Block)

	return &ast.BlockExpr{
		Parameters: p.parseParameterList(),
		Body:       p.parseBlockStmt(),
	}
}

func (p *parser) parseIdent(wantConst bool) *ast.Ident {
	tok := p.expect(token.Ident)

	ident := new(ast.Ident)
	ident.Pos = tok.Position
	ident.Value = tok.Literal
	if wantConst && !ident.IsConstant() {
		p.setError(tok.Position, "expected ident to be a constant")
	}

	return ident
}

func (p *parser) parseGroupExpr() ast.Expr {
	p.expect(token.LeftParen)
	defer p.expect(token.RightParen)

	return &ast.GroupExpr{Expr: p.parseExpr()}
}

func (p *parser) parseMemberSelector(base ast.Expr) ast.Expr {
	p.expect(token.Dot)

	mSel := new(ast.MemberSelector)
	mSel.Base = base
	mSel.Pos = base.Position()
	mSel.Member = p.parseIdent(false)

	return mSel
}

func (p *parser) parseCallExpr(base ast.Expr) ast.Expr {
	switch node := base.(type) {
	case *ast.MemberSelector:
		mCall := new(ast.MethodCallExpr)
		mCall.Pos = node.Position()
		mCall.Selector = node
		mCall.Arguments = p.parseArgumentList()
		return mCall

	case *ast.Ident:
		fCall := new(ast.FunctionCallExpr)
		fCall.Pos = node.Position()
		fCall.Name = node
		fCall.Arguments = p.parseArgumentList()
		return fCall

	default:
		p.setError(p.tok.Position, "invalid call expr")
		p.sync(startStmt)
		return new(ast.BadExpr)
	}
}

func (p *parser) parseArgumentList() (list []ast.Expr) {
	p.expect(token.LeftParen)
	for p.tok.Type != token.RightParen && p.tok.Type != token.EOF {
		list = append(list, p.parseExpr())

		if !p.consumeCommaOrExpect(token.RightParen) {
			return
		}
	}

	p.expect(token.RightParen)
	return
}

func (p *parser) consumeCommaOrExpect(next token.Type) bool {
	if p.consume(token.Comma) || p.tok.Type == next {
		return true
	}

	mesg := fmt.Sprintf("missing , or %s", next)
	p.setError(p.tok.Position, mesg)
	p.sync(closingToken)
	return false
}

func (p *parser) parseArrayLit() (ary *ast.ArrayLit) {
	ary = new(ast.ArrayLit)

	ary.LeftBracket = p.expect(token.LeftBracket)
	for p.tok.Type != token.RightBracket && p.tok.Type != token.EOF {
		ary.Elements = append(ary.Elements, p.parseExpr())

		if !p.consumeCommaOrExpect(token.RightBracket) {
			return
		}
	}

	ary.RightBracket = p.expect(token.RightBracket)
	return
}

func (p *parser) parseHashLit() *ast.HashLit {
	p.expect(token.LeftBrace)
	defer p.expect(token.RightBrace)

	return &ast.HashLit{Entries: p.parseHashEntries()}
}

func (p *parser) parseHashEntries() (list []*ast.HashEntry) {
	for p.tok.Type != token.RightBrace {
		list = append(list, p.parseHashEntry())

		if !p.consumeCommaOrExpect(token.RightBrace) {
			return
		}
	}

	return
}

func (p *parser) parseHashEntry() *ast.HashEntry {
	return &ast.HashEntry{
		Key:   p.parseExpr(),
		Colon: p.expect(token.Colon),
		Value: p.parseExpr(),
	}
}

func (p *parser) parseIndexExpr(expr ast.Expr) ast.Expr {
	p.expect(token.LeftBracket)
	defer p.expect(token.RightBracket)

	return &ast.IndexExpr{Expr: expr, Index: p.parseExpr()}
}

func (p *parser) parseSuperExpr() ast.Expr {
	superTok := p.expect(token.Super)

	expr := new(ast.SuperExpr)
	expr.Pos = superTok.Position
	expr.Arguments = p.parseArgumentList()
	return expr
}

func (p *parser) parseNewExpr() ast.Expr {
	tok := p.expect(token.New)

	expr := new(ast.NewExpr)
	expr.Pos = tok.Position
	expr.Type = p.parseType()
	expr.Arguments = p.parseArgumentList()
	return expr
}

func (p *parser) parseField() *ast.Field {
	field := new(ast.Field)
	field.Name = p.parseIdent(false)
	field.Type = p.parseIdent(true)

	if p.consume(token.Assign) {
		field.Value = p.parseExpr()
	}

	return field
}

func (p *parser) advance() {
	p.tok = p.lexer.NextToken()
}

func (p *parser) at(kind token.Type) bool {
	return p.tok.Type == kind
}

func (p *parser) setError(pos *token.Position, err string) {
	if p.err != nil {
		return
	}

	mesg := fmt.Sprintf("[Lin: %d Col: %d] syntax error: %s", pos.Line(), pos.Column(), err)
	p.err = errors.New(mesg)
}

func (p *parser) expect(expected token.Type) (tok *token.Token) {
	defer p.advance()

	if p.tok.Type != expected {
		p.setError(p.tok.Position, fmt.Sprintf("expected '%s', found '%s'", expected, p.tok.Type))
	}

	return p.tok
}

func (p *parser) consume(tok token.Type) bool {
	if p.tok.Type == tok {
		p.advance()
		return true
	}

	return false
}
