buffers
=======
A buffered pipe, sync reader/writer wrappers, and bounded circular/ring byte buffer for Go.

```
func NewPipe(size int) (PipeReader, PipeWriter)
```

   * makes a fixed-buffered pipe and returns 'io'ers for each end

```
func SyncReader(io.Reader) io.Reader
func SyncWriter(io.Writer) io.Writer
func SyncReadCloser(io.ReadCloser) io.ReadCloser
func SyncWriteCloser(io.WriteCloser) io.WriteCloser
func SyncPipeReader(PipeReader) PipeReader
func SyncPipeWriter(PipeWriter) PipeWriter
```
   * wraps underlying 'reader', 'writer' or 'pipe' providing synchronized access
   * also provides `DoAtomic` as a means to perform multiple accesses within a single unit
   * `PipeReader` and `PipeWriter` interfaces match same named structs in `io` package

```
func NewCircular(size int) *Circular
```

   * creates a fixed-size circular byte-buffer
   * is a Reader
   * is a Writer - `Write` can return `io.ErrShortWrite` with the count of partial bytes written
   * `Buffered()` returns count of bytes in buffer
   * `Reset()` clears contents
   * does not have synchronization


copyright (c) 2016, Keith Kim

license: MIT
