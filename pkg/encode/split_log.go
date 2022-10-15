package encode

import (
	"fmt"
	"os"
)

// building block; write blobs to a sequence of files.

// Writes either a dense
type SplitLogWriter struct {
	// byte size of target file size in bytes, eg 25M
	maxfile int
	// buffer the file an write it at once.
	b []byte
	// prefix is {table}_{version}
	prefix string
	// added to the prefix {table}_{version}_{offset}
	offset int
	//
}

func (s *SplitLogWriter) Close() error {
	err := os.WriteFile(fmt.Sprintf("%s_%d", s.prefix, s.offset), s.b, os.ModePerm)
	s.b = nil
	return err
}
func (s *SplitLogWriter) Write(value []byte) error {
	// write as much as we can,then start the next file
	var err error
	remain := s.maxfile - len(s.b)
	if len(value) < remain {
		s.b = append(s.b, value...)
	} else {
		s.b = append(s.b, value[0:remain]...)
		err = os.WriteFile(fmt.Sprintf("%s_%d", s.prefix, s.offset), s.b, os.ModePerm)
		s.offset += s.maxfile
		s.b = value[remain:]
	}
	return err
}
func (s *SplitLogWriter) WriteAll(p [][]byte) error {
	// write as much as we can,then start the next file
	for _, v := range p {
		err := s.Write(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenSplitLog(p string, maxfile int) *SplitLogWriter {
	return &SplitLogWriter{
		maxfile: maxfile,
		b:       make([]byte, 0, maxfile),
		prefix:  p,
		offset:  0,
	}
}
