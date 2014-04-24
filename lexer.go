package fountain

import (
	"strings"

	"github.com/exupero/state-lexer"
)

const (
	TokenDataKey lexer.TokenType = iota
	TokenDataValue

	TokenText
	TokenTextBold
	TokenTextItalic
	TokenTextUnderline

	TokenSpeaker
	TokenDialogue
	TokenParenthetical
)

func lexDataValue(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == lexer.Eof {
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
		if r == lexer.Eof {
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
	if r == lexer.Eof {
		return nil
	}
	if r == '\n' {
		lex.Accept("\n")
		lex.Ignore()
		return lexBody
	}
	return lexDataKey
}

func lexBody(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == lexer.Eof {
			break
		}
		if strings.IndexRune("abcdefghijklmnopqrstuvwxyz*_", r) >= 0 {
			lex.Backup()
			return lexText
		}
		if strings.IndexRune("\n", r) >= 0 {
			lex.Backup()
			return lexSpeaker
		}
	}
	return nil
}

func lexSpeaker(lex *lexer.Lexer) lexer.StateFn {
	lex.Emit(TokenSpeaker)
	lex.Accept("\n")
	lex.Ignore()
	return lexDialogue
}

func lexDialogue(lex *lexer.Lexer) lexer.StateFn {
	r := lex.NextRune()

	if r == lexer.Eof {
		return nil
	}

	if r == '(' {
		lex.Backup()
		return lexParenthetical
	}

	for {
		r = lex.NextRune()

		if r == lexer.Eof {
			lex.Emit(TokenDialogue)
			break
		}
		if r == '\n' {
			lex.Backup()
			return lexDialogueText
		}
	}
	return nil
}

func lexDialogueText(lex *lexer.Lexer) lexer.StateFn {
	lex.Emit(TokenDialogue)
	lex.Accept("\n")
	lex.Ignore()
	return lexDialogue
}

func lexParenthetical(lex *lexer.Lexer) lexer.StateFn {
	lex.Accept("(")
	lex.Ignore()

	lex.Until(")")
	lex.Emit(TokenParenthetical)

	lex.Accept(")")
	lex.Accept("\n")
	lex.Ignore()

	return lexDialogue
}

func lexText(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()
		if r == lexer.Eof {
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
		if r == lexer.Eof {
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
		if r == lexer.Eof {
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
		if r == lexer.Eof {
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
