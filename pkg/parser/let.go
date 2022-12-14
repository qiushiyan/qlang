package parser

import (
	"github.com/qiushiyan/qlang/pkg/ast"
	"github.com/qiushiyan/qlang/pkg/token"
)

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.curToken}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.nextTokenIs(token.ASSIGN_ARROW) && !p.nextTokenIs(token.ASSIGN) {
		return nil
	}

	p.advanceToken()
	p.advanceToken()
	// if !p.expectNextToken(token.ASSIGN) {
	// 	return nil
	// }

	// p.advanceToken()

	statement.Value = p.parseExpression(LOWEST)
	if p.endOfExpression() {
		p.advanceToken()
	}

	return statement
}
