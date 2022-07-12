package parser

import (
	"fmt"
	"io"
	"iracema/ast"
	"iracema/lexer"
	"iracema/token"
	"strings"
)

type parser struct {
	lexer  lexer.Lexer
	tok    *token.Token
	errors ErrorList
}

func (p *parser) init(source io.Reader) {
	p.lexer = lexer.New(source, p.addError)
	p.next()
}

func (p *parser) parse() *ast.File {
	return &ast.File{
		Stmts: p.parseStmtList(),
	}
}

func isLit(stmt ast.Stmt) bool {
	exprStmt, ok := stmt.(*ast.ExprStmt)
	if !ok {
		return false
	}

	_, ok = exprStmt.Expr.(*ast.BasicLit)
	return ok
}

func (p *parser) parseStmtList() (list []ast.Stmt) {
	for p.tok.Type != token.Eof && p.tok.Type != token.RightBrace {
		if len(list) != 0 && isLit(list[0]) {
			list[0] = p.parseStmt()
			continue
		}

		list = append(list, p.parseStmt())
	}

	return
}

var startStmt = map[token.Type]bool{
	token.Object: true,
	token.Fun:    true,
	token.If:     true,
	token.While:  true,
	token.Stop:   true,
	token.Return: true,
}

func (p *parser) advance(to map[token.Type]bool) {
	for ; !p.at(token.Eof); p.next() {
		if to[p.tok.Type] {
			break
		}
	}
}

func (p *parser) parseStmt() ast.Stmt {
	switch p.tok.Type {
	case token.Object:
		return p.parseObjectDecl()

	case token.Fun:
		return p.parseFunDecl()

	case token.If:
		return p.parseIfStmt()

	case token.While:
		return p.parseWhileStmt()

	case token.For:
		return p.parseForStmt()

	case token.Stop:
		return p.parseStopStmt()

	case token.Next:
		return p.parseNextStmt()

	case token.Return:
		return p.parseReturnStmt()

	case
		token.Ident, token.String, token.Bool, token.Int, token.Float,
		token.LeftParenthesis, token.Not, token.Plus, token.Minus,
		token.LeftBracket, token.LeftBrace, token.Nil, token.Block:
		return p.parseSimpleStmt()

	default:
		p.addError(p.tok.Position, "unknown token")
		p.advance(startStmt)
		return nil
	}
}

func (p *parser) parseObjectDecl() ast.Stmt {
	p.expect(token.Object)

	name := p.parseConst()

	var parent *ast.Ident
	if p.consume(token.Is) {
		parent = p.parseConst()
	}

	body := p.parseBlockStmt()

	return &ast.ObjectDecl{Name: name, Parent: parent, Body: body}
}

func (p *parser) parseFunDecl() ast.Stmt {
	p.expect(token.Fun)

	return &ast.FunDecl{
		Name:       p.parseIdent(),
		Parameters: p.parseParameterList(),
		Body:       p.parseBlockStmt(),
		Catches:    p.parseCatchList(),
	}
}

func (p *parser) parseCatchList() (list []*ast.CatchDecl) {
	for p.consume(token.Catch) {
		p.expect(token.LeftParenthesis)
		ref := p.parseIdent()
		p.expect(token.Colon)
		typ := p.parseConst()
		p.expect(token.RightParenthesis)
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
	predicate := p.parseExpr(true)
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
			p.addError(p.tok.Position, "syntax error: expected left brace or if statement")
			p.next()
		}
	}

	return &ast.IfStmt{Cond: predicate, Then: consequent, Else: alternative}
}

func (p *parser) parseWhileStmt() ast.Stmt {
	p.expect(token.While)

	return &ast.WhileStmt{Cond: p.parseExpr(true), Body: p.parseBlockStmt()}
}

func (p *parser) parseForStmt() ast.Stmt {
	p.expect(token.For)
	element := p.parseIdent()
	p.expect(token.In)
	iterable := p.parseExpr(false)
	body := p.parseBlockStmt()

	return &ast.ForStmt{Element: element, Iterable: iterable, Body: body}
}

func (p *parser) parseStopStmt() ast.Stmt {
	return &ast.StopStmt{
		Token: p.expect(token.Stop),
	}
}

func (p *parser) parseNextStmt() ast.Stmt {
	return &ast.NextStmt{
		Token: p.expect(token.Next),
	}
}

