package block

import (
	"fmt"
	"strings"

	"github.com/russross/blackfriday/v2"
)

type MarkdownConverter struct {
	rules []MarkdownRuleFunc
}

func NewMarkdownConverter(opts ...MarkdownOption) *MarkdownConverter {
	mc := &MarkdownConverter{}
	for _, o := range opts {
		o.Apply(mc)
	}
	return mc
}

func (mc *MarkdownConverter) newMarkdown() *blackfriday.Markdown {
	return blackfriday.New(blackfriday.WithExtensions(blackfriday.CommonExtensions))
}

type MarkdownRuleFunc func(b *Builder, node *blackfriday.Node, entering bool) (blackfriday.WalkStatus, bool)

type MarkdownOption interface {
	Apply(mc *MarkdownConverter)
}

type markdownOptionFn func(mc *MarkdownConverter)

func (f markdownOptionFn) Apply(mc *MarkdownConverter) {
	f(mc)
}

func WithMarkdownRules(rules ...MarkdownRuleFunc) MarkdownOption {
	return markdownOptionFn(func(mc *MarkdownConverter) {
		mc.rules = append(mc.rules, rules...)
	})
}

func (mc *MarkdownConverter) ToBlocks(s string) ([]Block, error) {
	root := mc.newMarkdown().Parse([]byte(s))

	var b Builder

	var lastLinkKey string
	root.Walk(func(node *blackfriday.Node, entering bool) blackfriday.WalkStatus {
		// if entering {
		// 	fmt.Printf("entering %s\n", node.Type.String())
		// } else {
		// 	fmt.Printf("exiting %s\n", node.Type.String())
		// }

		for _, rule := range mc.rules {
			if result, ok := rule(&b, node, entering); ok {
				return result
			}
		}

		switch node.Type {
		case blackfriday.Document:
			break
		case blackfriday.Paragraph:
			if node.Parent != nil &&
				(node.Parent.Type == blackfriday.Item || node.Parent.Type == blackfriday.BlockQuote) {
				if entering && node.Prev != nil {
					b.AppendText("\n\n")
				}
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
			text := strings.ReplaceAll(string(node.Literal), "\n", " ")
			b.AppendText(text)
		case blackfriday.Hardbreak:
			b.AppendText("\n")
		case blackfriday.Emph:
			if entering {
				b.StartMark("em")
			} else {
				b.EndMark("em")
			}
		case blackfriday.Strong:
			if entering {
				b.StartMark("strong")
			} else {
				b.EndMark("strong")
			}
		case blackfriday.Code:
			b.StartMark("code")
			b.AppendText(string(node.Literal))
			b.EndMark("code")
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
		case blackfriday.BlockQuote:
			if entering {
				b.StartBlock("blockquote")
			} else {
				b.EndBlock()
			}
		case blackfriday.CodeBlock:
			b.AddCustomBlock("code", &CodeContent{
				Language: string(node.Info),
				Code:     strings.TrimSuffix(string(node.Literal), "\n"),
			})
		case blackfriday.Link:
			if entering {
				lastLinkKey = b.AddMarkDef("link", &LinkData{
					Href: string(node.Destination),
				})
				b.StartMark(lastLinkKey)
			} else {
				b.EndMark(lastLinkKey)
			}
		}

		return blackfriday.GoToNext
	})

	return b.Blocks(), nil
}
