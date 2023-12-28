package handlers

import (
	"encoding/json"
	"go-media/internal/pkg/httperror"
	"go-media/internal/store"
	"net/http"
	"strconv"
)

type GetMediaHandler struct {
	mediaStore store.MediaStore
}

type NewGetMediaHandlerParams struct {
	MediaStore store.MediaStore
}

func NewGetMediaHandler(params NewGetMediaHandlerParams) *GetMediaHandler {
	return &GetMediaHandler{
		mediaStore: params.MediaStore,
	}
}

type getMediaResponse struct {
	Media []store.Media `json:"items"`
}

func (h *GetMediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	skipStr := r.URL.Query().Get("skip")
	limitStr := r.URL.Query().Get("limit")

	if skipStr == "" {
		skipStr = "0"
	}

	skip, err := strconv.Atoi(skipStr)

	if err != nil {
		httperror.Writef(w, http.StatusBadRequest, "error parsing skip: %s", err.Error())
		return
	}

	if limitStr == "" {
		limitStr = "20"
	}

	limit, err := strconv.Atoi(limitStr)

	if err != nil {
		httperror.Writef(w, http.StatusBadRequest, "error parsing limit: %s", err.Error())
		return
	}

	ctx := r.Context()
	media, err := h.mediaStore.GetMedia(ctx, store.GetMediaParams{
		Skip:  int64(skip),
		Limit: int64(limit),
	})

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error getting media: %s", err.Error())
		return
	}

	response := getMediaResponse{
		Media: media,
	}

	responseBytes, err := json.Marshal(response)

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error marshalling response: %s", err.Error())
		return
	}

	w.Write(responseBytes)
}
