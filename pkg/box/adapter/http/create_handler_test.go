package http

import (
	"context"
	"errors"
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
		serviceError     error
		expectedStatus   int
		expectedResponse string
	}{
		"no error": {
			file:           `{"content":"dGVzdDE=","path":"test.txt","mode":123}`,
			expectedStatus: http.StatusCreated,
		},
		"error unmarshal": {
			file:             `{not-good}`,
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `can't unmarshal body`,
		},
		"service error": {
			file:             `{"content":"dGVzdDE=","path":"test.txt","mode":755}`,
			serviceError:     errors.New("storage error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "internal error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mux := chi.NewMux()
			boxer := &boxerMock{
				WriteDocumentsFunc: func(ctx context.Context, file domain.File) error {
					return tc.serviceError
				},
			}
			handler := WriteDocumentHandler(boxer)
			mux.Post("/document", handler)

			req := httptest.NewRequest(http.MethodPost, "/document", strings.NewReader(tc.file))
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
