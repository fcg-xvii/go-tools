package bufio

import (
	"bytes"
	"io"
	"log"
)

var (
	ReadBufferSize = 1024
)

func DelimRemove(data, delim []byte) []byte {
	if bytes.HasSuffix(data, delim) {
		data = data[:len(data)-len(delim)+1]
	}
	return data
}

func NewReader(r io.Reader) *Reader {
	return &Reader{
		r: r,
	}
}

type Reader struct {
	r    io.Reader
	buf  bytes.Buffer
	seek int
}

func (s *Reader) fromBuf(data []byte) (res []byte) {
	if s.buf.Len() > 0 {
		res = append(s.buf.Bytes(), data...)
		s.buf.Reset()
		s.seek = 0
	} else {
		res = data
	}
	return
}

func (s *Reader) toBuf(data []byte) {
	s.seek = s.buf.Len()
	s.buf.Write(data)
}

func (s *Reader) scanBuf(delim []byte) (res []byte, check bool) {
	buf := s.buf.Bytes()[s.seek:]
	if index := bytes.Index(buf, delim); index >= 0 {
		rSize := s.seek + index + len(delim)
		check, res = true, make([]byte, rSize)
		s.buf.Read(res)
		s.seek = 0
	} else {
		s.seek = s.buf.Len()
	}
	return
}

func (s *Reader) scanNeeded() bool {
	return s.seek < s.buf.Len()
}

func (s *Reader) readBufferSize(delim []byte) int {
	if ReadBufferSize >= len(delim) {
		return ReadBufferSize
	}
	return len(delim)
}

func (s *Reader) ReadBytes(delim []byte) (res []byte, err error) {
	var check bool
	// scan internal buffer
	if s.scanNeeded() {
		log.Println("NEEDED")
		if res, check = s.scanBuf(delim); check {
			return
		}
	}
	// read from external buffer
	buf, count := make([]byte, s.readBufferSize(delim)), 0
	for {
		if count, err = s.r.Read(buf); count > 0 {
			s.buf.Write(buf[:count])
			for s.scanNeeded() {
				if res, check = s.scanBuf(delim); check {
					return
				}
			}
		}
		if err != nil {
			if err == io.EOF && s.buf.Len() > 0 {
				res, check = s.buf.Bytes(), true
				s.buf.Reset()
			}
			return
		}
	}
}
