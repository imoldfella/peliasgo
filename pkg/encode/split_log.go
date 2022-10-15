package encode

import (
	"fmt"
	"os"
)

type SplitLog struct {
	b       []byte
	prefix  string
	count   int
	maxfile int
}

func (s *SplitLog) Write(p []byte) error {
	// write as much as we can,then start the next file
	var err error
	remain := s.maxfile - len(s.b)
	if len(p) < remain {
		s.b = append(s.b, p...)
	} else {
		s.b = append(s.b, p[0:remain]...)
		err = os.WriteFile(fmt.Sprintf("%s_%d", s.prefix, s.count), s.b, os.ModePerm)
		s.b = p[remain:]
	}
	return err
}
func (s *SplitLog) WriteAll(p [][]byte) error {
	// write as much as we can,then start the next file
	for _, v := range p {
		err := s.Write(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenSplitLog(p string, maxfile int) *SplitLog {
	return &SplitLog{
		prefix: p,
		count:  0,
		b:      make([]byte, 0, maxfile),
	}
}
