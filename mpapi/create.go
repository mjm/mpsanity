package mpapi

import (
	"time"

	"github.com/mjm/mpsanity"
)

type CreateInput struct {
	Type  []string `json:"type"`
	Props Props    `json:"properties"`
}

func (in *CreateInput) Name() string {
	if vs := in.Props.Name; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (in *CreateInput) Content() string {
	if vs := in.Props.Content; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (in *CreateInput) Slug() string {
	if vs := in.Props.Slug; len(vs) > 0 {
		return vs[0]
	}
	return ""
}

func (in *CreateInput) Published() *time.Time {
	if vs := in.Props.Published; len(vs) > 0 {
		return &vs[0]
	}
	return nil
}

func (in *CreateInput) Photos() []mpsanity.Reference {
	var refs []mpsanity.Reference
	for _, p := range in.Props.Photo {
		refs = append(refs, mpsanity.Reference(p))
	}
	return refs
}

func (in *CreateInput) Syndication() []string {
	return in.Props.Syndication
}
