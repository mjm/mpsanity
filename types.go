package mpsanity

import (
	"encoding/json"
	"errors"
	"fmt"
)

var (
	ErrNotSlug      = errors.New("value is not a slug")
	ErrNotReference = errors.New("value is not a reference")
)

type Slug string

func (s Slug) MarshalJSON() ([]byte, error) {
	slug := internalSlug{
		Type:    "slug",
		Current: string(s),
	}
	return json.Marshal(slug)
}

func (s *Slug) UnmarshalJSON(b []byte) error {
	var slug internalSlug
	if err := json.Unmarshal(b, &slug); err != nil {
		return err
	}

	if slug.Type != "slug" {
		return fmt.Errorf("%s %w", slug.Type, ErrNotSlug)
	}

	*s = Slug(slug.Current)
	return nil
}

type internalSlug struct {
	Type    string `json:"_type"`
	Current string `json:"current"`
}

type Reference string

func (r Reference) MarshalJSON() ([]byte, error) {
	ref := internalReference{
		Type: "reference",
		Ref:  string(r),
	}
	return json.Marshal(ref)
}

func (r *Reference) UnmarshalJSON(b []byte) error {
	var ref internalReference
	if err := json.Unmarshal(b, &ref); err != nil {
		return err
	}

	if ref.Type != "reference" {
		return fmt.Errorf("%s %w", ref.Type, ErrNotReference)
	}

	*r = Reference(ref.Ref)
	return nil
}

type internalReference struct {
	Type string `json:"_type"`
	Ref  string `json:"_ref"`
}
