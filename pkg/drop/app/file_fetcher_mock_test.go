// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app

import (
	"github.com/diego1q2w/drop-box-it/pkg/drop/domain"
	"sync"
)

var (
	lockfileFetcherMockListFiles       sync.RWMutex
	lockfileFetcherMockPathExists      sync.RWMutex
	lockfileFetcherMockReadFileContent sync.RWMutex
)

// Ensure, that fileFetcherMock does implement fileFetcher.
// If this is not the case, regenerate this file with moq.
var _ fileFetcher = &fileFetcherMock{}

// fileFetcherMock is a mock implementation of fileFetcher.
//
//     func TestSomethingThatUsesfileFetcher(t *testing.T) {
//
//         // make and configure a mocked fileFetcher
//         mockedfileFetcher := &fileFetcherMock{
//             ListFilesFunc: func(root domain.Path) ([]domain.File, error) {
// 	               panic("mock out the ListFiles method")
//             },
//             PathExistsFunc: func(path domain.Path) (bool, bool) {
// 	               panic("mock out the PathExists method")
//             },
//             ReadFileContentFunc: func(path domain.Path) ([]byte, error) {
// 	               panic("mock out the ReadFileContent method")
//             },
//         }
//
//         // use mockedfileFetcher in code that requires fileFetcher
//         // and then make assertions.
//
//     }
type fileFetcherMock struct {
	// ListFilesFunc mocks the ListFiles method.
	ListFilesFunc func(root domain.Path) ([]domain.File, error)

	// PathExistsFunc mocks the PathExists method.
	PathExistsFunc func(path domain.Path) (bool, bool)

	// ReadFileContentFunc mocks the ReadFileContent method.
	ReadFileContentFunc func(path domain.Path) ([]byte, error)

	// calls tracks calls to the methods.
	calls struct {
		// ListFiles holds details about calls to the ListFiles method.
		ListFiles []struct {
			// Root is the root argument value.
			Root domain.Path
		}
		// PathExists holds details about calls to the PathExists method.
		PathExists []struct {
			// Path is the path argument value.
			Path domain.Path
		}
		// ReadFileContent holds details about calls to the ReadFileContent method.
		ReadFileContent []struct {
			// Path is the path argument value.
			Path domain.Path
		}
	}
}

// ListFiles calls ListFilesFunc.
func (mock *fileFetcherMock) ListFiles(root domain.Path) ([]domain.File, error) {
	if mock.ListFilesFunc == nil {
		panic("fileFetcherMock.ListFilesFunc: method is nil but fileFetcher.ListFiles was just called")
	}
	callInfo := struct {
		Root domain.Path
	}{
		Root: root,
	}
	lockfileFetcherMockListFiles.Lock()
	mock.calls.ListFiles = append(mock.calls.ListFiles, callInfo)
	lockfileFetcherMockListFiles.Unlock()
	return mock.ListFilesFunc(root)
}

// ListFilesCalls gets all the calls that were made to ListFiles.
// Check the length with:
//     len(mockedfileFetcher.ListFilesCalls())
func (mock *fileFetcherMock) ListFilesCalls() []struct {
	Root domain.Path
} {
	var calls []struct {
		Root domain.Path
	}
	lockfileFetcherMockListFiles.RLock()
	calls = mock.calls.ListFiles
	lockfileFetcherMockListFiles.RUnlock()
	return calls
}

// PathExists calls PathExistsFunc.
func (mock *fileFetcherMock) PathExists(path domain.Path) (bool, bool) {
	if mock.PathExistsFunc == nil {
		panic("fileFetcherMock.PathExistsFunc: method is nil but fileFetcher.PathExists was just called")
	}
	callInfo := struct {
		Path domain.Path
	}{
		Path: path,
	}
	lockfileFetcherMockPathExists.Lock()
	mock.calls.PathExists = append(mock.calls.PathExists, callInfo)
	lockfileFetcherMockPathExists.Unlock()
	return mock.PathExistsFunc(path)
}

// PathExistsCalls gets all the calls that were made to PathExists.
// Check the length with:
//     len(mockedfileFetcher.PathExistsCalls())
func (mock *fileFetcherMock) PathExistsCalls() []struct {
	Path domain.Path
} {
	var calls []struct {
		Path domain.Path
	}
	lockfileFetcherMockPathExists.RLock()
	calls = mock.calls.PathExists
	lockfileFetcherMockPathExists.RUnlock()
	return calls
}

// ReadFileContent calls ReadFileContentFunc.
func (mock *fileFetcherMock) ReadFileContent(path domain.Path) ([]byte, error) {
	if mock.ReadFileContentFunc == nil {
		panic("fileFetcherMock.ReadFileContentFunc: method is nil but fileFetcher.ReadFileContent was just called")
	}
	callInfo := struct {
		Path domain.Path
	}{
		Path: path,
	}
	lockfileFetcherMockReadFileContent.Lock()
	mock.calls.ReadFileContent = append(mock.calls.ReadFileContent, callInfo)
	lockfileFetcherMockReadFileContent.Unlock()
	return mock.ReadFileContentFunc(path)
}

// ReadFileContentCalls gets all the calls that were made to ReadFileContent.
// Check the length with:
//     len(mockedfileFetcher.ReadFileContentCalls())
func (mock *fileFetcherMock) ReadFileContentCalls() []struct {
	Path domain.Path
} {
	var calls []struct {
		Path domain.Path
	}
	lockfileFetcherMockReadFileContent.RLock()
	calls = mock.calls.ReadFileContent
	lockfileFetcherMockReadFileContent.RUnlock()
	return calls
}