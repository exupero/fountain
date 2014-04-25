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
		s <- lexer.Token{TokenStarDouble, "**"}
		s <- lexer.Token{TokenText, "CROWD"}
		s <- lexer.Token{TokenStarDouble, "**"}
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
		s <- lexer.Token{TokenStar, "*"}
		s <- lexer.Token{TokenText, "CROWD"}
		s <- lexer.Token{TokenStar, "*"}
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
		s <- lexer.Token{TokenUnderscore, "_"}
		s <- lexer.Token{TokenText, "CROWD"}
		s <- lexer.Token{TokenUnderscore, "_"}
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
		s <- lexer.Token{TokenStarDouble, "**"}
		s <- lexer.Token{TokenText, "A"}
		s <- lexer.Token{TokenStarDouble, "**"}
		s <- lexer.Token{TokenText, " "}
		s <- lexer.Token{TokenStar, "*"}
		s <- lexer.Token{TokenText, "CROWD"}
		s <- lexer.Token{TokenStar, "*"}
		s <- lexer.Token{TokenText, " "}
		s <- lexer.Token{TokenUnderscore, "_"}
		s <- lexer.Token{TokenText, "gathers"}
		s <- lexer.Token{TokenUnderscore, "_"}
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

func TestDialogueTextVariants(t *testing.T) {
	script := `Title: The One Day

BOY
**This** *is* a _sunny_ day!`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenSpeaker, "BOY"}
		s <- lexer.Token{TokenDialogue, ""}
		s <- lexer.Token{TokenStarDouble, "**"}
		s <- lexer.Token{TokenDialogue, "This"}
		s <- lexer.Token{TokenStarDouble, "**"}
		s <- lexer.Token{TokenDialogue, " "}
		s <- lexer.Token{TokenStar, "*"}
		s <- lexer.Token{TokenDialogue, "is"}
		s <- lexer.Token{TokenStar, "*"}
		s <- lexer.Token{TokenDialogue, " a "}
		s <- lexer.Token{TokenUnderscore, "_"}
		s <- lexer.Token{TokenDialogue, "sunny"}
		s <- lexer.Token{TokenUnderscore, "_"}
		s <- lexer.Token{TokenDialogue, " day!"}
	})
}

func TestDialogueAndText(t *testing.T) {
	script := `Title: The One Day

The sun shines...

BOY
This is a sunny day!
(beat)
But I think it will rain...

The rain starts...

...and then it pours.

BOY
I knew it.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, "The sun shines..."}
		s <- lexer.Token{TokenSpeaker, "BOY"}
		s <- lexer.Token{TokenDialogue, "This is a sunny day!"}
		s <- lexer.Token{TokenParenthetical, "beat"}
		s <- lexer.Token{TokenDialogue, "But I think it will rain..."}
		s <- lexer.Token{TokenText, "The rain starts..."}
		s <- lexer.Token{TokenText, "...and then it pours."}
		s <- lexer.Token{TokenSpeaker, "BOY"}
		s <- lexer.Token{TokenDialogue, "I knew it."}
	})
}

func TestIndentation(t *testing.T) {
	script := `Title: The One Day

    This is indented text.`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenIndent, "    "}
		s <- lexer.Token{TokenText, "This is indented text."}
	})
}

func TestComment(t *testing.T) {
	script := `Title: The One Day

[[Unfinished]]

The End`
	lexer.AssertStream(t, Tokenize, script, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenText, ""}
		s <- lexer.Token{TokenCommentOpen, "[["}
		s <- lexer.Token{TokenText, "Unfinished"}
		s <- lexer.Token{TokenCommentClose, "]]"}
		s <- lexer.Token{TokenText, ""}
		s <- lexer.Token{TokenText, "The End"}
	})
}
