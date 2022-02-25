package token

type Position struct {
	line   int
	column int
}

func (p *Position) Line() int   { return p.line }
func (p *Position) Column() int { return p.column }

func (p *Position) AddLine(offset int) {
	p.line++
	p.column = offset
}

func (p *Position) Snapshot(offset int) *Position {
	column := offset - p.column
	if column == 0 {
		column = 1
	}

	return &Position{line: p.line, column: column}
}

func NewPosition() *Position {
	return &Position{line: 1}
}
