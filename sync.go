package buffers

import (
	"io"
	"sync"
)

type syncHelper struct {
	sync.Locker
}

type syncReader struct {
	syncHelper
	reader io.Reader
}

type syncReadCloser struct {
	syncHelper
	readCloser io.ReadCloser
}

type syncPipeReader struct {
	syncHelper
	pipeReader PipeReader
}

type syncWriter struct {
	syncHelper
	writer io.Writer
}

type syncWriteCloser struct {
	syncHelper
	writeCloser io.WriteCloser
}

type syncPipeWriter struct {
	syncHelper
	pipeWriter PipeWriter
}

func SyncReader(reader io.Reader) io.Reader {
	return &syncReader{syncHelper{&sync.Mutex{}}, reader}
}
func SyncReadCloser(readCloser io.ReadCloser) io.ReadCloser {
	return &syncReadCloser{syncHelper{&sync.Mutex{}}, readCloser}
}
func SyncPipeReader(pipeReader PipeReader) PipeReader {
	return &syncPipeReader{syncHelper{&sync.Mutex{}}, pipeReader}
}

func SyncWriter(writer io.Writer) io.Writer {
	return &syncWriter{syncHelper{&sync.Mutex{}}, writer}
}
func SyncWriteCloser(writeCloser io.WriteCloser) io.WriteCloser {
	return &syncWriteCloser{syncHelper{&sync.Mutex{}}, writeCloser}
}
func SyncPipeWriter(pipeWriter PipeWriter) PipeWriter {
	return &syncPipeWriter{syncHelper{&sync.Mutex{}}, pipeWriter}
}

// Read

func (s syncReader) Read(dest []byte) (int, error) {
	return s.read(s.reader, dest)
}
func (s syncReadCloser) Read(dest []byte) (int, error) {
	return s.read(s.readCloser, dest)
}
func (s syncPipeReader) Read(dest []byte) (int, error) {
	return s.read(s.pipeReader, dest)
}

// Write

func (s syncWriter) Write(p []byte) (int, error) {
	return s.write(s.writer, p)
}
func (s syncWriteCloser) Write(p []byte) (int, error) {
	return s.write(s.writeCloser, p)
}
func (s syncPipeWriter) Write(p []byte) (int, error) {
	return s.write(s.pipeWriter, p)
}

// Close

func (s syncReadCloser) Close() error {
	return s.close(s.readCloser)
}
func (s syncWriteCloser) Close() error {
	return s.close(s.writeCloser)
}
func (s syncPipeReader) Close() error {
	return s.close(s.pipeReader)
}
func (s syncPipeWriter) Close() error {
	return s.close(s.pipeWriter)
}

// CloseWithError

func (s syncPipeReader) CloseWithError(err error) error {
	return s.closeWithError(s.pipeReader, err)
}
func (s syncPipeWriter) CloseWithError(err error) error {
	return s.closeWithError(s.pipeWriter, err)
}

// DoAtomic

func (s syncReader) DoAtomic(block func(reader io.Reader) (int, error)) (int, error) {
	s.Lock()
	defer s.Unlock()
	return block(s.reader)
}
func (s syncReadCloser) DoAtomic(block func(readCloser io.ReadCloser) (int, error)) (int, error) {
	s.Lock()
	defer s.Unlock()
	return block(s.readCloser)
}
func (s syncPipeReader) DoAtomic(block func(pipeReader PipeReader) (int, error)) (int, error) {
	s.Lock()
	defer s.Unlock()
	return block(s.pipeReader)
}

func (s syncWriter) DoAtomic(block func(writer io.Writer) (int, error)) (int, error) {
	s.Lock()
	defer s.Unlock()
	return block(s.writer)
}
func (s syncWriteCloser) DoAtomic(block func(writCloser io.WriteCloser) (int, error)) (int, error) {
	s.Lock()
	defer s.Unlock()
	return block(s.writeCloser)
}
func (s syncPipeWriter) DoAtomic(block func(pipeWriter PipeWriter) (int, error)) (int, error) {
	s.Lock()
	defer s.Unlock()
	return block(s.pipeWriter)
}

// helper

func (s syncHelper) read(r io.Reader, p []byte) (int, error) {
	s.Lock()
	defer s.Unlock()

	var n int
	for {
		c, err := r.Read(p)
		if err != nil {
			return 0, err
		}
		n += c
		if c == len(p) {
			return n, nil
		}
		p = p[c:]
	}
}

func (s syncHelper) write(w io.Writer, p []byte) (int, error) {
	s.Lock()
	defer s.Unlock()

	var n int
	for {
		c, err := w.Write(p)
		if err != nil {
			return 0, err
		}
		n += c
		if c == len(p) {
			return n, nil
		}
		p = p[c:]
	}
}

func (s syncHelper) close(c io.Closer) error {
	s.Lock()
	defer s.Unlock()
	return c.Close()
}

func (s syncHelper) closeWithError(c CloseWithError, err error) error {
	s.Lock()
	defer s.Unlock()
	return c.CloseWithError(err)
}
