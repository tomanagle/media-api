package handlers

import (
	"encoding/json"
	"go-media/internal/pkg/httperror"
	"go-media/internal/store"
	"net/http"
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

func (h *GetMediaHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	media, err := h.mediaStore.GetMedia(r.Context(), store.GetMediaParams{
		Skip:  0,
		Limit: 10,
	})

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error getting media: %s", err.Error())
		return
	}

	json.NewEncoder(w).Encode(media)
}
