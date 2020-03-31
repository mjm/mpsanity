package block

import (
	"fmt"

	"github.com/russross/blackfriday/v2"
)

func FromMarkdown(s string) ([]Block, error) {
	root := blackfriday.New().Parse([]byte(s))

	var b Builder

	root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		fmt.Println(node.Type.String())
		switch node.Type {
		case blackfriday.Document:
			break
		case blackfriday.Paragraph:
			if entering {
				b.StartBlock("normal")
			} else {
				b.EndBlock()
			}
		case blackfriday.Heading:
			if entering {
				style := fmt.Sprintf("h%d", node.Level)
				b.StartBlock(style)
			} else {
				b.EndBlock()
			}
		case blackfriday.Text:
			b.AppendText(string(node.Literal))
		case blackfriday.Emph:
			if entering {
				b.StartSpan("emphasis")
			} else {
				b.EndSpan()
			}
		case blackfriday.Strong:
			if entering {
				b.StartSpan("strong")
			} else {
				b.EndSpan()
			}
		case blackfriday.Code:
			b.StartSpan("code")
			b.AppendText(string(node.Literal))
			b.EndSpan()
		}

		return blackfriday.GoToNext
	})

	return b.Blocks(), nil
}
