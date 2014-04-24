package fountain

import (
	"testing"

	"github.com/exupero/state-lexer"
)

func TestData(t *testing.T) {
	script := `Title: The One Day
Credit: Written By`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenDataKey, "Credit"}
		s <- lexer.Token{TokenDataValue, "Written By"}
	})
}

func TestText(t *testing.T) {
	script := `Title: The One Day

A CROWD gathers.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A CROWD gathers."}
	})
}

func TestBold(t *testing.T) {
	script := `Title: The One Day

A **CROWD** gathers.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A "}
		s <- lexer.Token{TokenTextBold, "CROWD"}
		s <- lexer.Token{TokenText, " gathers."}
	})
}

func TestItalic(t *testing.T) {
	script := `Title: The One Day

A *CROWD* gathers.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A "}
		s <- lexer.Token{TokenTextItalic, "CROWD"}
		s <- lexer.Token{TokenText, " gathers."}
	})
}

func TestUnderline(t *testing.T) {
	script := `Title: The One Day

A _CROWD_ gathers.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "A "}
		s <- lexer.Token{TokenTextUnderline, "CROWD"}
		s <- lexer.Token{TokenText, " gathers."}
	})
}

func TestMultipleTextVariants(t *testing.T) {
	script := `Title: The One Day

**A** *CROWD* _gathers_.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
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
	script := `Title: The One Day

BOY
This is a sunny day!
(beat)
But I think it will rain...`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenSpeaker, "BOY"}
		s <- lexer.Token{TokenDialogue, "This is a sunny day!"}
		s <- lexer.Token{TokenParenthetical, "beat"}
		s <- lexer.Token{TokenDialogue, "But I think it will rain..."}
	})
}
