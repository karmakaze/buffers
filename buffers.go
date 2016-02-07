package buffers

import "io"

type CloseWithError interface {
	CloseWithError(error) error
}

type PipeReader interface {
	io.Reader
	io.Closer
	CloseWithError
}

type PipeWriter interface {
	io.Writer
	io.Closer
	CloseWithError
}
