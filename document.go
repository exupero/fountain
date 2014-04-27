package fountain

import (
	"fmt"
)

type Text struct {
	content string
	styles []string
}

func (t Text) String() string {
	return fmt.Sprintf("{\"%s\" %s}", t.content, t.styles)
}

type Line []Text

type Document struct {
	Title, Credit, Author, DraftDate string
	Data map[string]string
	Body []Line
}

