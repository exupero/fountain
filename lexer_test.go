package fountain

import (
	"testing"

	"github.com/exupero/state-lexer"
)

func TestData(t *testing.T) {
	lexer.AssertStream(t, "Title: The One Day", Tokenize, func(s chan lexer.Token) {
		s <- lexer.Token{TokenText, "Title"}
		s <- lexer.Token{TokenColon, ":"}
		s <- lexer.Token{TokenText, "The One Day"}
	})
}
