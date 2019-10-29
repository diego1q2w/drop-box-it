package infra

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"net/http"
	"net/url"
	"os"
)

type Client interface {
	Do(req *http.Request) (*http.Response, error)
}

type BoxClient struct {
	requester Client
	baseURL   string
}

//go:generate moq -out box_client_mock_test.go . BoxClienter
type BoxClienter interface {
	WriteDocument(ctx context.Context, file domain.File, content []byte) error
	DeleteDocument(ctx context.Context, file domain.File, content []byte) error
}

type createFileClient struct {
	Content []byte      `json:"content"`
	Mode    os.FileMode `json:"mode"`
}

func NewBoxClient(requester Client, baseURL string) BoxClienter {
	return &BoxClient{
		requester: requester,
		baseURL:   baseURL,
	}
}

func (c *BoxClient) WriteDocument(ctx context.Context, file domain.File, content []byte) error {
	uri, err := c.getURI(file)
	if err != nil {
		return err
	}

	createFile := createFileClient{
		Content: content,
		Mode:    file.Mode,
	}
	body, err := json.Marshal(createFile)
	if err != nil {
		return fmt.Errorf("error while marshalling body: %w", err)
	}

	req, err := http.NewRequest(http.MethodPost, uri, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	res, err := c.requester.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error while requesting: %w", err)
	}

	if res.StatusCode != http.StatusCreated {
		return fmt.Errorf("error writing documents got status %d", res.StatusCode)
	}

	return nil
}

func (c *BoxClient) DeleteDocument(ctx context.Context, file domain.File, content []byte) error {
	uri, err := c.getURI(file)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(http.MethodDelete, uri, nil)
	if err != nil {
		return fmt.Errorf("error while creating request: %w", err)
	}

	res, err := c.requester.Do(req.WithContext(ctx))
	if err != nil {
		return fmt.Errorf("error while requesting: %w", err)
	}

	if res.StatusCode != http.StatusNoContent {
		return fmt.Errorf("error deleting documents got status %d", res.StatusCode)
	}

	return nil
}

func (c *BoxClient) getURI(file domain.File) (string, error) {
	base64Path := base64.StdEncoding.EncodeToString(file.Path.ToBytes())
	uri, err := c.withPath(fmt.Sprintf("/document/%s", base64Path))

	return uri, err
}

func (c *BoxClient) withPath(path string) (string, error) {
	uri, err := url.Parse(c.baseURL)
	if err != nil {
		return "", fmt.Errorf("unable to parse string: %w", err)
	}

	uri.Path = path
	return uri.String(), nil
}
