package encode

import (
	"fmt"
	"os"
)

type SnapshotWriter struct {
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

// the header includes the start of each block so hard to write without writing those blocks first
func (s *SnapshotWriter) WriteHeader(ref []int) {

}

func (s *SnapshotWriter) Write(key, value []byte) error {
	// write as much as we can,then start the next file
	var err error
	remain := s.maxfile - len(s.b)
	if len(value) < remain {
		s.b = append(s.b, value...)
	} else {
		s.b = append(s.b, value[0:remain]...)
		err = os.WriteFile(fmt.Sprintf("%s_%d", s.prefix, s.offset), s.b, os.ModePerm)
		s.b = value[remain:]
	}
	return err
}
func (s *SnapshotWriter) WriteAll(p [][]byte) error {
	// write as much as we can,then start the next file
	for _, v := range p {
		err := s.Write(v)
		if err != nil {
			return err
		}
	}
	return nil
}

func OpenSplitLog(p string, maxfile int) *SnapshotWriter {
	return &SnapshotWriter{
		prefix: p,
		offset: 0,
		b:      make([]byte, 0, maxfile),
	}
}
