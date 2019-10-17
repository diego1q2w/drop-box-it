package app

import (
	"context"
	"fmt"
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSyncFiles(t *testing.T) {
	testCases := map[string]struct {
		existsPath          bool
		isDir               bool
		listFiles           []domain.File
		listingErr          error
		initialFileStatus   map[domain.Path]fileStatus
		expectedFinalStatus map[domain.Path]fileStatus
		expectedError       error
	}{
		"if path does not exists error expected": {
			existsPath:    false,
			expectedError: fmt.Errorf("the path 'root' does not exists"),
		},
		"if path is not a directory error expected": {
			existsPath:    true,
			isDir:         false,
			expectedError: fmt.Errorf("the path priovided is not a directory"),
		},
		"if error listing error expected": {
			existsPath:    true,
			isDir:         true,
			listingErr:    fmt.Errorf("test"),
			expectedError: fmt.Errorf("error while listing files: test"),
		},
		"should update the state accordingly": {
			existsPath: true,
			isDir:      true,
			listFiles: []domain.File{
				{Path: "test1"},
				{Path: "test3"},
				{Path: "test4"},
			},
			initialFileStatus: map[domain.Path]fileStatus{
				"test1": {
					status: domain.Created,
				},
				"test5": {
					status: domain.Created,
				},
			},
			expectedFinalStatus: map[domain.Path]fileStatus{
				"test1": {
					status: domain.Updated,
					file:   domain.File{Path: "test1"},
				},
				"test3": {
					status: domain.Created,
					file:   domain.File{Path: "test3"},
				},
				"test4": {
					status: domain.Created,
					file:   domain.File{Path: "test4"},
				},
				"test5": {
					status: domain.Deleted,
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			fetcher := &fileFetcherMock{
				ListFilesFunc: func(root domain.Path) (files []domain.File, e error) {
					return tc.listFiles, tc.listingErr
				},
				PathExistsFunc: func(path domain.Path) (b bool, b2 bool) {
					return tc.existsPath, tc.isDir
				},
			}

			boxClient := &boxClientMock{}
			dropper := NewDropper(fetcher, boxClient, "root")
			dropper.filesStatus = tc.initialFileStatus

			err := dropper.SyncFiles(context.Background())
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tc.expectedError) {
				t.Errorf("exepected error: %s, got: %s", tc.expectedError, err)
			}

			assert.Equal(t, tc.expectedFinalStatus, dropper.filesStatus)
		})
	}
}

func TestSendUpdatesFinalStateCorrect(t *testing.T) {
	files := map[domain.Path]fileStatus{
		"test1": {
			status: domain.Updated,
			file:   domain.File{Path: "test1"},
		},
		"test3": {
			status: domain.Created,
			file:   domain.File{Path: "test3"},
		},
		"test4": {
			status: domain.Created,
			file:   domain.File{Path: "test4"},
		},
		"test5": {
			status: domain.Deleted,
			file:   domain.File{Path: "test5"},
		},
		"test6": {
			status: domain.Deleted,
			file:   domain.File{Path: "test6"},
		},
	}
	expected := map[domain.Path]fileStatus{
		"test1": {
			status: domain.Synced,
			file:   domain.File{Path: "test1"},
		},
		"test3": {
			status: domain.Synced,
			file:   domain.File{Path: "test3"},
		},
		"test4": {
			status: domain.Synced,
			file:   domain.File{Path: "test4"},
		},
	}
	fetcher := &fileFetcherMock{
		ReadFileContentFunc: func(path domain.Path) (bytes []byte, e error) {
			return []byte(`test`), nil
		},
	}

	boxClient := &boxClientMock{
		DeleteDocumentFunc: func(ctx context.Context, file domain.File, content []byte) error {
			return nil
		},
		WriteDocumentFunc: func(ctx context.Context, file domain.File, content []byte) error {
			return nil
		},
	}
	dropper := NewDropper(fetcher, boxClient, "")
	dropper.filesStatus = files

	dropper.SendUpdates(context.Background())
	assert.Equal(t, expected, dropper.filesStatus)
}
