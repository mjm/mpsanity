package mpapi

import (
	"time"
)

type Props struct {
	Name        []string    `json:"name,omitempty"`
	Content     []string    `json:"content,omitempty"`
	Slug        []string    `json:"mp-slug"`
	Published   []time.Time `json:"published"`
	Photo       []string    `json:"photo"`
	Syndication []string    `json:"syndication"`
}
