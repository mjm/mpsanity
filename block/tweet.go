package block

import (
	"regexp"

	"github.com/russross/blackfriday/v2"
)

const TypeTweet = "tweet"

type TweetContent struct {
	URL string `json:"url"`
}

var tweetLinkRegex = regexp.MustCompile("^https://twitter.com/.*/status/\\d+")

func TweetMarkdownRule(b *Builder, node *blackfriday.Node, entering bool) (blackfriday.WalkStatus, bool) {
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

	if !tweetLinkRegex.Match(child.Destination) {
		return blackfriday.GoToNext, false
	}

	b.AddCustomBlock(TypeTweet, &TweetContent{
		URL: string(child.Destination),
	})

	return blackfriday.SkipChildren, true
}
