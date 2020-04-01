package prompt

import (
	"bytes"
	"foundry/cli/logger"
	"io"
	"sync"
	"time"
)

// Buffer is a thread safe wrapper for buffer
type Buffer struct {
	buf bytes.Buffer
	mut sync.Mutex
}

// NewBuffer returns a pointer to new Buffer
func NewBuffer() *Buffer {
	return &Buffer{buf: bytes.Buffer{}}
}

func (b *Buffer) Write(p []byte) (n int, err error) {
	b.mut.Lock()
	defer b.mut.Unlock()
	return b.buf.Write(p)
}

func (b *Buffer) Read(bufCh chan<- []byte, stopCh <-chan struct{}) {
	for {
		select {
		case <-stopCh:
			return
		default:
			b.mut.Lock()

			buf := make([]byte, 1024)
			n, err := b.buf.Read(buf)

			if err == nil {
				bufCh <- buf[:n]
			} else if err != io.EOF {
				logger.Fdebugln(err)
				logger.LoglnFatal(err)
			}

			b.mut.Unlock()
		}
		time.Sleep(time.Millisecond * 10)
	}
}
