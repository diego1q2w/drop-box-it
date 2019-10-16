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

func TestDeleteHandler(t *testing.T) {
	testCases := map[string]struct {
		path             string
		serviceError     error
		expectedStatus   int
		expectedPath     domain.Path
		expectedResponse string
	}{
		"no error": {
			path:           "dGVzdDE=",
			expectedPath:   "test1",
			expectedStatus: http.StatusNoContent,
		},
		"error decoding": {
			path:             "00000----",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: `can't decode path`,
		},
		"service error": {
			path:             "dGVzdDE=",
			serviceError:     errors.New("storage error"),
			expectedStatus:   http.StatusInternalServerError,
			expectedPath:     "test1",
			expectedResponse: "internal error",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			mux := chi.NewMux()
			boxer := &boxerMock{
				DeleteDocumentsFunc: func(ctx context.Context, path domain.Path) error {
					assert.Equal(t, tc.expectedPath, path)
					return tc.serviceError
				},
			}
			handler := DeleteDocumentHandler(boxer)
			mux.Delete("/document/{path}", handler)

			req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/document/%s", tc.path), nil)
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
