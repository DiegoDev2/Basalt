package parser

import (
	"fmt"

	"github.com/DiegoDev2/basalt/pkg/ast"
)

// Parser parses tokens into an AST.
type Parser struct {
	l         *Lexer
	curToken  Token
	peekToken Token
	errors    []string
}

// NewParser creates a new Parser instance.
func NewParser(l *Lexer) *Parser {
	p := &Parser{l: l}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

// ParseFile parses the entire input into a File AST node.
func (p *Parser) ParseFile() (*ast.File, error) {
	file := &ast.File{}

	for p.curToken.Type != EOF {
		switch p.curToken.Literal {
		case "hclconfig":
			config, err := p.parseConfig()
			if err != nil {
				return nil, err
			}
			file.Config = config
		case "table":
			table, err := p.parseTable()
			if err != nil {
				return nil, err
			}
			file.Tables = append(file.Tables, table)
		case "resource":
			resource, err := p.parseResource()
			if err != nil {
				return nil, err
			}
			file.Resources = append(file.Resources, resource)
		default:
			return nil, p.errorf("unexpected top-level identifier: %s", p.curToken.Literal)
		}
		p.nextToken()
	}

	return file, nil
}

func (p *Parser) parseConfig() (*ast.Config, error) {
	config := &ast.Config{}
	// curToken is "hclconfig"

	if !p.expectPeek(LBRACE) {
		return nil, p.errorf("expected '{' after hclconfig")
	}
	// curToken is now "{"
	p.nextToken() // move to first key

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(IDENT) {
			key := p.curToken.Literal
			if !p.expectPeek(COLON) {
				return nil, p.errorf("expected ':' after config key %s", key)
			}
			p.nextToken()
			val := p.curToken.Literal

			switch key {
			case "db":
				config.DB = val
			case "auth":
				config.Auth = val
			case "framework":
				config.Framework = val
			case "lang":
				config.Lang = val
			}
		}
		p.nextToken()
	}

	return config, nil
}

func (p *Parser) parseTable() (*ast.Table, error) {
	table := &ast.Table{}
	p.nextToken() // skip "table"

	if !p.curTokenIs(STRING) {
		return nil, p.errorf("expected table name as string")
	}
	table.Name = p.curToken.Literal

	if !p.expectPeek(LBRACE) {
		return nil, p.errorf("expected '{' after table name")
	}
	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(IDENT) {
			field, err := p.parseField()
			if err != nil {
				return nil, err
			}
			table.Fields = append(table.Fields, field)
		} else {
			p.nextToken()
		}
	}

	return table, nil
}

func (p *Parser) parseField() (*ast.Field, error) {
	field := &ast.Field{Name: p.curToken.Literal}

	if !p.expectPeek(COLON) {
		return nil, p.errorf("expected ':' after field name %s", field.Name)
	}
	p.nextToken() // move to type

	if !p.curTokenIs(IDENT) {
		return nil, p.errorf("expected field type for %s", field.Name)
	}
	field.Type = p.curToken.Literal

	for p.peekTokenIs(AT) {
		p.nextToken() // move to @
		p.nextToken() // move to decorator name
		decorator := &ast.Decorator{Name: p.curToken.Literal}
		if p.peekTokenIs(LPAREN) {
			p.nextToken() // move to (
			p.nextToken() // move to arg
			decorator.Arg = p.curToken.Literal
			if !p.expectPeek(RPAREN) {
				return nil, p.errorf("expected ')' after decorator argument")
			}
		}
		field.Decorators = append(field.Decorators, decorator)
	}

	p.nextToken()
	return field, nil
}

func (p *Parser) parseResource() (*ast.Resource, error) {
	resource := &ast.Resource{}
	p.nextToken() // skip "resource"

	if !p.curTokenIs(STRING) {
		return nil, p.errorf("expected resource name as string")
	}
	resource.Name = p.curToken.Literal

	if !p.expectPeek(LBRACE) {
		return nil, p.errorf("expected '{' after resource name")
	}
	p.nextToken()

	for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
		if p.curTokenIs(IDENT) {
			switch p.curToken.Literal {
			case "table":
				if !p.expectPeek(COLON) {
					return nil, p.errorf("expected ':' after table key")
				}
				p.nextToken()
				resource.Table = p.curToken.Literal
				p.nextToken()
			case "endpoints":
				if !p.expectPeek(LBRACE) {
					return nil, p.errorf("expected '{' after endpoints")
				}
				p.nextToken()
				for !p.curTokenIs(RBRACE) && !p.curTokenIs(EOF) {
					if p.curTokenIs(HTTP_METHOD) {
						method := p.curToken.Literal
						p.nextToken()
						if !p.curTokenIs(PATH) {
							return nil, p.errorf("expected path after HTTP method %s", method)
						}
						resource.Endpoints = append(resource.Endpoints, &ast.Endpoint{
							Method: method,
							Path:   p.curToken.Literal,
						})
					}
					p.nextToken()
				}
				p.nextToken()
			default:
				p.nextToken()
			}
		} else {
			p.nextToken()
		}
	}

	return resource, nil
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	return false
}

func (p *Parser) errorf(format string, a ...interface{}) error {
	msg := fmt.Sprintf(format, a...)
	return fmt.Errorf("Parse error %d:%d: %s", p.curToken.Line, p.curToken.Column, msg)
}
