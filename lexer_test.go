package fountain

import (
	"testing"

	"github.com/exupero/state-lexer"
)

func TestData(t *testing.T) {
	script := `Title: The One Day
Credit: Written By`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenDataKey, "Credit"},
		lexer.Token{TokenDataValue, "Written By"},
	})
}

func TestText(t *testing.T) {
	script := `Title: The One Day

A CROWD gathers.`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, "A CROWD gathers."},
	})
}

func TestBold(t *testing.T) {
	script := `Title: The One Day

A **CROWD** gathers.`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, "A "},
		lexer.Token{TokenStarDouble, "**"},
		lexer.Token{TokenText, "CROWD"},
		lexer.Token{TokenStarDouble, "**"},
		lexer.Token{TokenText, " gathers."},
	})
}

func TestItalic(t *testing.T) {
	script := `Title: The One Day

A *CROWD* gathers.`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, "A "},
		lexer.Token{TokenStar, "*"},
		lexer.Token{TokenText, "CROWD"},
		lexer.Token{TokenStar, "*"},
		lexer.Token{TokenText, " gathers."},
	})
}

func TestUnderline(t *testing.T) {
	script := `Title: The One Day

A _CROWD_ gathers.`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, "A "},
		lexer.Token{TokenUnderscore, "_"},
		lexer.Token{TokenText, "CROWD"},
		lexer.Token{TokenUnderscore, "_"},
		lexer.Token{TokenText, " gathers."},
	})
}

func TestMultipleTextVariants(t *testing.T) {
	script := `Title: The One Day

**A** *CROWD* _gathers_.`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, ""},
		lexer.Token{TokenStarDouble, "**"},
		lexer.Token{TokenText, "A"},
		lexer.Token{TokenStarDouble, "**"},
		lexer.Token{TokenText, " "},
		lexer.Token{TokenStar, "*"},
		lexer.Token{TokenText, "CROWD"},
		lexer.Token{TokenStar, "*"},
		lexer.Token{TokenText, " "},
		lexer.Token{TokenUnderscore, "_"},
		lexer.Token{TokenText, "gathers"},
		lexer.Token{TokenUnderscore, "_"},
		lexer.Token{TokenText, "."},
	})
}

func TestDialogue(t *testing.T) {
	script := `Title: The One Day

BOY
This is a sunny day!
(beat)
But I think it will rain...`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenSpeaker, "BOY"},
		lexer.Token{TokenDialogue, "This is a sunny day!"},
		lexer.Token{TokenParenthetical, "beat"},
		lexer.Token{TokenDialogue, "But I think it will rain..."},
	})
}

func TestDialogueTextVariants(t *testing.T) {
	script := `Title: The One Day

BOY
**This** *is* a _sunny_ day!`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenSpeaker, "BOY"},
		lexer.Token{TokenDialogue, ""},
		lexer.Token{TokenStarDouble, "**"},
		lexer.Token{TokenDialogue, "This"},
		lexer.Token{TokenStarDouble, "**"},
		lexer.Token{TokenDialogue, " "},
		lexer.Token{TokenStar, "*"},
		lexer.Token{TokenDialogue, "is"},
		lexer.Token{TokenStar, "*"},
		lexer.Token{TokenDialogue, " a "},
		lexer.Token{TokenUnderscore, "_"},
		lexer.Token{TokenDialogue, "sunny"},
		lexer.Token{TokenUnderscore, "_"},
		lexer.Token{TokenDialogue, " day!"},
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
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, "The sun shines..."},
		lexer.Token{TokenParagraph, "\n\n"},
		lexer.Token{TokenSpeaker, "BOY"},
		lexer.Token{TokenDialogue, "This is a sunny day!"},
		lexer.Token{TokenParenthetical, "beat"},
		lexer.Token{TokenDialogue, "But I think it will rain..."},
		lexer.Token{TokenParagraph, "\n\n"},
		lexer.Token{TokenText, "The rain starts..."},
		lexer.Token{TokenParagraph, "\n\n"},
		lexer.Token{TokenText, "...and then it pours."},
		lexer.Token{TokenParagraph, "\n\n"},
		lexer.Token{TokenSpeaker, "BOY"},
		lexer.Token{TokenDialogue, "I knew it."},
	})
}

func TestIndentation(t *testing.T) {
	script := `Title: The One Day

    This is indented text.`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenIndent, "    "},
		lexer.Token{TokenText, "This is indented text."},
	})
}

func TestComment(t *testing.T) {
	script := `Title: The One Day

[[Unfinished]]

The End`
	lexer.AssertStream(t, Tokenize, script, []lexer.Token{
		lexer.Token{TokenDataKey, "Title"},
		lexer.Token{TokenDataValue, "The One Day"},
		lexer.Token{TokenParagraph, "\n"},
		lexer.Token{TokenText, ""},
		lexer.Token{TokenCommentOpen, "[["},
		lexer.Token{TokenText, "Unfinished"},
		lexer.Token{TokenCommentClose, "]]"},
		lexer.Token{TokenText, ""},
		lexer.Token{TokenParagraph, "\n\n"},
		lexer.Token{TokenText, "The End"},
	})
}
