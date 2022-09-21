package parser

import (
	"fmt"

	"github.com/qiushiyan/peach/pkg/ast"
	"github.com/qiushiyan/peach/pkg/lexer"
	"github.com/qiushiyan/peach/pkg/token"
)

const (
	_ int = iota
	LOWEST
	EQUALS  // == LESSGREATER >, >=, <, <=
	SUM     // +
	PRODUCT // *
	PREFIX  // -X or !X
	CALL    // myFunction(X)
)

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

type Parser struct {
	l         *lexer.Lexer
	curToken  token.Token
	nextToken token.Token
	errors    []string

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
}

func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{l: l, errors: []string{}}

	// Read two tokens, so curToken and nextToken are both set
	p.advancedToken()
	p.advancedToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) advancedToken() {
	p.curToken = p.nextToken
	p.nextToken = p.l.NextToken()
}

func (p *Parser) nextTokenError(t token.TokenType) {
	msg := fmt.Sprintf("expected next token to be %s, got %s instead",
		t, p.nextToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{Statements: []ast.Statement{}}
	for p.curToken.Type != token.EOF {
		statement := p.parseStatement()
		if statement != nil {
			program.Statements = append(program.Statements, statement)
		}
		p.advancedToken()
	}
	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.LET:
		return p.parseLetStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	statement := &ast.LetStatement{Token: p.curToken}

	if !p.expectNextToken(token.IDENTIFIER) {
		return nil
	}

	statement.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectNextToken(token.ASSIGN) {
		return nil
	}

	p.advancedToken()
	for !p.endOfExpression() {
		p.advancedToken()
	}

	return statement
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	statement := &ast.ReturnStatement{Token: p.curToken}
	p.advancedToken()

	for !p.endOfExpression() {
		p.advancedToken()
	}

	return statement
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	statement := &ast.ExpressionStatement{Token: p.curToken}

	p.advancedToken()

	for !p.endOfExpression() {
		p.advancedToken()
	}

	return statement
}

func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}
func (p *Parser) nextTokenIs(t token.TokenType) bool {
	return p.nextToken.Type == t
}

func (p *Parser) expectNextToken(t token.TokenType) bool {
	if p.nextTokenIs(t) {
		p.advancedToken()
		return true
	} else {
		p.nextTokenError(t)
		return false
	}
}

func (p *Parser) endOfExpression() bool {
	return p.curTokenIs(token.SEMICOLON) || p.curTokenIs(token.LINEBREAK) || p.curTokenIs(token.EOF)
}
