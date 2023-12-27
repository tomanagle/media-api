package handlers

import (
	"bytes"
	idMock "go-media/internal/id/mock"
	storageMock "go-media/internal/storage/mock"
	"go-media/internal/store"
	storeMock "go-media/internal/store/mock"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestUploadFile(t *testing.T) {

	testCases := []struct {
		name string

		getVariations       func(ratio float64) []store.Variation
		createSessionResult *session.Session
		createSessionError  error
		uploadResult        *s3manager.UploadOutput
		uploadError         error

		expectedUploadFileName string
		mockID                 string

		createMediaResult   *store.Media
		createMediaError    error
		expectedCreateMedia store.CreateMediaParams

		expectedStatusCode int
		expectedResponse   []byte
	}{
		{
			name:               "success - resize and upload file",
			expectedStatusCode: http.StatusOK,
			getVariations: func(ratio float64) []store.Variation {
				return []store.Variation{
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
				}
			},
			createSessionResult: &session.Session{},
			createSessionError:  nil,
			uploadResult: &s3manager.UploadOutput{
				Location: "https://s3.amazonaws.com/go-media/123/123-small.png",
			},
			uploadError:            nil,
			expectedUploadFileName: "123/123-small.png",
			mockID:                 "123",
			expectedCreateMedia: store.CreateMediaParams{
				MediaID: "123",
				Variations: []store.Variation{
					{
						Name:     "small.png",
						Location: "https://s3.amazonaws.com/go-media/123/123-small.png",
						Width:    200,
						Height:   200,
					},
				},
			},
			createMediaResult: &store.Media{
				MediaID: "123",
				Variations: []store.Variation{
					{
						Name:     "small.png",
						Location: "https://s3.amazonaws.com/go-media/123/123-small.png",
						Width:    200,
						Height:   200,
					},
				},
			},
			createMediaError: nil,
			expectedResponse: []byte(`{"mediaId":"123","variations":[{"name":"small.png","location":"https://s3.amazonaws.com/go-media/123/123-small.png","width":200,"height":200}],"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}`),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			asset := assert.New(t)

			mockID := &idMock.MockID{}

			mockID.On("New").Return(tc.mockID, nil)

			mockStorage := &storageMock.MockStorage{}

			mockStorage.On("NewSession").Return(tc.createSessionResult, tc.createSessionError)
			mockStorage.On("Upload", tc.createSessionResult, mock.Anything, tc.expectedUploadFileName).Return(tc.uploadResult, tc.uploadError)

			mockMediaStore := &storeMock.MockMediaStore{}
			mockMediaStore.On("CreateMedia", tc.expectedCreateMedia).Return(tc.createMediaResult, tc.createMediaError)

			body := new(bytes.Buffer)
			mw := multipart.NewWriter(body)
			file, err := os.Open("./testdata/logo.png")
			if err != nil {
				t.Fatal(err)
			}

			w, err := mw.CreateFormFile("file", file.Name())

			if err != nil {
				t.Fatal(err)
			}

			if _, err := io.Copy(w, file); err != nil {
				t.Fatal(err)
			}

			mw.Close()

			request := httptest.NewRequest("POST", "/upload", body)
			request.Header.Set("Content-Type", mw.FormDataContentType())

			handler := NewUploadFileHandler(NewUploadFileHandlerParams{
				S3:            mockStorage,
				MediaStore:    mockMediaStore,
				ID:            mockID,
				GetVariations: tc.getVariations,
			})
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)

			respBody, err := io.ReadAll(response.Body)
			asset.NoError(err)

			asset.Equal(tc.expectedStatusCode, response.Code)
			asset.JSONEq(string(tc.expectedResponse), string(respBody))

			mockID.AssertExpectations(t)
			mockStorage.AssertExpectations(t)
			mockMediaStore.AssertExpectations(t)
		})
	}

}
