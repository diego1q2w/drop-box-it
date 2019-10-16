package infra

import (
	"context"
	"errors"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestWriteDocumentClient_NoError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	boxClient := NewBoxClient(http.DefaultClient, ts.URL)
	err := boxClient.WriteDocument(context.Background(), domain.File{
		Path: "path/text1.txt",
		Mode: os.ModePerm,
	}, []byte("test something"))
	assert.NoError(t, err)
}

func TestWriteDocumentClient_WithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	boxClient := NewBoxClient(http.DefaultClient, ts.URL)
	err := boxClient.WriteDocument(context.Background(), domain.File{
		Path: "path/text1.txt",
		Mode: os.ModePerm,
	}, []byte("test something"))
	expectedError := errors.New("error writing documents got status 400")

	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestWriteDocumentClient_BadUrl(t *testing.T) {
	boxClient := NewBoxClient(http.DefaultClient, "this-is-wrong")
	err := boxClient.WriteDocument(context.Background(), domain.File{
		Path: "path/text1.txt",
		Mode: os.ModePerm,
	}, []byte("test something"))
	expectedError := errors.New(`error while requesting: Post /document/cGF0aC90ZXh0MS50eHQ=: unsupported protocol scheme ""`)

	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}

func TestWriteDocumentClient_IsTheRequestWellFormed(t *testing.T) {
	expectedURL := "/document/cGF0aC90ZXh0MS50eHQ="
	expectedMethod := "POST"
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, r.Method, expectedMethod)
		assert.Equal(t, r.URL.Path, expectedURL)
		w.WriteHeader(http.StatusCreated)
	}))
	defer ts.Close()

	boxClient := NewBoxClient(http.DefaultClient, ts.URL)
	err := boxClient.WriteDocument(context.Background(), domain.File{
		Path: "path/text1.txt",
		Mode: os.ModePerm,
	}, []byte("test something"))
	assert.NoError(t, err)
}

func TestDeleteDocumentClient_NoError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	}))
	defer ts.Close()

	boxClient := NewBoxClient(http.DefaultClient, ts.URL)
	err := boxClient.DeleteDocument(context.Background(), domain.File{
		Path: "path/text1.txt",
		Mode: os.ModePerm,
	}, []byte("test something"))
	assert.NoError(t, err)
}

func TestDeleteDocumentClient_WithError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer ts.Close()

	boxClient := NewBoxClient(http.DefaultClient, ts.URL)
	err := boxClient.DeleteDocument(context.Background(), domain.File{
		Path: "path/text1.txt",
		Mode: os.ModePerm,
	}, []byte("test something"))
	expectedError := errors.New("error deleting documents got status 400")

	if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", expectedError) {
		t.Errorf("expected error %v, got %v", expectedError, err)
	}
}
