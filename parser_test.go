package fountain

import "testing"

func TestDocData(t *testing.T) {
	script := `Title: The One Day
Credit: Written By
Author: Some Body
Draft Date: 02/14/14
Quality: Pretty Good`
	doc := Parse(script)

	if doc.Title != "The One Day" {
		t.Errorf("Title is not '%s', but is '%s'", "The One Day", doc.Title)
	}

	if doc.Credit != "Written By" {
		t.Errorf("Credit is not '%s', but is '%s'", "Written By", doc.Credit)
	}

	if doc.Author != "Some Body" {
		t.Errorf("Author is not '%s', but is '%s'", "Some Body", doc.Author)
	}

	if doc.DraftDate != "02/14/14" {
		t.Errorf("Draft Date is not '%s', but is '%s'", "02/14/14", doc.DraftDate)
	}

	if doc.Data["Quality"] != "Pretty Good" {
		t.Errorf("Data[Quality] is not '%s', but is '%s'", "Pretty Good", doc.Data["Quality"])
	}
}

func TestDocTextVariants(t *testing.T) {
	script := `Title: The One Day

The MEN ran down the *street*. *They **jumped** into the _ditch_.*`

	assertBody(t, script, []Paragraph{
		Paragraph{
			lines: []Line {
				Line{
					chunks: []Chunk{
						Chunk{content: "The MEN ran down the ", styles: []string{}},
						Chunk{content: "street", styles: []string{"italic"}},
						Chunk{content: ". ", styles: []string{}},
						Chunk{content: "They ", styles: []string{"italic"}},
						Chunk{content: "jumped", styles: []string{"bold", "italic"}},
						Chunk{content: " into the ", styles: []string{"italic"}},
						Chunk{content: "ditch", styles: []string{"italic", "underline"}},
						Chunk{content: ".", styles: []string{"italic"}},
						Chunk{content: "", styles: []string{}},
					},
					typ: "action",
				},
			},
			typ: "action",
		},
	})
}

func TestDocTextParagraphs(t *testing.T) {
	script := `Title: The One Day

The BOYS cheered.

The WOMEN sang.

The MEN stood.`

	assertBody(t, script, []Paragraph{
		Paragraph{
			lines: []Line{
				Line{
					chunks: []Chunk{
						Chunk{content: "The BOYS cheered.", styles: []string{}},
					},
					typ: "action",
				},
			},
			typ: "action",
		},
		Paragraph{
			lines: []Line{
				Line{
					chunks: []Chunk{
						Chunk{content: "The WOMEN sang.", styles: []string{}},
					},
					typ: "action",
				},
			},
			typ: "action",
		},
		Paragraph{
			lines: []Line{
				Line{
					chunks: []Chunk{
						Chunk{content: "The MEN stood.", styles: []string{}},
					},
					typ: "action",
				},
			},
			typ: "action",
		},
	})
}

func TestDocDialogue(t *testing.T) {
	script := `Title: The One Day

BOY
(triumphantly)
I think that's a great idea!`

	assertBody(t, script, []Paragraph{
		Paragraph{
			lines: []Line{
				Line{
					chunks: []Chunk{
						Chunk{content: "BOY", styles: []string{}},
					},
					typ: "speaker",
				},
				Line{
					chunks: []Chunk{
						Chunk{content: "triumphantly", styles: []string{}},
					},
					typ: "parenthetical",
				},
				Line{
					chunks: []Chunk{
						Chunk{content: "I think that's a great idea!", styles: []string{}},
					},
					typ: "dialogue",
				},
			},
			typ: "dialogue",
		},
	})
}

func assertBody(t *testing.T, script string, expectedParagraphs []Paragraph) {
	doc := Parse(script)

	mismatch := func() {
		t.Errorf(`Body wrong.
Expected: %s
Actual:   %s
Source:
-------
%s`, expectedParagraphs, doc.Body, script)
	}

	if len(doc.Body) != len(expectedParagraphs) {
		mismatch()
		return
	}

	for i, expectedParagraph := range expectedParagraphs {
		actualParagraph := doc.Body[i]

		if len(actualParagraph.lines) != len(expectedParagraph.lines) {
			mismatch()
			return
		}

		for j, expectedLine := range expectedParagraph.lines {
			actualLine := doc.Body[i].lines[j]

			if actualLine.typ != expectedLine.typ {
				mismatch()
				return
			}

			if len(actualLine.chunks) != len(expectedLine.chunks) {
				mismatch()
				return
			}

			for k, expectedChunk := range expectedLine.chunks {
				actualChunk := actualLine.chunks[k]

				if actualChunk.content != expectedChunk.content {
					mismatch()
					return
				}

				if len(actualChunk.styles) != len(expectedChunk.styles) {
					mismatch()
					return
				}

				for m, expectedStyle := range expectedChunk.styles {
					actualStyle := actualChunk.styles[m]

					if actualStyle != expectedStyle {
						mismatch()
						return
					}
				}
			}
		}
	}
}
