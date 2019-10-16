// Code generated by moq; DO NOT EDIT.
// github.com/matryer/moq

package app

import (
	"github.com/diego1q2w/drop-box-it/pkg/box/domain"
	"sync"
)

var (
	lockfileBoxerMockDeleteFile sync.RWMutex
	lockfileBoxerMockWriteFile  sync.RWMutex
)

// Ensure, that fileBoxerMock does implement fileBoxer.
// If this is not the case, regenerate this file with moq.
var _ fileBoxer = &fileBoxerMock{}

// fileBoxerMock is a mock implementation of fileBoxer.
//
//     func TestSomethingThatUsesfileBoxer(t *testing.T) {
//
//         // make and configure a mocked fileBoxer
//         mockedfileBoxer := &fileBoxerMock{
//             DeleteFileFunc: func(path domain.Path) error {
// 	               panic("mock out the DeleteFile method")
//             },
//             WriteFileFunc: func(file domain.File) error {
// 	               panic("mock out the WriteFile method")
//             },
//         }
//
//         // use mockedfileBoxer in code that requires fileBoxer
//         // and then make assertions.
//
//     }
type fileBoxerMock struct {
	// DeleteFileFunc mocks the DeleteFile method.
	DeleteFileFunc func(path domain.Path) error

	// WriteFileFunc mocks the WriteFile method.
	WriteFileFunc func(file domain.File) error

	// calls tracks calls to the methods.
	calls struct {
		// DeleteFile holds details about calls to the DeleteFile method.
		DeleteFile []struct {
			// Path is the path argument value.
			Path domain.Path
		}
		// WriteFile holds details about calls to the WriteFile method.
		WriteFile []struct {
			// File is the file argument value.
			File domain.File
		}
	}
}

// DeleteFile calls DeleteFileFunc.
func (mock *fileBoxerMock) DeleteFile(path domain.Path) error {
	if mock.DeleteFileFunc == nil {
		panic("fileBoxerMock.DeleteFileFunc: method is nil but fileBoxer.DeleteFile was just called")
	}
	callInfo := struct {
		Path domain.Path
	}{
		Path: path,
	}
	lockfileBoxerMockDeleteFile.Lock()
	mock.calls.DeleteFile = append(mock.calls.DeleteFile, callInfo)
	lockfileBoxerMockDeleteFile.Unlock()
	return mock.DeleteFileFunc(path)
}

// DeleteFileCalls gets all the calls that were made to DeleteFile.
// Check the length with:
//     len(mockedfileBoxer.DeleteFileCalls())
func (mock *fileBoxerMock) DeleteFileCalls() []struct {
	Path domain.Path
} {
	var calls []struct {
		Path domain.Path
	}
	lockfileBoxerMockDeleteFile.RLock()
	calls = mock.calls.DeleteFile
	lockfileBoxerMockDeleteFile.RUnlock()
	return calls
}

// WriteFile calls WriteFileFunc.
func (mock *fileBoxerMock) WriteFile(file domain.File) error {
	if mock.WriteFileFunc == nil {
		panic("fileBoxerMock.WriteFileFunc: method is nil but fileBoxer.WriteFile was just called")
	}
	callInfo := struct {
		File domain.File
	}{
		File: file,
	}
	lockfileBoxerMockWriteFile.Lock()
	mock.calls.WriteFile = append(mock.calls.WriteFile, callInfo)
	lockfileBoxerMockWriteFile.Unlock()
	return mock.WriteFileFunc(file)
}

// WriteFileCalls gets all the calls that were made to WriteFile.
// Check the length with:
//     len(mockedfileBoxer.WriteFileCalls())
func (mock *fileBoxerMock) WriteFileCalls() []struct {
	File domain.File
} {
	var calls []struct {
		File domain.File
	}
	lockfileBoxerMockWriteFile.RLock()
	calls = mock.calls.WriteFile
	lockfileBoxerMockWriteFile.RUnlock()
	return calls
}
