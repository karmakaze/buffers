package buffers

import "io"

type Circular struct {
	data    []byte
	readAt  int
	writeAt int
}

func NewCircular(size int) *Circular {
	// plus 1 to allow write of 'size' with 'writeAt' < 'readAt' after wrap-around
	return &Circular{make([]byte, size + 1), 0, 0}
}

func (b *Circular) Buffered() int {
	if b.readAt <= b.writeAt {
		return b.writeAt - b.readAt
	}
	return len(b.data) - b.readAt + b.writeAt
}

func (b *Circular) Reset() {
	b.readAt = 0
	b.writeAt = 0
}

func (b *Circular) Read(dest []byte) (n int, err error) {
	if b.writeAt < b.readAt {
		n = b.read(dest, len(b.data) - b.readAt)
	}
	if b.readAt < b.writeAt {
		n += b.read(dest[n:], b.writeAt - b.readAt)
	}
	return n, nil
}

func (b *Circular) Write(src []byte) (n int, err error) {
	if b.readAt <= b.writeAt {
		end := len(b.data)
		if b.readAt == 0 {
			end -= 1
		}
		n = b.write(src, end - b.writeAt)
	}
	if b.writeAt < b.readAt - 1 {
		n += b.write(src[n:], b.readAt - 1 - b.writeAt)
	}
	if n < len(src) {
		err = io.ErrShortWrite
	}
	return n, err
}

func (b *Circular) read(dest []byte, c int) int {
	c = copy(dest, b.data[b.readAt:b.readAt + c])
	b.readAt += c
	if b.readAt == len(b.data) {
		b.readAt = 0
	}
	return c
}

func (b *Circular) write(src []byte, c int) int {
	c = copy(b.data[b.writeAt:b.writeAt + c], src)
	b.writeAt += c
	if b.writeAt == len(b.data) {
		b.writeAt = 0
	}
	return c
}
