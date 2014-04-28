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
			Lines: []Line {
				Line{
					Chunks: []Chunk{
						Chunk{Content: "The MEN ran down the "},
						Chunk{Content: "street", Styles: []string{"italic"}},
						Chunk{Content: ". "},
						Chunk{Content: "They ", Styles: []string{"italic"}},
						Chunk{Content: "jumped", Styles: []string{"bold", "italic"}},
						Chunk{Content: " into the ", Styles: []string{"italic"}},
						Chunk{Content: "ditch", Styles: []string{"italic", "underline"}},
						Chunk{Content: ".", Styles: []string{"italic"}},
						Chunk{Content: ""},
					},
					Type: "action",
				},
			},
			Type: "action",
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
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "The BOYS cheered."},
					},
					Type: "action",
				},
			},
			Type: "action",
		},
		Paragraph{
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "The WOMEN sang."},
					},
					Type: "action",
				},
			},
			Type: "action",
		},
		Paragraph{
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "The MEN stood."},
					},
					Type: "action",
				},
			},
			Type: "action",
		},
	})
}

func TestDocDialogue(t *testing.T) {
	script := `Title: The One Day

BOY
(triumphantly)
*I think _that's_ a **great** idea*!

GIRL
Sure it is...`

	assertBody(t, script, []Paragraph{
		Paragraph{
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "BOY"},
					},
					Type: "speaker",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: "triumphantly"},
					},
					Type: "parenthetical",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: ""},
						Chunk{Content: "I think ", Styles: []string{"italic"}},
						Chunk{Content: "that's", Styles: []string{"italic", "underline"}},
						Chunk{Content: " a ", Styles: []string{"italic"}},
						Chunk{Content: "great", Styles: []string{"bold", "italic"}},
						Chunk{Content: " idea", Styles: []string{"italic"}},
						Chunk{Content: "!"},
					},
					Type: "dialogue",
				},
			},
			Type: "dialogue",
		},
		Paragraph{
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "GIRL"},
					},
					Type: "speaker",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: "Sure it is..."},
					},
					Type: "dialogue",
				},
			},
		},
	})
}

func TestMultilineAction(t *testing.T) {
	script := `Title: The One Day

This action
takes place over
multiple lines.`
	assertBody(t, script, []Paragraph{
		Paragraph{
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "This action"},
					},
					Type: "action",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: "takes place over"},
					},
					Type: "action",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: "multiple lines."},
					},
					Type: "action",
				},
			},
			Type: "action",
		},
	})
}

func TestParentheticalAfterDialogue(t *testing.T) {
	script := `Title: The One Day

BOY
This is a great day!
(pause)
Or is it...`
	assertBody(t, script, []Paragraph{
		Paragraph{
			Lines: []Line{
				Line{
					Chunks: []Chunk{
						Chunk{Content: "BOY"},
					},
					Type: "speaker",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: "This is a great day!"},
					},
					Type: "dialogue",
				},
				Line{
					Chunks: []Chunk{
						Chunk{Content: "pause"},
					},
					Type: "parenthetical",
				},
				Line {
					Chunks: []Chunk{
						Chunk{Content: "Or is it..."},
					},
					Type: "dialogue",
				},
			},
			Type: "dialogue",
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

		if len(actualParagraph.Lines) != len(expectedParagraph.Lines) {
			mismatch()
			return
		}

		for j, expectedLine := range expectedParagraph.Lines {
			actualLine := doc.Body[i].Lines[j]

			if actualLine.Type != expectedLine.Type {
				mismatch()
				return
			}

			if len(actualLine.Chunks) != len(expectedLine.Chunks) {
				mismatch()
				return
			}

			for k, expectedChunk := range expectedLine.Chunks {
				actualChunk := actualLine.Chunks[k]

				if actualChunk.Content != expectedChunk.Content {
					mismatch()
					return
				}

				if len(actualChunk.Styles) != len(expectedChunk.Styles) {
					mismatch()
					return
				}

				for m, expectedStyle := range expectedChunk.Styles {
					actualStyle := actualChunk.Styles[m]

					if actualStyle != expectedStyle {
						mismatch()
						return
					}
				}
			}
		}
	}
}
