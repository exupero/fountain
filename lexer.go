package fountain

import (
	"github.com/exupero/state-lexer"
)

const (
	TokenText lexer.TokenType = iota

	TokenColon
)

func lexText(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
	}
	lex.Emit(TokenText)
	return nil
}

func lexColon(lex *lexer.Lexer) lexer.StateFn {
	lex.Accept(":")
	lex.Emit(TokenColon)

	lex.Accept(" ")
	lex.Ignore()

	return lexText
}

func lexData(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
		if r == ':' {
			lex.Backup()
			break
		}
	}
	lex.Emit(TokenText)
	return lexColon
}

func Tokenize(src string) *lexer.Lexer {
	lex := lexer.NewLexer(src)
	go lex.Run(lexData)
	return lex
}
