package fountain

import (
	"strings"

	"github.com/exupero/state-lexer"
)

const (
	TokenDataKey lexer.TokenType = iota
	TokenDataValue

	TokenParagraph
	TokenText
	TokenStar
	TokenStarDouble
	TokenUnderscore
	TokenIndent

	TokenSpeaker
	TokenDialogue
	TokenParenthetical

	TokenCommentOpen
	TokenCommentClose
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
		lex.Backup()
		return lexBody
	}

	return lexDataKey
}

func lexBody(lex *lexer.Lexer) lexer.StateFn {
	if lex.Peek() == ' ' {
		return lexIndent
	}

	if lex.Peek() == '\n' {
		return lexParagraph
	}

	for {
		r := lex.NextRune()
		if r == lexer.Eof {
			break
		}
		if strings.IndexRune("abcdefghijklmnopqrstuvwxyz*_()[]", r) >= 0 {
			lex.Backup()
			return lexText
		}
		if r == '\n' {
			lex.Backup()
			return lexSpeaker
		}
	}
	return nil
}

func lexParagraph(lex *lexer.Lexer) lexer.StateFn {
	lex.AcceptRun("\n")
	lex.Emit(TokenParagraph)
	return lexBody
}

func lexIndent(lex *lexer.Lexer) lexer.StateFn {
	lex.AcceptRun(" ")
	lex.Emit(TokenIndent)
	return lexBody
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
		lex.Backup()
		return lexParagraph
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

		if r == '[' {
			lex.Backup()
			lex.Emit(TokenDialogue)
			lex.Accept("[")
			lex.Accept("[")
			lex.Emit(TokenCommentOpen)
			return lexDialogueText
		}

		if r == ']' {
			lex.Backup()
			lex.Emit(TokenDialogue)
			lex.Accept("]")
			lex.Accept("]")
			lex.Emit(TokenCommentClose)
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
			return lexParagraph
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

		if r == '[' {
			lex.Backup()
			lex.Emit(TokenText)
			lex.Accept("[")
			lex.Accept("[")
			lex.Emit(TokenCommentOpen)
			return lexText
		}

		if r == ']' {
			lex.Backup()
			lex.Emit(TokenText)
			lex.Accept("]")
			lex.Accept("]")
			lex.Emit(TokenCommentClose)
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
