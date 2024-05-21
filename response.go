package goboxer

import (
	"io"
	"net/http"
)

type Response struct {
	Request      *Request
	ContentType  string
	Headers      http.Header
	Body         io.ReadCloser
	ResponseCode int
	RTTInMillis  int64
}
