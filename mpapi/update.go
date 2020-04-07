package mpapi

import (
	"encoding/json"
)

type UpdateInput struct {
	URL     string `json:"url"`
	Replace Props  `json:"replace"`
	Add     Props  `json:"add"`
	// Delete is either a []string or a Props
	Delete interface{} `json:"delete"`
}

// use a new type to get the standard unmarshalling behavior
type stdUpdateInput UpdateInput

func (input *UpdateInput) UnmarshalJSON(b []byte) error {
	var actualDelete interface{}

	var deleteKeys struct {
		Delete []string `json:"delete"`
	}
	if err := json.Unmarshal(b, &deleteKeys); err != nil {
		var deleteProps struct {
			Delete Props `json:"delete"`
		}
		if err := json.Unmarshal(b, &deleteProps); err != nil {
			return err
		}
		actualDelete = deleteProps
	} else {
		actualDelete = deleteKeys
	}

	stdInput := (*stdUpdateInput)(input)
	if err := json.Unmarshal(b, stdInput); err != nil {
		return err
	}
	input.Delete = actualDelete
	return nil
}
