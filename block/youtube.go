package block

import (
	"regexp"

	"github.com/russross/blackfriday/v2"
)

const TypeYouTube = "youtube"

type YouTubeContent struct {
	URL string `json:"url"`
}

var youTubeLinkRegex = regexp.MustCompile("^https://www.youtube.com/watch\\?v=\\w+")

func YouTubeMarkdownRule(b *Builder, node *blackfriday.Node, entering bool) (blackfriday.WalkStatus, bool) {
	if !entering || node.Type != blackfriday.Paragraph {
		return blackfriday.GoToNext, false
	}

	child := node.FirstChild
	if child != nil && child.Type == blackfriday.Text && len(child.Literal) == 0 {
		child = child.Next
	}
	if child == nil || child.Type != blackfriday.Link || child.Next != nil {
		return blackfriday.GoToNext, false
	}

	if !youTubeLinkRegex.Match(child.Destination) {
		return blackfriday.GoToNext, false
	}

	b.AddCustomBlock(TypeYouTube, &YouTubeContent{
		URL: string(child.Destination),
	})

	return blackfriday.SkipChildren, true
}
