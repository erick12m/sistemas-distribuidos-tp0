package common

import (
	"net"

	log "github.com/sirupsen/logrus"
)

type Stream struct {
	conn net.Conn
	Id   string
}

func initStream(conn net.Conn, id string) *Stream {
	stream := &Stream{
		conn: conn,
		Id:   id,
	}
	log.Infof("Stream created with id %v", id)
	return stream
}

func (s *Stream) Close() {
	s.conn.Close()
	log.Infof("Stream with id %v closed", s.Id)
}

func (s *Stream) Write(data []byte) error {
	totalWrite := 0

	for totalWrite < len(data) {
		n, err := s.conn.Write(data[totalWrite:])
		if err != nil {
			log.Errorf("Error writing to stream with id %v: %v", s.Id, err)
			return err
		}
		totalWrite += n
	}
	log.Infof("Data written to stream with id %v", s.Id)
	return nil
}

func (s *Stream) Read(size int) ([]byte, error) {
	totalRead := 0

	message := make([]byte, size)
	for totalRead < size {
		n, err := s.conn.Read(message[totalRead:])
		if err != nil {
			log.Errorf("Error reading from stream with id %v: %v", s.Id, err)
			return nil, err
		}
		totalRead += n
	}
	log.Infof("Data read from stream with id %v", s.Id)
	return message, nil
}
