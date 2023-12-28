package handlers

import (
	"go-media/internal/store"
	storeMock "go-media/internal/store/mock"
	"io"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetMed(t *testing.T) {

	testCases := []struct {
		name                   string
		expectedResponse       []byte
		expectedGetMediaParams store.GetMediaParams
		getMediaResult         []store.Media
		getMediaError          error
	}{
		{
			name:             "success - get media with default pagination",
			expectedResponse: []byte(`{"items":[{"mediaId":"1","variations":[{"name":"small","location":"","width":0,"height":0}],"createdAt":"0001-01-01T00:00:00Z","updatedAt":"0001-01-01T00:00:00Z"}]}`),
			getMediaResult: []store.Media{
				{
					MediaID: "1",
					Variations: []store.Variation{
						{
							Name: "small",
						},
					},
				},
			},
			expectedGetMediaParams: store.GetMediaParams{
				Skip:  0,
				Limit: 20,
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {

			asset := assert.New(t)
			mockMediaStore := &storeMock.MockMediaStore{}

			mockMediaStore.On("GetMedia", mock.Anything, tc.expectedGetMediaParams).Return(tc.getMediaResult, tc.getMediaError)

			request := httptest.NewRequest("GET", "/media", nil)

			handler := NewGetMediaHandler(NewGetMediaHandlerParams{
				MediaStore: mockMediaStore,
			})
			response := httptest.NewRecorder()
			handler.ServeHTTP(response, request)
			respBody, err := io.ReadAll(response.Body)
			asset.NoError(err)
			asset.JSONEq(string(tc.expectedResponse), string(respBody))

			mockMediaStore.AssertExpectations(t)
		})
	}

}
