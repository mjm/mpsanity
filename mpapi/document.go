package mpapi

import (
	"context"
	"time"

	"github.com/gosimple/slug"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/mpsanity"
	"github.com/mjm/mpsanity/block"
)

var ErrNotEntry = status.Error(codes.InvalidArgument, "post is not an entry")

type DocumentBuilder interface {
	BuildDocument(ctx context.Context, input *CreateInput) ([]interface{}, error)
}

type DefaultDocumentBuilder struct {
	MarkdownConverter *block.MarkdownConverter
}

func (d *DefaultDocumentBuilder) BuildDocument(_ context.Context, input *CreateInput) ([]interface{}, error) {
	if input.Type[0] != "entry" {
		return nil, ErrNotEntry
	}

	var doc defaultDocument
	results := []interface{}{&doc}

	if len(input.Props.Slug) > 0 {
		doc.Slug = mpsanity.Slug(input.Props.Slug[0])
	}

	if len(input.Props.Name) == 0 {
		doc.Type = "micropost"
	} else {
		doc.Type = "post"
		doc.Title = input.Props.Name[0]
		if doc.Slug == "" {
			doc.Slug = mpsanity.Slug(slug.Make(doc.Title))
		}
	}

	if len(input.Props.Content) > 0 {
		content := input.Props.Content[0]
		out, err := d.MarkdownConverter.ToBlocks(content)
		if err != nil {
			return nil, err
		}
		doc.Body = out

		if doc.Slug == "" {
			doc.Slug = mpsanity.Slug(slug.Make(block.ToPlainText(doc.Body)))
		}
	}

	// TODO photo

	if len(input.Props.Published) == 0 {
		doc.PublishedAt = time.Now()
	} else {
		doc.PublishedAt = input.Props.Published[0]
	}

	doc.Syndication = input.Props.Syndication

	return results, nil
}

type defaultDocument struct {
	Type        string        `json:"type"`
	Title       string        `json:"title,omitempty"`
	Body        []block.Block `json:"body"`
	Slug        mpsanity.Slug `json:"slug"`
	PublishedAt time.Time     `json:"publishedAt"`
	Syndication []string      `json:"syndication,omitempty"`
}