func (p *parser) parseReturnStmt() ast.Stmt {
	return &ast.ReturnStmt{
		Token: p.expect(token.Return),
		Expr:  p.parseExpr(false),
	}
}

func (p *parser) parseParameterList() (list []*ast.Ident) {
	if !p.at(token.LeftParenthesis) {
		return
	}

	p.expect(token.LeftParenthesis)
	for p.tok.Type != token.RightParenthesis && p.tok.Type != token.Eof {
		param := p.parseIdent()
		if param.IsAttr() {
			p.addError(param.Token.Position, "syntax error: argument cannot be an instance variable")
			continue
		}

		list = append(list, param)
		if !p.atComma(token.RightParenthesis) {
			break
		}

		p.next()
	}

	p.expect(token.RightParenthesis)

	return
}

func (p *parser) parseSimpleStmt() ast.Stmt {
	leftExpr := p.parseLeftExprList()

	switch p.tok.Type {
	case token.Assign:
		tok := p.tok
		p.next()
		rightExpr := p.parseRightExprList()
		return &ast.AssignStmt{Left: leftExpr, Token: tok, Right: rightExpr}
	}

	return &ast.ExprStmt{Expr: leftExpr[0]}
}

func (p *parser) parseLeftExprList() []ast.Expr {
	return p.parseExprList(true)
}

func (p *parser) parseRightExprList() []ast.Expr {
	return p.parseExprList(false)
}

func (p *parser) parseExprList(leftHand bool) (list []ast.Expr) {
	list = append(list, p.parseExpr(leftHand))

	for p.tok.Type == token.Comma {
		p.next()
		list = append(list, p.parseExpr(leftHand))
	}

	return
}

func (p *parser) parseExpr(leftHand bool) ast.Expr {
	return p.parseBinaryExpr(leftHand, token.LowestPrecedence+1)
}

func (p *parser) parseBinaryExpr(leftHand bool, precedence int) ast.Expr {
	left := p.parseUnaryExpr(leftHand)

	for {
		if p.tok.Precedence() < precedence {
			return left
		}

		tok := p.expect(p.tok.Type)
		right := p.parseBinaryExpr(false, tok.Precedence()+1)

		left = &ast.BinaryExpr{Left: left, Operator: tok, Right: right}
	}
}

func (p *parser) parseUnaryExpr(leftHand bool) (expr ast.Expr) {
	switch p.tok.Type {
	case token.Not, token.Plus, token.Minus:
		expr = &ast.UnaryExpr{
			Operator: p.expect(p.tok.Type),
			Expr:     p.parseUnaryExpr(leftHand),
		}
	default:
		expr = p.parsePrimaryExpr(leftHand)
	}

	return
}

func (p *parser) parsePrimaryExpr(leftHand bool) (expr ast.Expr) {
	expr = p.parseOperand(leftHand)

	for {
		switch p.tok.Type {
		case token.Dot:
			expr = p.parseCallExpr(expr)

		case token.LeftParenthesis:
			expr = &ast.CallExpr{
				Method:    expr.(*ast.Ident),
				Arguments: p.parseArgumentList(),
			}

		case token.LeftBracket:
			expr = p.parseIndexExpr(expr)
		default:
			return
		}
	}
}

func (p *parser) parseOperand(leftHand bool) (expr ast.Expr) {
	switch p.tok.Type {
	case
		token.Int, token.Float, token.String, token.Bool,
		token.Nil:
		expr = &ast.BasicLit{Token: p.tok, Value: p.tok.Literal}
		p.next()

	case token.Block:
		expr = p.parseBlockExpr()

	case token.Ident:
		expr = p.parseIdent()

	case token.LeftParenthesis:
		expr = p.parseGroupExpr()

	case token.LeftBracket:
		expr = p.parseArrayLit()

	case token.LeftBrace:
		expr = p.parseHashLit()

	default:
		p.addError(p.tok.Position, fmt.Sprintf("no parse implemented for (%q) just yet\n", p.tok.Type))
		p.next()
	}

	return
}

func (p *parser) parseBlockExpr() *ast.BlockExpr {
	p.expect(token.Block)

	return &ast.BlockExpr{
		Parameters: p.parseParameterList(),
		Body:       p.parseBlockStmt(),
	}
}

