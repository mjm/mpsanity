package block

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarkdownToBlocks(t *testing.T) {
	cases := []struct {
		name   string
		input  string
		output []Block
	}{
		{
			name:  "simple paragraph",
			input: `This is some text.`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "This is some text.",
								},
							},
						},
					},
				},
			},
		},
		{
			name:  "paragraph with formatting",
			input: "This is _some text_ with **formatting**. Including `code.`",
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "This is ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "some text",
									Marks: []string{"emphasis"},
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text: " with ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "formatting",
									Marks: []string{"strong"},
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text: ". Including ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "code.",
									Marks: []string{"code"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "multiple paragraphs",
			input: `This is some text.

And this is some more text.

And how about _some more?_`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "This is some text.",
								},
							},
						},
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "And this is some more text.",
								},
							},
						},
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "And how about ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "some more?",
									Marks: []string{"emphasis"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "headings",
			input: `# This is a heading

And a paragraph underneath

## And another heading`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "h1",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "This is a heading",
								},
							},
						},
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "And a paragraph underneath",
								},
							},
						},
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "h2",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "And another heading",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "unordered lists",
			input: `* Unordered list
* Second item
    * Second level item
    * Another one
* Back to the first level
    * End on a second level item

Now there's a paragraph at the end.
`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Unordered list",
								},
							},
						},
						ListItem: "bullet",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Second item",
								},
							},
						},
						ListItem: "bullet",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Second level item",
								},
							},
						},
						ListItem: "bullet",
						Level:    2,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Another one",
								},
							},
						},
						ListItem: "bullet",
						Level:    2,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Back to the first level",
								},
							},
						},
						ListItem: "bullet",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "End on a second level item",
								},
							},
						},
						ListItem: "bullet",
						Level:    2,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Now there's a paragraph at the end.",
								},
							},
						},
					},
				},
			},
		},
		{
			name: "ordered lists",
			input: `1. First item
2. Second item
3. Third item`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "First item",
								},
							},
						},
						ListItem: "number",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Second item",
								},
							},
						},
						ListItem: "number",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Third item",
								},
							},
						},
						ListItem: "number",
						Level:    1,
					},
				},
			},
		},
		{
			name: "lists with paragraphs",
			input: `* This is a first paragraph

    This should be part of the same list items

* Another list item
* A final list item`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "This is a first paragraph\n\nThis should be part of the same list items",
								},
							},
						},
						ListItem: "bullet",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Another list item",
								},
							},
						},
						ListItem: "bullet",
						Level:    1,
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "A final list item",
								},
							},
						},
						ListItem: "bullet",
						Level:    1,
					},
				},
			},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := FromMarkdown(c.input)
			assert.NoError(t, err)
			assert.Equal(t, c.output, out)
		})
	}
}
