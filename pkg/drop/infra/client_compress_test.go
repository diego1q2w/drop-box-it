package infra

import (
	"errors"
	"fmt"
	"github.com/stretchr/testify/assert"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"
)

func TestClientCompress(t *testing.T) {
	testCases := map[string]struct {
		originalBody  io.Reader
		clientError   error
		expectedBody  string
		expectedError error
	}{
		"client error, error expected": {
			originalBody:  strings.NewReader(`{"test":2,"foo":"bar"}`),
			clientError:   errors.New("test"),
			expectedError: errors.New("test"),
			expectedBody:  "\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\xaaV*I-.Q\xb22\xd2QJ\xcb\xcfW\xb2RJJ,R\xaa\x05\x04\x00\x00\xff\xff\xb5bj\xa4\x16\x00\x00\x00",
		},
		"all cool": {
			originalBody: strings.NewReader(`{"test":2,"foo":"bar"}`),
			expectedBody: "\x1f\x8b\b\x00\x00\x00\x00\x00\x00\xff\xaaV*I-.Q\xb22\xd2QJ\xcb\xcfW\xb2RJJ,R\xaa\x05\x04\x00\x00\xff\xff\xb5bj\xa4\x16\x00\x00\x00",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			client := &mockClient{err: tc.clientError}
			clientCompress := NewClientCompress(client)

			req, err := http.NewRequest(http.MethodGet, "http://something.com", tc.originalBody)
			assert.NoError(t, err)

			_, err = clientCompress.Do(req)
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tc.expectedError) {
				t.Errorf("expected error: '%v', got: '%v'", tc.expectedError, err)
			}
			assert.Equal(t, string(client.doBody), tc.expectedBody)
		})
	}
}

type mockClient struct {
	err    error
	doBody []byte
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return &http.Response{}, err
	}

	m.doBody = body
	return nil, m.err
}
