package fountain

import (
	"testing"

	"github.com/exupero/state-lexer"
)

func TestData(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\nCredit: Written By", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenDataKey, "Credit"}
		s <- lexer.Token{TokenDataValue, "Written By"}
	})
}

func TestText(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\n\nA CROWD gathers.", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A CROWD gathers."}
	})
}

func TestBold(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\n\nA **CROWD** gathers.", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A "}
		s <- lexer.Token{TokenTextBold, "CROWD"}
		s <- lexer.Token{TokenText, " gathers."}
	})
}

func TestItalic(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\n\nA *CROWD* gathers.", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A "}
		s <- lexer.Token{TokenTextItalic, "CROWD"}
		s <- lexer.Token{TokenText, " gathers."}
	})
}

func TestUnderline(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\n\nA _CROWD_ gathers.", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A "}
		s <- lexer.Token{TokenTextUnderline, "CROWD"}
		s <- lexer.Token{TokenText, " gathers."}
	})
}

func TestMultipleTextVariants(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\n\n**A** *CROWD* _gathers_.", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, ""}
		s <- lexer.Token{TokenTextBold, "A"}
		s <- lexer.Token{TokenText, " "}
		s <- lexer.Token{TokenTextItalic, "CROWD"}
		s <- lexer.Token{TokenText, " "}
		s <- lexer.Token{TokenTextUnderline, "gathers"}
		s <- lexer.Token{TokenText, "."}
	})
}

func TestDialogue(t *testing.T) {
	lexer.AssertStream(t, Tokenize, "Title: The One Day\n\nBOY\nThis is a sunny day!\n(beat)\nBut I think it will rain...", func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenSpeaker, "BOY"}
		s <- lexer.Token{TokenDialogue, "This is a sunny day!"}
		s <- lexer.Token{TokenParenthetical, "beat"}
		s <- lexer.Token{TokenDialogue, "But I think it will rain..."}
	})
}
