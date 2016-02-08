package buffers

import (
	"fmt"
	"io"
	"sync"
)

type Pipe struct {
	buffer   *Circular
	rClosed  error // nil while reading
	wClosed  error // nil while writing
	cond     *sync.Cond
}

type pipeReader struct {
	pipe *Pipe
}

type pipeWriter struct {
	pipe *Pipe
}

func NewPipe(size int) (PipeReader, PipeWriter) {
	pipe := &Pipe{NewCircular(size), nil, nil, sync.NewCond(&sync.Mutex{})}
	return pipeReader{pipe}, pipeWriter{pipe}
}

func (r pipeReader) Read(dest []byte) (int, error) {
	r.pipe.cond.L.Lock()
	defer r.pipe.cond.L.Unlock()

	for {
		if r.pipe.rClosed != nil {
			return 0, io.ErrClosedPipe
		}

		n, err := r.pipe.buffer.Read(dest)
		if n != 0 {
			r.pipe.cond.Broadcast()
			return n, err
		}
		if err != nil {
			return 0, err
		}
		if r.pipe.wClosed != nil && r.pipe.buffer.Buffered() == 0 {
			return 0, r.pipe.wClosed
		}
		if len(dest) == 0 {
			return 0, nil
		}

		r.pipe.cond.Wait()
	}
}

func (r pipeReader) Close() error {
	return r.CloseWithError(io.ErrClosedPipe)
}

func (r pipeReader) CloseWithError(err error) error {
	if err == nil {
		return fmt.Errorf("Can't BufferedPipeReader.CloseWithError(nil)")
	}

	r.pipe.cond.L.Lock()
	defer r.pipe.cond.L.Unlock()

	if r.pipe.rClosed == nil {
		r.pipe.rClosed = err
		r.pipe.cond.Broadcast()
	} else if r.pipe.rClosed != err && (err == nil || err.Error() != r.pipe.rClosed.Error()) {
		return fmt.Errorf("Ignoring BufferedPipeReader.CloseWithError(%v); already closed with (%v)", err, r.pipe.rClosed)
	}
	return nil
}

func (w pipeWriter) Write(src []byte) (int, error) {
	w.pipe.cond.L.Lock()
	defer w.pipe.cond.L.Unlock()

	for {
		if w.pipe.wClosed != nil {
			return 0, io.ErrClosedPipe
		}
		if w.pipe.rClosed != nil {
			return 0, w.pipe.rClosed
		}

		n, err := w.pipe.buffer.Write(src)
		if n != 0 {
			w.pipe.cond.Broadcast()
			return n, err
		}
		if err != nil {
			return 0, err
		}
		if len(src) == 0 {
			return 0, nil
		}

		w.pipe.cond.Wait()
	}
}

func (w pipeWriter) Close() error {
	return w.CloseWithError(io.EOF)
}

func (w pipeWriter) CloseWithError(err error) error {
	if err == nil {
		return fmt.Errorf("Can't BufferedPipeWriter.CloseWithError(nil)")
	}

	w.pipe.cond.L.Lock()
	defer w.pipe.cond.L.Unlock()

	if w.pipe.wClosed == nil {
		w.pipe.wClosed = err
		w.pipe.cond.Broadcast()
	} else if w.pipe.wClosed != err && (err == nil || err.Error() != w.pipe.wClosed.Error()) {
		return fmt.Errorf("Ignoring BufferedPipeWriter.CloseWithError(%v); already closed with (%v)", err, w.pipe.wClosed)
	}
	return nil
}
