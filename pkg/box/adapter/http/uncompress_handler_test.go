package http

import (
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestUncompressHandler(t *testing.T) {
	testCases := map[string]struct {
		body          io.Reader
		contentHeader string
		expectedBody  string
	}{
		"if content header is set it should wrap the reader": {
			body:          strings.NewReader("\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\xaaVJ\xce\xcf+I\xcd+Q\xb2RJq\x0f\xabJqv\xaaJ22\xccIq\xcf(H\xcaM\xb6U\xd2Q\xca\xcdOIU\xb2254\xac\x05\x04\x00\x00\xff\xffg\u007f\xd4|-\x00\x00\x00"),
			contentHeader: "gzip",
			expectedBody:  "{\"content\":\"dGVzdCBzb21ldGhpbmc=\",\"mode\":511}",
		},
		"if content header is an unknown option should pass the body": {
			body:          strings.NewReader("something cool"),
			contentHeader: "nothing",
			expectedBody:  "something cool",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {

			targetHandler := func(writer http.ResponseWriter, request *http.Request) {
				body, err := ioutil.ReadAll(request.Body)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedBody, string(body))
			}

			uHandler := UncompressHandler(targetHandler)
			req, err := http.NewRequest("GET", "http://something.com", tc.body)
			assert.NoError(t, err)
			req.Header.Set("Content-Encoding", "gzip")

			w := httptest.NewRecorder()
			uHandler(w, req)
		})
	}
}
