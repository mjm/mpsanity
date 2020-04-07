package mpapi

import (
	"context"
	"time"

	"github.com/gosimple/slug"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/mjm/mpsanity"
	"github.com/mjm/mpsanity/block"
	"github.com/mjm/mpsanity/patch"
)

var ErrNotEntry = status.Error(codes.InvalidArgument, "post is not an entry")

type DocumentBuilder interface {
	BuildDocument(ctx context.Context, input *CreateInput) (Document, error)
	UpdateDocument(ctx context.Context, input *UpdateInput) ([]patch.Patch, error)
}

type DefaultDocumentBuilder struct {
	MarkdownConverter *block.MarkdownConverter
}

func (d *DefaultDocumentBuilder) BuildDocument(_ context.Context, input *CreateInput) (Document, error) {
	if input.Type[0] != "entry" {
		return nil, ErrNotEntry
	}

	var doc defaultDocument

	if slug := input.Slug(); slug != "" {
		doc.Slug = mpsanity.Slug(slug)
	}

	if name := input.Name(); name != "" {
		doc.Type = "post"
		doc.Title = name
		if doc.Slug == "" {
			doc.Slug = mpsanity.Slug(slug.Make(doc.Title))
		}
	} else {
		doc.Type = "micropost"
	}

	if content := input.Content(); content != "" {
		out, err := d.MarkdownConverter.ToBlocks(content)
		if err != nil {
			return nil, err
		}
		doc.Body = out

		if doc.Slug == "" {
			doc.Slug = mpsanity.Slug(slug.Make(block.ToPlainText(doc.Body)))
		}
	}

	if doc.Slug == "" {
		doc.Slug = mpsanity.Slug(randomString(10))
	}

	for _, photo := range input.Photos() {
		doc.Body = append(doc.Body, block.Block{
			Type: "mainImage",
			Content: map[string]interface{}{
				"alt":   "Photo",
				"asset": photo,
			},
		})
	}

	if pub := input.Published(); pub != nil {
		doc.PublishedAt = *pub
	} else {
		doc.PublishedAt = time.Now()
	}

	doc.Slug = mpsanity.Slug(doc.PublishedAt.Format("2006-01-02") + "-" + string(doc.Slug))

	doc.Syndication = input.Syndication()

	return &doc, nil
}

func (d *DefaultDocumentBuilder) UpdateDocument(ctx context.Context, input *UpdateInput) ([]patch.Patch, error) {
	var ps []patch.Patch

	if len(input.Replace.Name) > 0 {
		ps = append(ps, patch.Set("title", input.Replace.Name[0]))
	}
	// TODO update content (this is kinda hard when photos are stored inside it)
	if len(input.Replace.Syndication) > 0 {
		ps = append(ps, patch.Set("syndication", input.Replace.Syndication))
	}

	if len(input.Add.Syndication) > 0 {
		var items []interface{}
		for _, u := range input.Add.Syndication {
			items = append(items, u)
		}
		ps = append(ps,
			patch.SetIfMissing("syndication", make([]string, 0)),
			patch.InsertAfter("syndication[-1]", items...))
	}

	return ps, nil
}

type defaultDocument struct {
	Type        string        `json:"_type"`
	Title       string        `json:"title,omitempty"`
	Body        []block.Block `json:"body"`
	Slug        mpsanity.Slug `json:"slug"`
	PublishedAt time.Time     `json:"publishedAt"`
	Syndication []string      `json:"syndication,omitempty"`
}

type Document interface {
	URLPath() string
}

func (d *defaultDocument) URLPath() string {
	return "/" + string(d.Slug)
}
