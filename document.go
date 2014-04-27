package fountain

type Chunk struct {
	content string
	styles []string
}

type Line struct {
	chunks []Chunk
	typ string
}

type Paragraph struct {
	lines []Line
	typ string
}

type Document struct {
	Title, Credit, Author, DraftDate string
	Data map[string]string
	Body []Paragraph
}

