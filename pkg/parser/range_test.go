package parser

import (
	"strings"
	"testing"

	"github.com/qiushiyan/qlang/pkg/ast"
	"github.com/qiushiyan/qlang/pkg/lexer"
)

func TestRangeParsing(t *testing.T) {
	input := "1:10"

	l := lexer.New(strings.NewReader(input))
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d", 1, len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	if !testRangeExpression(t, statement.Expression, 1, 10) {
		return
	}
}

func testRangeExpression(t *testing.T, expr ast.Expression, start, end float64) bool {
	rangeExpr, ok := expr.(*ast.RangeExpression)
	if !ok {
		t.Fatalf("expr not *ast.RangeExpression. got=%T", expr)
	}
	exprStart := rangeExpr.Start.(*ast.NumberLiteral).Value
	if exprStart != start {
		t.Errorf("range.Start want=%v, got=%v", exprStart, start)
		return false
	}

	exprEnd := rangeExpr.End.(*ast.NumberLiteral).Value
	if exprEnd != end {
		t.Errorf("range.End want=%v, got=%v", exprEnd, end)
		return false
	}

	return true
}