func (p *parser) parseIdent() *ast.Ident {
	tok := p.expect(token.Ident)

	return &ast.Ident{Token: tok, Value: tok.Literal}
}

func (p *parser) parseConst() *ast.Ident {
	tok := p.expect(token.Ident)

	ident := &ast.Ident{Value: tok.Literal}
	if !ident.IsConstant() {
		p.addError(tok.Position, "syntax error: expected ident to be a constant")
	}

	return ident
}

func (p *parser) parseGroupExpr() ast.Expr {
	p.expect(token.LeftParenthesis)
	expr := p.parseExpr(false)
	p.expect(token.RightParenthesis)

	return &ast.GroupExpr{Expr: expr}
}

func (p *parser) parseCallExpr(receiver ast.Expr) ast.Expr {
	p.expect(token.Dot)

	return &ast.CallExpr{
		Receiver:  receiver,
		Method:    p.parseIdent(),
		Arguments: p.parseArgumentList(),
	}
}

func (p *parser) parseArgumentList() (list []ast.Expr) {
	if !p.at(token.LeftParenthesis) {
		return
	}

	p.expect(token.LeftParenthesis)
	for p.tok.Type != token.RightParenthesis && p.tok.Type != token.Eof {
		list = append(list, p.parseExpr(false))
		if !p.atComma(token.RightParenthesis) {
			break
		}

		p.next()
	}

	p.expect(token.RightParenthesis)

	return
}

func (p *parser) parseArrayLit() *ast.ArrayLit {
	leftBracket := p.expect(token.LeftBracket)

	var list []ast.Expr
	for p.tok.Type != token.RightBracket && p.tok.Type != token.Eof {
		list = append(list, p.parseExpr(false))
		if !p.atComma(token.RightBracket) {
			break
		}

		p.next()
	}

	rightBracket := p.expect(token.RightBracket)

	return &ast.ArrayLit{LeftBracket: leftBracket, Elements: list, RightBracket: rightBracket}
}

func (p *parser) parseHashLit() *ast.HashLit {
	return &ast.HashLit{
		LeftBrace:  p.expect(token.LeftBrace),
		Elements:   p.KeyValuePairList(),
		RightBrace: p.expect(token.RightBrace),
	}
}

func (p *parser) KeyValuePairList() (list []*ast.KeyValueExpr) {
	for p.tok.Type != token.RightBrace {

		list = append(list, p.KeyValuePair())
		if !p.atComma(token.RightBrace) {
			break
		}

		p.next()
	}

	return
}

func (p *parser) KeyValuePair() *ast.KeyValueExpr {
	return &ast.KeyValueExpr{
		Key:   p.parseExpr(false),
		Colon: p.expect(token.Colon),
		Value: p.parseExpr(false),
	}
}

func (p *parser) parseIndexExpr(expr ast.Expr) ast.Expr {
	return &ast.IndexExpr{
		Expr:         expr,
		LeftBracket:  p.expect(token.LeftBracket),
		Index:        p.parseExpr(false),
		RightBracket: p.expect(token.RightBracket),
	}
}

func (p *parser) next() {
	p.tok = p.lexer.NextToken()
}

func (p *parser) at(kind token.Type) bool {
	return p.tok.Type == kind
}

func (p *parser) addError(pos *token.Position, err string) {
	var b strings.Builder
	fmt.Fprintf(&b, "[Lin: %d Col: %d] ", pos.Line(), pos.Column())
	b.WriteString(err)

	p.errors = append(p.errors, &Error{Msg: b.String()})
}

func (p *parser) expect(expected token.Type) (tok *token.Token) {
	if p.tok.Type != expected {
		p.addError(p.tok.Position, fmt.Sprintf("syntax error: expected '%s', found '%s'", expected, p.tok.Type))
	}

	tok = p.tok
	p.next()
	return
}

func (p *parser) atComma(next token.Type) bool {
	if p.tok.Is(token.Comma) {
		return true
	}

	if p.tok.Is(next) {
		return false
	}

	p.addError(p.tok.Position, "syntax error: missing ','")
	p.next() // cosume invalid token
	return false
}

func (p *parser) consume(tok token.Type) bool {
	if p.tok.Type == tok {
		p.next()
		return true
	}

	return false
}
