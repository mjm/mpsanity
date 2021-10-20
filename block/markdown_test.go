package block

import (
	"encoding/json"
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
						MarkDefs: []MarkDef{},
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
									Marks: []string{"em"},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
									Marks: []string{"em"},
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
			},
		},
		{
			name: "line breaks",
			input: `Here is some text  
with some forced
line breaks.  
Let's see how they go.`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Here is some text\nwith some forced line breaks.\nLet's see how they go.",
								},
							},
						},
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
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
						MarkDefs: []MarkDef{},
					},
				},
			},
		},
		{
			name: "blockquotes",
			input: `> A first line of blockquoted material

A paragraph in-between

> First quoted paragraph
>
> Second quoted paragraph`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "blockquote",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "A first line of blockquoted material",
								},
							},
						},
						MarkDefs: []MarkDef{},
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
									Text: "A paragraph in-between",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
				{
					Type: "block",
					Content: &BlockContent{
						Style: "blockquote",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "First quoted paragraph\n\nSecond quoted paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
			},
		},
		{
			name: "indented code blocks",
			input: `A normal paragraph

    const foo = 1
    const bar = 2

A following paragraph`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "A normal paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
				{
					Type: "code",
					Content: &CodeContent{
						Code: "const foo = 1\nconst bar = 2",
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
									Text: "A following paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
			},
		},
		{
			name:  "fenced code blocks",
			input: "A normal paragraph\n\n```\nconst foo = 1\nconst bar = 2\n```\n\nA following paragraph",
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "A normal paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
				{
					Type: "code",
					Content: &CodeContent{
						Code: "const foo = 1\nconst bar = 2",
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
									Text: "A following paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
			},
		},
		{
			name:  "fenced code blocks with language",
			input: "A normal paragraph\n\n```js\nconst foo = 1\nconst bar = 2\n```\n\nA following paragraph",
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "A normal paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
				{
					Type: "code",
					Content: &CodeContent{
						Language: "js",
						Code:     "const foo = 1\nconst bar = 2",
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
									Text: "A following paragraph",
								},
							},
						},
						MarkDefs: []MarkDef{},
					},
				},
			},
		},
		{
			name: "links",
			input: `Here are some [link examples](http://example.org/).

Can we [include _formatting_](/foo/bar) inside the links?

Links can be [done][] with footnotes too.

[done]: /foo/bar`,
			output: []Block{
				{
					Type: "block",
					Content: &BlockContent{
						Style: "normal",
						Children: []Block{
							{
								Type: "span",
								Content: &SpanContent{
									Text: "Here are some ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "link examples",
									Marks: []string{"mark1"},
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text: ".",
								},
							},
						},
						MarkDefs: []MarkDef{
							{
								Type: "link",
								Key:  "mark1",
								Data: &LinkData{
									Href: "http://example.org/",
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
									Text: "Can we ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "include ",
									Marks: []string{"mark1"},
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "formatting",
									Marks: []string{"em", "mark1"},
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text: " inside the links?",
								},
							},
						},
						MarkDefs: []MarkDef{
							{
								Type: "link",
								Key:  "mark1",
								Data: &LinkData{
									Href: "/foo/bar",
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
									Text: "Links can be ",
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text:  "done",
									Marks: []string{"mark1"},
								},
							},
							{
								Type: "span",
								Content: &SpanContent{
									Text: " with footnotes too.",
								},
							},
						},
						MarkDefs: []MarkDef{
							{
								Type: "link",
								Key:  "mark1",
								Data: &LinkData{
									Href: "/foo/bar",
								},
							},
						},
					},
				},
			},
		},
	}

	mc := NewMarkdownConverter()

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			out, err := mc.ToBlocks(c.input)
			assert.NoError(t, err)
			assert.Equal(t, c.output, out)
		})
	}
}

func TestMarkdownTweetRule(t *testing.T) {
	mc := NewMarkdownConverter(WithMarkdownRules(TweetMarkdownRule))

	out, err := mc.ToBlocks(`This is some content with an embedded tweet.

https://twitter.com/some_user/status/1234567890

And some more content afterwards.`)
	assert.NoError(t, err)

	assert.Equal(t, []Block{
		{
			Type: "block",
			Content: &BlockContent{
				Style: "normal",
				Children: []Block{
					{
						Type: "span",
						Content: &SpanContent{
							Text: "This is some content with an embedded tweet.",
						},
					},
				},
				MarkDefs: []MarkDef{},
			},
		},
		{
			Type: "tweet",
			Content: &TweetContent{
				URL: "https://twitter.com/some_user/status/1234567890",
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
							Text: "And some more content afterwards.",
						},
					},
				},
				MarkDefs: []MarkDef{},
			},
		},
	}, out)
}

func TestMarkdownYouTubeRule(t *testing.T) {
	mc := NewMarkdownConverter(WithMarkdownRules(YouTubeMarkdownRule))

	out, err := mc.ToBlocks(`This is some content with an embedded YouTube video.

https://www.youtube.com/watch?v=TamwFUUd9Yk

And some more content afterwards.`)
	assert.NoError(t, err)

	assert.Equal(t, []Block{
		{
			Type: "block",
			Content: &BlockContent{
				Style: "normal",
				Children: []Block{
					{
						Type: "span",
						Content: &SpanContent{
							Text: "This is some content with an embedded YouTube video.",
						},
					},
				},
				MarkDefs: []MarkDef{},
			},
		},
		{
			Type: "youtube",
			Content: &YouTubeContent{
				URL: "https://www.youtube.com/watch?v=TamwFUUd9Yk",
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
							Text: "And some more content afterwards.",
						},
					},
				},
				MarkDefs: []MarkDef{},
			},
		},
	}, out)
}

func TestBlockJSONRoundTrip(t *testing.T) {
	b := Block{
		Type: "block",
		Content: &BlockContent{
			Style: "normal",
			Children: []Block{
				{
					Type: "span",
					Content: &SpanContent{
						Text: "Some content",
					},
				},
			},
			MarkDefs: []MarkDef{
				{
					Type: "link",
					Key:  "asdf",
					Data: &LinkData{
						Href: "http://example.org",
					},
				},
			},
		},
	}

	data, err := json.Marshal(b)
	assert.NoError(t, err)

	var m map[string]interface{}
	assert.NoError(t, json.Unmarshal(data, &m))

	assert.Equal(t, map[string]interface{}{
		"_type": "block",
		"style": "normal",
		"children": []interface{}{
			map[string]interface{}{
				"_type": "span",
				"text":  "Some content",
			},
		},
		"markDefs": []interface{}{
			map[string]interface{}{
				"_type": "link",
				"_key":  "asdf",
				"href":  "http://example.org",
			},
		},
	}, m)
}
