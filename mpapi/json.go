package mpapi

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel/api/key"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNoURL     = status.Error(codes.InvalidArgument, "no url provided for update")
	ErrWrongBase = status.Error(codes.InvalidArgument, "url is not from this site")
)

func (h *MicropubHandler) handleMicropubJSON(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "handleMicropubJSON")
	defer span.End()

	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		respondWithError(ctx, w, err)
		return
	}

	span.SetAttributes(contentLengthKey(len(data)))

	var typeVal struct {
		Action string   `json:"action"`
		Type   []string `json:"type"`
	}
	if err := json.Unmarshal(data, &typeVal); err != nil {
		err = status.Error(codes.InvalidArgument, err.Error())
		respondWithError(ctx, w, err)
		return
	}

	if typeVal.Action == "update" {
		span.SetAttributes(requestTypeUpdate)

		var input UpdateInput
		if err := json.Unmarshal(data, &input); err != nil {
			err = status.Error(codes.InvalidArgument, err.Error())
			respondWithError(ctx, w, err)
			return
		}

		span.SetAttributes(urlKey(input.URL))

		if input.URL == "" {
			respondWithError(ctx, w, ErrNoURL)
			return
		}

		if !strings.HasPrefix(input.URL, h.baseURL+"/") {
			respondWithError(ctx, w, ErrWrongBase)
		}

		patches, err := h.docBuilder.UpdateDocument(ctx, &input)
		if err != nil {
			respondWithError(ctx, w, err)
			return
		}

		// TODO maybe move query construction into document builder
		slug := strings.TrimPrefix(strings.TrimSuffix(input.URL, "/"), h.baseURL+"/")
		span.SetAttributes(slugKey(slug))
		q := fmt.Sprintf(`*[slug.current == %q]`, slug)
		span.SetAttributes(key.String("sanity.query", q))
		if err := h.Sanity.Txn().PatchQuery(q, patches...).Commit(ctx); err != nil {
			respondWithError(ctx, w, err)
			return
		}

		notifyTitle := fmt.Sprintf("Update %s", slug)
		h.notifyWebhook(ctx, notifyTitle)

		w.WriteHeader(http.StatusNoContent)
	} else if len(typeVal.Type) > 0 {
		// create
		span.SetAttributes(requestTypeCreate, typeKey(typeVal.Type[0]))

		var input CreateInput
		if err := json.Unmarshal(data, &input); err != nil {
			err = status.Error(codes.InvalidArgument, err.Error())
			respondWithError(ctx, w, err)
			return
		}

		input.Type[0] = strings.TrimPrefix(input.Type[0], "h-")

		if len(input.Props.Photo) > 0 {
			rs, err := h.fetchImageAssets(ctx, input.Props.Photo)
			if err != nil {
				respondWithError(ctx, w, err)
				return
			}

			imgIDs, err := h.uploadImageAssets(ctx, rs)
			if err != nil {
				respondWithError(ctx, w, err)
				return
			}

			input.Props.Photo = imgIDs
		}

		h.createDocument(ctx, w, &input)
	}
}
