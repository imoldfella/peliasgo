package encode

import (
	"fmt"
	"io"
	"os"
)

type SplitLog struct {
	b      []byte
	w      io.Writer
	prefix string
	count  int
}

// Close implements io.WriteCloser
func (o *SplitLog) Close() error {
	o.Flush()
	return nil
}

func (o *SplitLog) Flush() error {
	if len(o.b) == 0 {
		return nil
	}
	return os.WriteFile(fmt.Sprintf("%s_%d", o.prefix, o.count), o.b, os.ModePerm)
	o.b = o.b[:0]
	o.count++
	return nil
}

// Write implements io.Writer
func (*SplitLog) Write(p []byte) (n int, err error) {
	// write as much as we can,then start the next file
	return 0, nil
}
func (*SplitLog) WriteAll(p [][]byte) (n int, err error) {
	// write as much as we can,then start the next file
	return 0, nil
}

var _ io.WriteCloser = (*SplitLog)(nil)

func OpenSplitLog(p string, maxfile int) *SplitLog {
	return &SplitLog{
		prefix: p,
		count:  0,
		b:      make([]byte, maxfile),
	}
}
