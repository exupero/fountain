package fountain

import (
	"strings"

	"github.com/exupero/state-lexer"
)

const (
	TokenDataKey lexer.TokenType = iota
	TokenDataValue

	TokenText
	TokenStar
	TokenStarDouble
	TokenUnderscore

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
		lex.Emit(TokenDialogue)
		return nil
	}

	if r == '(' {
		lex.Backup()
		return lexParenthetical
	}

	if r == '\n' {
		lex.Ignore()
		return lexBody
	}

	lex.Backup()
	return lexDialogueText
}

func lexDialogueText(lex *lexer.Lexer) lexer.StateFn {
	for {
		r := lex.NextRune()

		if r == lexer.Eof {
			lex.Emit(TokenDialogue)
			return nil
		}

		if r == '\n' {
			lex.Backup()
			lex.Emit(TokenDialogue)
			lex.Accept("\n")
			lex.Ignore()
			return lexDialogue
		}

		if r == '*' {
			lex.Backup()
			lex.Emit(TokenDialogue)
			lex.Accept("*")

			r = lex.NextRune()
			if r == '*' {
				lex.Emit(TokenStarDouble)
				return lexDialogueText
			}
			lex.Backup()
			lex.Emit(TokenStar)
			return lexDialogueText
		}

		if r == '_' {
			lex.Backup()
			lex.Emit(TokenDialogue)
			lex.Accept("_")
			lex.Emit(TokenUnderscore)
			return lexDialogueText
		}
	}
	return nil
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
			lex.Emit(TokenText)
			return nil
		}

		if r == '\n' {
			lex.Backup()
			lex.Emit(TokenText)
			lex.Accept("\n")
			lex.Ignore()

			r = lex.NextRune()
			if r == '\n' {
				lex.Ignore()
			} else {
				lex.Backup()
			}
			return lexBody
		}

		if r == '*' {
			lex.Backup()
			lex.Emit(TokenText)
			lex.Accept("*")

			r = lex.NextRune()
			if r == '*' {
				lex.Emit(TokenStarDouble)
				return lexText
			}
			lex.Backup()
			lex.Emit(TokenStar)
			return lexText
		}

		if r == '_' {
			lex.Backup()
			lex.Emit(TokenText)
			lex.Accept("_")
			lex.Emit(TokenUnderscore)
			return lexText
		}
	}
	return nil
}

func Tokenize(src string) *lexer.Lexer {
	lex := lexer.NewLexer(src)
	go lex.Run(lexDataBlock)
	return lex
}
