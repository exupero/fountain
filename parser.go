package fountain

import (
	"github.com/exupero/state-lexer"
)

type Document struct {
	Title, Credit, Author, DraftDate string
	Data map[string]string
}

type Parser struct {
	lexer *lexer.Lexer
	Doc *Document
	Value string
}

type state func(*Parser) state

func Parse(src string) *Document {
	parser := &Parser{
		lexer: Tokenize(src),
		Doc: &Document{
			Data: make(map[string]string),
		},
	}
	for state := parseDoc; state != nil; {
		state = state(parser)
	}
	return parser.Doc
}

func (p *Parser) Accept(tokenType lexer.TokenType) bool {
	tok, ok := p.lexer.Next()
	if !ok {
		return false
	}
	if tok.Type == tokenType {
		p.Value = tok.Value
		return true
	}
	return false
}

func parseDoc(p *Parser) state {
	return parseData
}

func parseData(p *Parser) state {
	for {
		if !p.Accept(TokenDataKey) {
			break
		}
		key := p.Value

		if !p.Accept(TokenDataValue) {
			break
		}
		value := p.Value

		if key == "Title" {
			p.Doc.Title = value
			continue
		}

		if key == "Credit" {
			p.Doc.Credit = value
			continue
		}

		if key == "Author" {
			p.Doc.Author = value
			continue
		}

		if key == "Draft Date" {
			p.Doc.DraftDate = value
			continue
		}

		p.Doc.Data[key] = value
	}
	return nil
}
