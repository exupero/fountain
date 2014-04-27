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

	assertBody(t, script, []Line{
		Line{
			Text{content: "The MEN ran down the ", styles: []string{}},
			Text{content: "street", styles: []string{"italic"}},
			Text{content: ". ", styles: []string{}},
			Text{content: "They ", styles: []string{"italic"}},
			Text{content: "jumped", styles: []string{"bold", "italic"}},
			Text{content: " into the ", styles: []string{"italic"}},
			Text{content: "ditch", styles: []string{"italic", "underline"}},
			Text{content: ".", styles: []string{"italic"}},
			Text{content: "", styles: []string{}},
		},
	})
}

func TestDocTextParagraphs(t *testing.T) {
	script := `Title: The One Day

The BOYS cheered.

The WOMEN sang.

The MEN stood.`

	assertBody(t, script, []Line{
		Line{Text{content: "The BOYS cheered.", styles: []string{}}},
		Line{Text{content: "The WOMEN sang.", styles: []string{}}},
		Line{Text{content: "The MEN stood.", styles: []string{}}},
	})
}

func assertBody(t *testing.T, script string, expectedLines []Line) {
	doc := Parse(script)

	mismatch := func() {
		t.Errorf(`Body wrong.
Expected: %s
Actual:   %s
Source:
-------
%s`, expectedLines, doc.Body, script)
	}

	if len(doc.Body) != len(expectedLines) {
		mismatch()
		return
	}

	for i, expectedLine := range expectedLines {
		if len(doc.Body[i]) != len(expectedLine) {
			mismatch()
			return
		}

		for j, expectedChunk := range expectedLine {
			actual := doc.Body[i][j]
			if actual.content != expectedChunk.content {
				mismatch()
				return
			}

			if len(actual.styles) != len(expectedChunk.styles) {
				mismatch()
				return
			}

			for k, expectedStyle := range expectedChunk.styles {
				if actual.styles[k] != expectedStyle {
					mismatch()
					return
				}
			}
		}
	}
}
