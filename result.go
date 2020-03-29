package mpsanity

import (
	"encoding/json"
)

type queryResult struct {
	Duration float64          `json:"ms"`
	Query    string           `json:"query"`
	Result   *json.RawMessage `json:"result"`
}
