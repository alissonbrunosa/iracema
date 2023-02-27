package parser

import (
	"errors"
	"fmt"
	"io"
	"iracema/ast"
	"iracema/lexer"
	"iracema/token"
)

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
	var stmts []ast.Stmt

	var imports []string
	for p.consume(token.Use) {
		name := p.expect(token.String)
		imports = append(imports, name.Literal)

		p.consume(token.NewLine)
	}

	for !isDone(p.tok) {
		stmts = append(stmts, p.parseStmt())

		if !p.consume(token.NewLine) && p.tok.Type != token.EOF {
			err := fmt.Sprintf("unexpected %s, expecting EOF or new line", p.tok)
			p.setError(p.tok.Position, err)
			p.sync(startStmt)
			continue
		}
	}

	return &ast.File{Stmts: stmts, Imports: imports}
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

	case token.Const:
		return p.parseConstDecl()

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
		return p.parseStopStmt()

	case token.Next:
		return p.parseNextStmt()

	case token.Return:
		return p.parseReturnStmt()

	case
		token.Ident, token.String, token.Bool, token.Int, token.Float,
		token.LeftParen, token.Not, token.Plus, token.Minus,
		token.LeftBracket, token.LeftBrace, token.None,
		token.Super, token.This:
		return p.parseSimpleStmt()

	default:
		return nil
	}
}

func (p *parser) parseObjectDecl() ast.Stmt {
	p.expect(token.Object)

	obj := new(ast.ObjectDecl)
	obj.Name = p.parseConst()

	obj.TypeParamList = p.parseParamTypeList()

	if p.consume(token.Is) {
		obj.Parent = p.parseType()
	}

	p.expect(token.LeftBrace)
	for p.tok.Type != token.RightBrace {
		switch p.tok.Type {
		case token.Var:
			obj.FieldList = append(obj.FieldList, p.parseVarDecl())

		case token.Const:
			obj.ConstantList = append(obj.ConstantList, p.parseConstDecl())

		case token.Fun:
			obj.FunctionList = append(obj.FunctionList, p.parseFunDecl())
		case token.NewLine:
			p.advance()
			continue

		default:
			mesg := fmt.Sprintf("unexpected %s, expecting FunDecl or VarDecl", p.tok)
			p.setError(p.tok.Position, mesg)
			return obj
		}
	}

	p.expect(token.RightBrace)

	return obj
}

func (p *parser) parseParamTypeList() (list []*ast.TypeParam) {
	if !p.consume(token.Less) {
		return
	}

	list = append(list, p.parseParamType())
	for p.consume(token.Comma) {
		list = append(list, p.parseParamType())
	}

	p.expect(token.Great)
	return
}

func (p *parser) parseParamType() *ast.TypeParam {
	field := new(ast.TypeParam)
	field.Name = p.parseIdent()

	if p.consume(token.Is) {
		field.Type = p.parseConst()
	}

	return field
}

