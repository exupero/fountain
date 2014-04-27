package fountain

import (
	"github.com/exupero/state-lexer"
)

type Parser struct {
	lexer *lexer.Lexer
	Doc *Document
}

type state func(*Parser) state

func Parse(src string) *Document {
	parser := &Parser{
		lexer: Tokenize(src),
		Doc: &Document{
			Data: make(map[string]string),
			Body: []Line{},
		},
	}
	for state := parseDoc; state != nil; {
		state = state(parser)
	}
	return parser.Doc
}

func (p *Parser) Next() (lexer.Token, bool) {
	return p.lexer.Next()
}

func parseDoc(p *Parser) state {
	return parseData
}

func parseData(p *Parser) state {
	for {
		tok, ok := p.Next()
		if !ok || tok.Type != TokenDataKey {
			if tok.Type == TokenParagraph {
				return parseParagraph
			}
		}
		key := tok.Value

		tok, ok = p.Next()
		if !ok || tok.Type != TokenDataValue {
			return nil
		}
		value := tok.Value

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

type styleManager struct {
	bold, italic, underline bool
}

func (s *styleManager) list() []string {
	styles := []string{}
	if s.bold { styles = append(styles, "bold") }
	if s.italic { styles = append(styles, "italic") }
	if s.underline { styles = append(styles, "underline") }
	return styles
}

func parseParagraph(p *Parser) state {
	line := []Text{}
	style := styleManager{false, false, false}

	for {
		tok, ok := p.Next()
		if !ok || tok.Type == TokenParagraph {
			break
		}
		if tok.Type == TokenText {
			line = append(line, Text{content: tok.Value, styles: style.list()})
		}
		if tok.Type == TokenStarDouble {
			style.bold = !style.bold
		}
		if tok.Type == TokenStar {
			style.italic = !style.italic
		}
		if tok.Type == TokenUnderscore {
			style.underline = !style.underline
		}
	}

	p.Doc.Body = append(p.Doc.Body, line)
	return nil
}
