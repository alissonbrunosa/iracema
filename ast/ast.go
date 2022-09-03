package ast

import (
	"iracema/token"
	"strings"
)

type Node interface {
	aNode()
	String() string
	Position() *token.Position
}

type node struct {
	Pos *token.Position
}

func (*node) aNode() {}
func (n *node) Position() *token.Position {
	return n.Pos
}

//
// Statement
//

type Stmt interface {
	Node
	aStmt()
}

type stmt struct{ node }

func (*stmt) aStmt() {}

type File struct {
	Name       string
	Imports    []string
	ObjectList []*ObjectDecl
	VarList    []*VarDecl
	FunList    []*FunDecl

	stmt
}

func (f *File) String() string {
	var buf strings.Builder

	for _, v := range f.VarList {
		buf.WriteString(v.String())
	}

	buf.WriteByte('\n')
	for _, f := range f.FunList {
		buf.WriteString(f.String())
	}

	buf.WriteByte('\n')
	for _, o := range f.ObjectList {
		buf.WriteString(o.String())
	}

	return buf.String()
}

type Import struct {
	Token *token.Token
	Name  string

	stmt
}

func (i *Import) Position() *token.Position { return i.Token.Position }

func (i *Import) String() string { return i.Name }

type BlockStmt struct {
	Stmts      []Stmt
	RightBrace *token.Position

	stmt
}

func (b *BlockStmt) String() string {
	var buf strings.Builder

	for _, s := range b.Stmts {
		buf.WriteString(s.String())
	}

	return buf.String()
}

type StopStmt struct {
	stmt
}

func (*StopStmt) String() string { return "stop" }

type NextStmt struct {
	stmt
}

func (*NextStmt) String() string { return "next" }

type ObjectDecl struct {
	Name      *Ident
	Parent    *Ident
	FieldList []*VarDecl
	FunList   []*FunDecl

	Body *BlockStmt

	stmt
}

func (*ObjectDecl) String() string { return "ObjectDecl" }

type AssignStmt struct {
	Token *token.Token
	Left  []Expr
	Right []Expr

	stmt
}

func (*AssignStmt) String() string { return "AssignStmt" }

type Field struct {
	Name  *Ident
	Type  *Ident
	Value Expr

	node
}

func (f *Field) String() string { return "ast.Field" }

type VarDecl struct {
	Token *token.Token
	Name  *Ident
	Type  *Ident
	Value Expr

	stmt
}

func (*VarDecl) String() string { return "LetDecl" }

type FunDecl struct {
	Name       *Ident
	Parameters []*Field
	Return     *Ident
	Body       *BlockStmt
	Catches    []*CatchDecl

	stmt
}

func (*FunDecl) String() string { return "FunDecl" }

type CatchDecl struct {
	Ref  *Ident
	Type *Ident
	Body *BlockStmt

	stmt
}

func (*CatchDecl) String() string { return "CatchDecl" }

type IfStmt struct {
	Cond Expr
	Then *BlockStmt
	Else Stmt

	stmt
}

func (*IfStmt) String() string { return "IfStmt" }

type ForStmt struct {
	Element  *Ident
	Iterable Expr
	Body     *BlockStmt

	stmt
}

func (*ForStmt) String() string { return "ForStmt" }

type WhileStmt struct {
	Cond Expr
	Body *BlockStmt

	stmt
}

func (*WhileStmt) String() string { return "WhileStmt" }

type CaseClause struct {
	Value Expr
	Body  *BlockStmt

	node
}

func (*CaseClause) String() string { return "CaseClause" }

type SwitchStmt struct {
	Key     Expr
	Cases   []*CaseClause
	Default *CaseClause

	stmt
}

func (*SwitchStmt) String() string { return "SwitchStmt" }

type ExprStmt struct {
	Expr Expr

	stmt
}

func (e *ExprStmt) String() string {
	return e.Expr.String()
}

type ReturnStmt struct {
	Token *token.Token
	Value Expr

	stmt
}

func (r *ReturnStmt) String() string {
	var buf strings.Builder
	buf.WriteString("return ")
	buf.WriteString(r.Value.String())

	return buf.String()
}

//
// Expressions
//

type Expr interface {
	Node
	aExpr()
}

type expr struct{ node }

func (*expr) aExpr() {}

type Ident struct {
	Value string

	expr
}

func (i *Ident) IsConstant() bool { return 'A' <= i.Value[0] && i.Value[0] <= 'Z' }
func (*Ident) String() string     { return "Ident" }

type UnaryExpr struct {
	Operator *token.Token
	Expr     Expr

	expr
}

func (u *UnaryExpr) String() string {
	var buf strings.Builder
	buf.WriteByte('(')
	buf.WriteString(u.Operator.String())
	buf.WriteString(u.Expr.String())
	buf.WriteByte(')')

	return buf.String()
}

type BinaryExpr struct {
	Operator *token.Token
	Left     Expr
	Right    Expr

	expr
}

func (b *BinaryExpr) String() string {
	var buf strings.Builder
	buf.WriteByte('(')
	buf.WriteString(b.Left.String())
	buf.WriteString(b.Operator.String())
	buf.WriteString(b.Right.String())
	buf.WriteByte(')')

	return buf.String()
}

type BasicLit struct {
	T     token.Type
	Value string

	expr
}

func (b *BasicLit) String() string   { return b.Value }
func (b *BasicLit) Type() token.Type { return b.T }

type ArrayLit struct {
	LeftBracket  *token.Token
	Elements     []Expr
	RightBracket *token.Token

	expr
}

func (*ArrayLit) String() string { return "ArrayLit" }

type HashEntry struct {
	Key   Expr
	Colon *token.Token
	Value Expr

	node
}

func (*HashEntry) String() string { return "KeyValueExpr" }

type HashLit struct {
	Entries []*HashEntry

	expr
}

func (*HashLit) String() string { return "HashLit" }

type IndexExpr struct {
	Expr  Expr
	Index Expr

	expr
}

func (*IndexExpr) String() string { return "IndexExpr" }

type BlockExpr struct {
	Parameters []*Field
	Body       *BlockStmt

	expr
}

func (*BlockExpr) String() string { return "CodeBlock" }

type MethodCallExpr struct {
	Selector  *MemberSelector
	Arguments []Expr

	expr
}

func (*MethodCallExpr) String() string { return "MethodCallExpr" }

type FunctionCallExpr struct {
	Name      *Ident
	Arguments []Expr

	expr
}

func (*FunctionCallExpr) String() string { return "FunctionCallExpr" }

type BadExpr struct {
	Expr Expr

	expr
}

func (*BadExpr) String() string { return "BadExpr" }

type GroupExpr struct {
	Expr Expr

	expr
}

func (g *GroupExpr) String() string { return g.Expr.String() }

type SuperExpr struct {
	Arguments    []Expr
	ExplicitArgs bool

	expr
}

func (*SuperExpr) String() string { return "SuperExpr" }

type MemberSelector struct {
	Base   Expr
	Member *Ident

	expr
}

func (*MemberSelector) String() string { return "FieldSel" }

type NewExpr struct {
	Type      *Ident
	Arguments []Expr

	expr
}

func (*NewExpr) String() string { return "NewExpr" }
