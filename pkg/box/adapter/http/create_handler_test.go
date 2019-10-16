package http

import (
	"context"
	"errors"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateHandler(t *testing.T) {
	testCases := map[string]struct {
		file             string
		path             string
		serviceError     error
		expectedStatus   int
		expectedFile     domain.File
		expectedResponse string
	}{
		"no error": {
			path: "dGVzdC50eHQ=",
			file: `{"content":"dGVzdDE=","path":"test.txt","mode":123}`,
			expectedFile: domain.File{
				Path:    "test.txt",
				Mode:    123,
				Content: []byte("test1"),
			},
			expectedStatus: http.StatusCreated,
		},
		"error decoding": {
			path:             "---=-000",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `can't decode path`,
		},
		"error unmarshal": {
			path:             "dGVzdC50eHQ=",
			file:             `{not-good}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `can't unmarshal body`,
		},
		"service error": {
			path:           "dGVzdC50eHQ=",
			file:           `{"content":"dGVzdDE=","path":"test.txt","mode":755}`,
			serviceError:   errors.New("storage error"),
			expectedStatus: http.StatusInternalServerError,
			expectedFile: domain.File{
				Path:    "test.txt",
				Mode:    755,
				Content: []byte("test1"),
			},
			expectedResponse: "internal error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mux := chi.NewMux()
			boxer := &boxerMock{
				WriteDocumentsFunc: func(ctx context.Context, file domain.File) error {
					assert.Equal(t, tc.expectedFile, file)
					return tc.serviceError
				},
			}
			handler := WriteDocumentHandler(boxer)
			mux.Post("/document/{path}", handler)

			req := httptest.NewRequest(http.MethodPost, fmt.Sprintf("/document/%s", tc.path), strings.NewReader(tc.file))
			w := httptest.NewRecorder()
			mux.ServeHTTP(w, req)
			resp := w.Result()

			body, err := ioutil.ReadAll(resp.Body)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedStatus, resp.StatusCode)
			assert.Equal(t, tc.expectedResponse, strings.TrimSpace(string(body)))
		})
	}
}
