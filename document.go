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

func (p *Paragraph) IsDialogue() bool {
	return p.Type == "dialogue"
}

func (p *Paragraph) Speaker() string {
	if p.IsDialogue() {
		return p.Lines[0].Chunks[0].Content
	}
	return ""
}

func (p *Paragraph) Dialogue() string {
	if p.IsDialogue() {
		dialogue := ""
		for _, line := range p.Lines {
			if line.Type == "dialogue" {
				for _, chunk := range line.Chunks {
					if !(len(chunk.Styles) == 1 && chunk.Styles[0] == "comment") {
						dialogue += " " + chunk.Content
					}
				}
			}
		}
		return dialogue
	}
	return ""
}

type Document struct {
	Title, Credit, Author, DraftDate string
	Data map[string]string
	Body []Paragraph
}
