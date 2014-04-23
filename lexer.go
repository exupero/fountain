package fountain

import (
	"github.com/exupero/state-lexer"
)

const (
	TokenDataKey lexer.TokenType = iota
	TokenDataValue
)

func lexDataValue(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
		if r == '\n' {
			lex.Backup()
			break
		}
	}
	lex.Emit(TokenDataValue)

	lex.Accept("\n")
	lex.Ignore()
	return lexDataBlock
}

func lexDataKey(lex *lexer.Lexer) lexer.StateFn {
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
	lex.Emit(TokenDataKey)

	lex.AcceptRun(": ")
	lex.Ignore()
	return lexDataValue
}

func lexDataBlock(lex *lexer.Lexer) lexer.StateFn {
	r := lex.Peek()
	if r == -1 {
		return nil
	}
	if r == '\n' {
		return nil
	}
	return lexDataKey
}

func Tokenize(src string) *lexer.Lexer {
	lex := lexer.NewLexer(src)
	go lex.Run(lexDataBlock)
	return lex
}
