package mpsanity

import (
	"encoding/json"
	"errors"
	"fmt"
)

var ErrNotSlug = errors.New("value is not a slug")

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
