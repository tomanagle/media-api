package handlers

import (
	"encoding/json"
	"fmt"
	"go-media/internal/id"
	"go-media/internal/pkg/httperror"
	"go-media/internal/storage"
	"go-media/internal/store"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type UploadFileHandler struct {
	s3            storage.Storage
	id            id.ID
	mediaStore    store.MediaStore
	getVariations func(ratio float64) []store.Variation
}

type NewUploadFileHandlerParams struct {
	S3            storage.Storage
	ID            id.ID
	MediaStore    store.MediaStore
	GetVariations func(ratio float64) []store.Variation
}

func NewUploadFileHandler(params NewUploadFileHandlerParams) *UploadFileHandler {
	return &UploadFileHandler{
		s3:            params.S3,
		id:            params.ID,
		mediaStore:    params.MediaStore,
		getVariations: params.GetVariations,
	}
}

func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")

	if err != nil {
		httperror.Writef(w, http.StatusBadRequest, "error reading file: %s", err.Error())
		return
	}
	defer file.Close()

	id, err := h.id.New()

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error generating id: %s", err.Error())
		return
	}

	// make a folder for this image and it's variants
	err = os.MkdirAll(id, os.ModePerm)

	if err != nil && !os.IsExist(err) {
		httperror.Writef(w, http.StatusInternalServerError, "error creating folder: %s", err.Error())
		return
	}

	// cleanup the local files
	defer func() {
		parentDir := id
		err = os.RemoveAll(parentDir)
		if err != nil {
			fmt.Println("error deleting folder", err)
		}
	}()

	fileName := header.Filename
	filePath := fmt.Sprintf("%s/%s", id, fileName)

	// write the file to disk
	f, err := os.Create(filePath)

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error writing file: %s", err.Error())
		return
	}

	defer f.Close()

	_, err = io.Copy(f, file)

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error copying file: %s", err.Error())
		return
	}

	img, err := imgio.Open(filePath)
	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error opening file: %s", err.Error())
		return
	}

	width := img.Bounds().Dx()
	height := img.Bounds().Dy()

	// find the image ratio so we don't make funny shapes when we resize
	ratio := float64(width) / float64(height)

	variations := h.getVariations(ratio)

	resultsChan := make(chan store.Variation, len(variations))
	errorChan := make(chan error, len(variations))

	wg := sync.WaitGroup{}

	for _, v := range variations {
		wg.Add(1)

		go func(v store.Variation) {
			defer wg.Done()
			resized := transform.Resize(img, v.Width, v.Height, transform.Linear)

			err := imgio.Save(fmt.Sprintf("%s/%s-%s.png", id, id, v.Name), resized, imgio.PNGEncoder())

			if err != nil {
				errorChan <- err
				return
			}

			file, err := os.Open(fmt.Sprintf("%s/%s-%s.png", id, id, v.Name))

			if err != nil {
				errorChan <- err
				return
			}
			defer file.Close()
			session, err := h.s3.NewSession()
			if err != nil {
				errorChan <- err
			}

			res, err := h.s3.Upload(session, file, fmt.Sprintf("%s/%s-%s.png", id, id, v.Name))

			if err != nil {
				errorChan <- err
			}

			v.Name = fmt.Sprintf("%s.png", v.Name)
			v.Location = res.Location
			resultsChan <- v

		}(v)

	}

	go func() {
		wg.Wait()
		close(resultsChan)
		close(errorChan)
	}()

	uploadVariations := []store.Variation{}
	for result := range resultsChan {
		uploadVariations = append(uploadVariations, result)
	}

	media, err := h.mediaStore.CreateMedia(r.Context(), store.CreateMediaParams{
		MediaID:    id,
		Variations: uploadVariations,
	})

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error creating media: %s", err.Error())
		return
	}

	if err := json.NewEncoder(w).Encode(media); err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error encoding media: %s", err.Error())
		return
	}

}
