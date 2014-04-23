package fountain

import (
	"github.com/exupero/state-lexer"
)

const (
	TokenDataKey lexer.TokenType = iota
	TokenDataValue

	TokenText
	TokenTextBold
	TokenTextItalic
	TokenTextUnderline
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
		lex.AcceptRun("\n")
		lex.Ignore()
		return lexText
	}
	return lexDataKey
}

func lexText(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
		if r == '*' {
			if lex.Peek() == '*' {
				lex.Backup()
				lex.Emit(TokenText)
				return lexTextBold
			} else {
				lex.Backup()
				lex.Emit(TokenText)
				return lexTextItalic
			}
		}
		if r == '_' {
			lex.Backup()
			lex.Emit(TokenText)
			return lexTextUnderline
		}
	}
	lex.Emit(TokenText)
	return nil
}

func lexTextBold(lex *lexer.Lexer) lexer.StateFn {
	lex.Accept("*")
	lex.Accept("*")
	lex.Ignore()

	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
		if r == '*' && lex.Peek() == '*' {
			lex.Backup()
			break
		}
	}
	lex.Emit(TokenTextBold)

	lex.Accept("*")
	lex.Accept("*")
	lex.Ignore()
	return lexText
}

func lexTextItalic(lex *lexer.Lexer) lexer.StateFn {
	lex.Accept("*")
	lex.Ignore()

	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
		if r == '*' {
			lex.Backup()
			break
		}
	}
	lex.Emit(TokenTextItalic)

	lex.Accept("*")
	lex.Ignore()
	return lexText
}

func lexTextUnderline(lex *lexer.Lexer) lexer.StateFn {
	lex.Accept("_")
	lex.Ignore()

	for {
		r := lex.NextRune()
		if r == -1 {
			break
		}
		if r == '_' {
			lex.Backup()
			break
		}
	}
	lex.Emit(TokenTextUnderline)

	lex.Accept("_")
	lex.Ignore()
	return lexText
}

func Tokenize(src string) *lexer.Lexer {
	lex := lexer.NewLexer(src)
	go lex.Run(lexDataBlock)
	return lex
}