func (p *parser) parseVarDecl() *ast.VarDecl {
	p.expect(token.Var)

	decl := new(ast.VarDecl)
	decl.Name = p.parseIdent()

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

func (p *parser) parseConstDecl() *ast.ConstDecl {
	p.expect(token.Const)

	decl := new(ast.ConstDecl)
	decl.Name = p.parseIdent()

	if p.consume(token.Assign) {
		decl.Value = p.parseExpr()
	} else {
		decl.Type = p.parseType()
		p.expect(token.Assign)
		decl.Value = p.parseExpr()
	}

	return decl
}

func (p *parser) parseType() ast.Type {
	switch p.tok.Type {
	case token.Ident:
		name := p.parseConst()
		if p.at(token.Less) {
			return p.parseParameterizedType(name)
		}

		return name
	case token.Fun:
		return p.parseFunctionType(false, false)

	default:
		err := fmt.Sprintf("expected type, found %s", p.tok)
		p.setError(p.tok.Position, err)
		p.advance()
	}

	return nil
}

func (p *parser) parseParameterizedType(name *ast.Ident) ast.Type {
	p.expect(token.Less)

	t := new(ast.ParameterizedType)
	t.Name = name
	t.TypeArguments = append(t.TypeArguments, p.parseType())

	for p.consume(token.Comma) {
		t.TypeArguments = append(t.TypeArguments, p.parseType())
	}
	p.expect(token.Great)

	return t
}

func (p *parser) parseFunctionType(wantName, wantParamNames bool) *ast.FunctionType {
	p.expect(token.Fun)

	sig := new(ast.FunctionType)
	if wantName {
		sig.Name = p.parseIdent()
	}

	sig.ParameterList = p.parseParameterList(wantParamNames)
	if p.consume(token.Arrow) {
		sig.Return = p.parseType()
	}

	return sig
}

func (p *parser) parseFunDecl() *ast.FunDecl {
	fun := new(ast.FunDecl)
	fun.Type = p.parseFunctionType(true, true)
	fun.Body = p.parseBlockStmt()
	fun.Catches = p.parseCatchList()
	return fun
}

func (p *parser) parseCatchList() (list []*ast.CatchDecl) {
	for p.consume(token.Catch) {
		p.expect(token.LeftParen)
		ref := p.parseIdent()
		p.expect(token.Colon)
		typ := p.parseConst()
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
	p.expect(token.LeftBrace)
	stmts := p.parseStmtList()
	p.expect(token.RightBrace)

	return &ast.BlockStmt{Stmts: stmts}
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
	element := p.parseIdent()
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

func (p *parser) parseStopStmt() ast.Stmt {
	return &ast.StopStmt{Token: p.expect(token.Stop)}
}

func (p *parser) parseNextStmt() ast.Stmt {
	return &ast.NextStmt{Token: p.expect(token.Next)}
}

func (p *parser) parseReturnStmt() ast.Stmt {
	retToken := p.expect(token.Return)

	var value ast.Expr
	if p.tok.Type != token.NewLine && p.tok.Type != token.RightBrace {
		value = p.parseExpr()
	}

	return &ast.ReturnStmt{Token: retToken, Value: value}
}

func (p *parser) parseParameterList(wantNames bool) (list []*ast.VarDecl) {
	p.expect(token.LeftParen)

	for p.tok.Type != token.RightParen && p.tok.Type != token.EOF {
		list = append(list, p.parseParameter(wantNames))

		if !p.consumeCommaOrExpect(token.RightParen) {
			return
		}
	}

	p.expect(token.RightParen)

	return
}

func (p *parser) parseParameter(wantName bool) *ast.VarDecl {
	decl := new(ast.VarDecl)

	if wantName {
		decl.Name = p.parseIdent()
	}

	decl.Type = p.parseType()
	if p.consume(token.Assign) {
		decl.Value = p.parseExpr()
	}

	return decl
}

func (p *parser) parseSimpleStmt() ast.Stmt {
	leftExpr := p.parseExprList()

	switch p.tok.Type {
	case token.Assign:
		return &ast.AssignStmt{
			Left:  leftExpr,
			Token: p.expect(token.Assign),
			Right: p.parseExprList(),
		}
	}

	return &ast.ExprStmt{Expr: leftExpr[0]}
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
	left := p.parseUnaryExpr()

	for p.tok.Precedence() > precedence {
		tok := p.expect(p.tok.Type)
		right := p.parseBinaryExpr(tok.Precedence())

		left = &ast.BinaryExpr{Left: left, Operator: tok, Right: right}
	}

	return left
}

func (p *parser) parseUnaryExpr() (expr ast.Expr) {
	switch p.tok.Type {
	case token.Not, token.Plus, token.Minus:
		expr = &ast.UnaryExpr{
			Operator: p.expect(p.tok.Type),
			Expr:     p.parseUnaryExpr(),
		}
	default:
		expr = p.parsePrimaryExpr()
	}

	return
}

func (p *parser) parsePrimaryExpr() (expr ast.Expr) {
	expr = p.parseOperand()

	for {
		switch p.tok.Type {
		case token.Dot:
			member := new(ast.MemberExpr)
			member.Base = expr
			p.expect(token.Dot)
			member.Name = p.parseIdent()
			expr = member

		case token.LeftParen:
			expr = &ast.CallExpr{Function: expr, Arguments: p.parseArgumentList()}

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

	case token.Fun:
		return p.parseFunLiteral()

	case token.Ident:
		return p.parseIdent()

	case token.LeftParen:
		return p.parseGroupExpr()

	case token.LeftBracket:
		return p.parseArrayLit()

	case token.LeftBrace:
		return p.parseHashLit()

	case token.This:
		this := new(ast.BasicLit)
		this.Token = p.tok
		p.advance()

		return this

	case token.Super:
		return p.parseSuperExpr()

	default:
		mesg := fmt.Sprintf("unexpected %s, expecting expression", p.tok)
		p.setError(p.tok.Position, mesg)
		p.sync(closingToken)
		return new(ast.BadExpr)
	}
}

func (p *parser) parseBasicLit() (lit *ast.BasicLit) {
	defer p.advance()

	switch p.tok.Type {
	case token.String:
		return &ast.BasicLit{Token: p.tok, Value: readEscape(p.tok.Literal)}
	default:
		return &ast.BasicLit{Token: p.tok, Value: p.tok.Literal}
	}
}

func (p *parser) parseFunLiteral() (fun *ast.FunLiteral) {
	fun = new(ast.FunLiteral)
	fun.Type = p.parseFunctionType(false, true)
	fun.Body = p.parseBlockStmt()
	return
}

func (p *parser) parseIdent() *ast.Ident {
	tok := p.expect(token.Ident)

	return &ast.Ident{Token: tok, Value: tok.Literal}
}

func (p *parser) parseConst() *ast.Ident {
	tok := p.expect(token.Ident)

	ident := &ast.Ident{Value: tok.Literal}
	if !ident.IsConstant() {
		p.setError(tok.Position, "expected ident to be a constant")
	}

	return ident
}

func (p *parser) parseGroupExpr() ast.Expr {
	p.expect(token.LeftParen)
	expr := p.parseExpr()
	p.expect(token.RightParen)

	return &ast.GroupExpr{Expr: expr}
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

func (p *parser) parseHashLit() *ast.MapLit {
	return &ast.MapLit{
		LeftBrace:  p.expect(token.LeftBrace),
		Entries:    p.parseHashEntries(),
		RightBrace: p.expect(token.RightBrace),
	}
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
	return &ast.IndexExpr{
		Expr:         expr,
		LeftBracket:  p.expect(token.LeftBracket),
		Index:        p.parseExpr(),
		RightBracket: p.expect(token.RightBracket),
	}
}

func (p *parser) parseSuperExpr() ast.Expr {
	return &ast.SuperExpr{
		Token:        p.expect(token.Super),
		ExplicitArgs: p.at(token.LeftParen),
		Arguments:    p.parseArgumentList(),
	}
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
