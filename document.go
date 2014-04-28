package fountain

type Chunk struct {
	Content string
	Styles []string
}

type Line struct {
	Chunks []Chunk
	Type string
}

type Paragraph struct {
	Lines []Line
	Type string
}

type Document struct {
	Title, Credit, Author, DraftDate string
	Data map[string]string
	Body []Paragraph
}

