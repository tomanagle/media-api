package handlers

import (
	"fmt"
	"go-media/internal/pkg/httperror"
	"go-media/internal/storage/s3"
	"go-media/internal/store"
	"io"
	"net/http"
	"os"

	"github.com/google/uuid"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
)

type UploadFileHandler struct {
	s3         *s3.S3
	mediaStore store.MediaStore
}

type NewUploadFileHandlerParams struct {
	S3         *s3.S3
	MediaStore store.MediaStore
}

func NewUploadFileHandler(params NewUploadFileHandlerParams) *UploadFileHandler {
	return &UploadFileHandler{
		s3:         params.S3,
		mediaStore: params.MediaStore,
	}
}

func (h *UploadFileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")

	if err != nil {
		httperror.Writef(w, http.StatusBadRequest, "error reading file: %s", err.Error())
		return
	}
	defer file.Close()

	id := uuid.New()

	// make a folder for this image and it's variants
	err = os.MkdirAll(id.String(), os.ModePerm)

	if err != nil && !os.IsExist(err) {
		httperror.Writef(w, http.StatusInternalServerError, "error creating folder: %s", err.Error())
		return
	}

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

	ratio := float64(width) / float64(height)

	variations := []store.Variation{
		{
			Name:  "small",
			Width: 200,
			Height: func() int {
				if ratio > 1 {
					return int(200 / ratio)
				}
				return int(200 * ratio)
			}(),
		},
		{
			Name:  "medium",
			Width: 400,
			Height: func() int {
				if ratio > 1 {
					return int(400 / ratio)
				}
				return int(400 * ratio)
			}(),
		},
		{
			Name:  "large",
			Width: 600,
			Height: func() int {
				if ratio > 1 {
					return int(600 / ratio)
				}
				return int(600 * ratio)
			}(),
		},
	}

	for _, v := range variations {
		resized := transform.Resize(img, v.Width, v.Height, transform.Linear)

		if err := imgio.Save(fmt.Sprintf("%s/%s-%s.png", id, id, v.Name), resized, imgio.PNGEncoder()); err != nil {
			fmt.Println(err)
			return
		}
	}

	session, err := h.s3.NewSession()

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error creating session: %s", err.Error())
		return
	}

	for index, v := range variations {
		file, err := os.Open(fmt.Sprintf("%s/%s-%s.png", id, id, v.Name))

		if err != nil {
			httperror.Writef(w, http.StatusInternalServerError, "error opening file: %s", err.Error())
			return
		}

		defer file.Close()

		res, err := h.s3.Upload(session, file, fmt.Sprintf("%s/%s-%s.png", id, id, v.Name))

		if err != nil {
			httperror.Writef(w, http.StatusInternalServerError, "error uploading file: %s", err.Error())
			return
		}

		variations[index].Name = fmt.Sprintf("%s.png", v.Name)
		variations[index].Location = res.Location
	}

	err = os.RemoveAll(id.String())

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error removing folder: %s", err.Error())
		return
	}

	_, err = h.mediaStore.CreateMedia(store.CreateMediaParams{
		MediaID:    id.String(),
		Variations: variations,
	})

	if err != nil {
		httperror.Writef(w, http.StatusInternalServerError, "error creating media: %s", err.Error())
		return
	}

	w.Write([]byte("OK"))
}
