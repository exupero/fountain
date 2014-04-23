package fountain

import (
	"testing"

	"github.com/exupero/state-lexer"
)

func TestData(t *testing.T) {
	lexer.AssertStream(t, "Title: The One Day\nCredit: Written By", Tokenize, func(s chan lexer.Token) {
		s <- lexer.Token{TokenDataKey, "Title"}
		s <- lexer.Token{TokenDataValue, "The One Day"}
		s <- lexer.Token{TokenDataKey, "Credit"}
		s <- lexer.Token{TokenDataValue, "Written By"}
	})
}
