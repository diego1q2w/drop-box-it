package infra

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ClientCompress struct {
	client Client
}

func NewClientCompress(client Client) *ClientCompress {
	return &ClientCompress{client: client}
}

func (c *ClientCompress) Do(req *http.Request) (*http.Response, error) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, fmt.Errorf("unable to read body: %w", err)
	}

	body, err = c.compress(body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Encoding", "gzip")

	req, err = http.NewRequest(req.Method, req.URL.String(), bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("error while creating compress request: %w", err)
	}

	return c.client.Do(req)
}

func (c *ClientCompress) compress(content []byte) ([]byte, error) {
	var contentWriter bytes.Buffer
	gz := gzip.NewWriter(&contentWriter)
	if _, err := gz.Write(content); err != nil {
		return nil, fmt.Errorf("unable to compress body: %w", err)
	}

	gz.Close()

	return contentWriter.Bytes(), nil
}
