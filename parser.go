package fountain

import (
	"github.com/exupero/state-lexer"
)

type Parser struct {
	lexer *lexer.Lexer
	Doc *Document
	hasNext bool
	next lexer.Token
}

type state func(*Parser) state

func Parse(src string) *Document {
	parser := &Parser{
		lexer: Tokenize(src),
		Doc: &Document{
			Data: make(map[string]string),
			Body: []Paragraph{},
		},
	}
	for state := parseDoc; state != nil; {
		state = state(parser)
	}
	return parser.Doc
}

func (p *Parser) Next() (lexer.Token, bool) {
	if p.hasNext {
		p.hasNext = false
		return p.next, true
	}
	return p.lexer.Next()
}

func (p *Parser) Peek() (lexer.Token, bool) {
	tok, ok := p.lexer.Next()
	if ok {
		p.hasNext = true
		p.next = tok
	}
	return tok, ok
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

func parseParagraph(p *Parser) state {
	tok, ok := p.Peek()
	if !ok {
		return nil
	}
	if tok.Type == TokenSpeaker {
		return parseDialogue
	}
	return parseAction
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

func parseAction(p *Parser) state {
	style := styleManager{false, false, false}
	lines := []Line{}
	chunks := []Chunk{}

	defer func() {
		lines = append(lines, Line{Chunks: chunks, Type: "action"})
		paragraph := Paragraph{Lines: lines, Type: "action"}
		p.Doc.Body = append(p.Doc.Body, paragraph)
	}()

	for {
		tok, ok := p.Next()
		if !ok {
			return nil
		}
		if tok.Type == TokenParagraph {
			if tok.Value == "\n\n" {
				return parseParagraph
			}
			lines = append(lines, Line{Chunks: chunks, Type: "action"})
			chunks = []Chunk{}
		}
		if tok.Type == TokenText {
			chunks = append(chunks, Chunk{Content: tok.Value, Styles: style.list()})
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
	return nil
}

func parseDialogue(p *Parser) state {
	lines := []Line{}

	defer func() {
		paragraph := Paragraph{
			Lines: lines,
			Type: "dialogue",
		}
		p.Doc.Body = append(p.Doc.Body, paragraph)
	}()

	for {
		tok, ok := p.Next()
		if !ok {
			return nil
		}
		if tok.Type == TokenParagraph {
			return parseParagraph
		}

		if tok.Type == TokenSpeaker {
			line := Line{
				Chunks: []Chunk{
					Chunk{Content: tok.Value},
				},
				Type: "speaker",
			}
			lines = append(lines, line)
		}
		if tok.Type == TokenParenthetical {
			line := Line{
				Chunks: []Chunk{
					Chunk{Content: tok.Value},
				},
				Type: "parenthetical",
			}
			lines = append(lines, line)
		}
		if tok.Type == TokenDialogue {
			line := parseDialogueText(p, tok)
			lines = append(lines, line)
		}
	}
	return nil
}

func parseDialogueText(p *Parser, tok lexer.Token) Line {
	style := styleManager{false, false, false,}
	chunks := []Chunk{Chunk{Content: tok.Value}}

	for {
		// Check before consuming.
		tok, ok := p.Peek()
		if !ok || tok.Type == TokenParagraph || tok.Type == TokenSpeaker || tok.Type == TokenParenthetical {
			break
		}

		tok, _ = p.Next()
		if tok.Type == TokenDialogue {
			chunks = append(chunks, Chunk{Content: tok.Value, Styles: style.list()})
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

	return Line{
		Chunks: chunks,
		Type: "dialogue",
	}
}
