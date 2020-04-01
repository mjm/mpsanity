package block

import (
	"fmt"

	"github.com/russross/blackfriday/v2"
)

func FromMarkdown(s string) ([]Block, error) {
	root := blackfriday.New().Parse([]byte(s))

	var b Builder

	root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		if entering {
			fmt.Printf("entering %s\n", node.Type.String())
		} else {
			fmt.Printf("exiting %s\n", node.Type.String())
		}
		switch node.Type {
		case blackfriday.Document:
			break
		case blackfriday.Paragraph:
			if node.Parent != nil && node.Parent.Type == blackfriday.Item {
				break
			}

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
		case blackfriday.List:
			if entering {
				if node.ListFlags&blackfriday.ListTypeOrdered != 0 {
					b.StartList("number")
				} else {
					b.StartList("bullet")
				}
			} else {
				b.EndList()
			}
		case blackfriday.Item:
			if entering {
				b.StartListItem()
			} else {
				b.EndListItem()
			}
		}

		return blackfriday.GoToNext
	})

	return b.Blocks(), nil
}
