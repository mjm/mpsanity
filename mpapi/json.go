package mpapi

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

func (h *MicropubHandler) handleMicropubJSON(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleMicropubJSON")
	defer span.End()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		span.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	span.SetAttributes(contentLengthKey(len(data)))

	var typeVal struct {
		Action string   `json:"action"`
		Type   []string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeVal); err != nil {
		span.RecordError(ctx, err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if typeVal.Action == "update" {
		span.SetAttributes(requestTypeUpdate)
		// TODO handle update
	} else if len(typeVal.Type) > 0 {
		// create
		span.SetAttributes(requestTypeCreate, typeKey(typeVal.Type[0]))

		var input CreateInput
		if err := json.Unmarshal(data, &input); err != nil {
			span.RecordError(ctx, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		input.Type[0] = strings.TrimPrefix(input.Type[0], "h-")

		if len(input.Props.Photo) > 0 {
			rs, err := h.fetchImageAssets(ctx, input.Props.Photo)
			if err != nil {
				span.RecordError(ctx, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			imgIDs, err := h.uploadImageAssets(ctx, rs)
			if err != nil {
				span.RecordError(ctx, err)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			input.Props.Photo = imgIDs
		}

		doc, err := h.docBuilder.BuildDocument(ctx, &input)
		if err != nil {
			span.RecordError(ctx, err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		h.createDocument(ctx, w, doc)
	}
}

type CreateInput struct {
	Type  []string `json:"type"`
	Props struct {
		Name        []string    `json:"name,omitempty"`
		Content     []string    `json:"content,omitempty"`
		Slug        []string    `json:"mp-slug"`
		Published   []time.Time `json:"published"`
		Photo       []string    `json:"photo"`
		Syndication []string    `json:"syndication"`
	} `json:"properties"`
}
