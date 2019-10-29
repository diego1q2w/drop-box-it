package http

import (
	"compress/gzip"
	"io"
	"net/http"
)

func UncompressHandler(targetHandler http.HandlerFunc) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		var body io.ReadCloser
		var err error

		switch request.Header.Get("Content-Encoding") {
		case "gzip":
			body, err = gzip.NewReader(request.Body)
			if err != nil {
				http.Error(writer, "can't read gzip body", http.StatusBadRequest)
				return
			}
			break
		default:
			body = request.Body
		}

		request.Body = body
		targetHandler(writer, request)
	}
}
