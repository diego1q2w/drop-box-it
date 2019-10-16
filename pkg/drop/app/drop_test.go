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
		initialFileStatus   map[domain.Path]domain.FileStatus
		expectedFinalStatus map[domain.Path]domain.FileStatus
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
			initialFileStatus: map[domain.Path]domain.FileStatus{
				"test1": domain.Created,
				"test5": domain.Created,
			},
			expectedFinalStatus: map[domain.Path]domain.FileStatus{
				"test1": domain.Updated,
				"test3": domain.Created,
				"test4": domain.Created,
				"test5": domain.Deleted,
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
			dropper := NewDropper(fetcher)
			dropper.filesStatus = tc.initialFileStatus

			err := dropper.SyncFiles(context.Background(), "root")
			if fmt.Sprintf("%v", err) != fmt.Sprintf("%v", tc.expectedError) {
				t.Errorf("exepected error: %s, got: %s", tc.expectedError, err)
			}

			assert.Equal(t, dropper.filesStatus, tc.expectedFinalStatus)
		})
	}
}
