package common

import (
	"encoding/binary"

	log "github.com/sirupsen/logrus"
)

type ConnectionHandler struct {
	conn  *Stream
	Id    string
	msgID int
}

func InitConnectionHandler(conn *Stream, id string) *ConnectionHandler {
	handler := &ConnectionHandler{
		conn:  conn,
		Id:    id,
		msgID: 1,
	}
	return handler
}

func (h *ConnectionHandler) Close() {
	h.conn.Close()
}

func (h *ConnectionHandler) Write(message string) error {
	messageLength := len(message)
	messageLenghtBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(messageLenghtBytes, uint32(messageLength))
	err := h.conn.Write(messageLenghtBytes)
	if err != nil {
		return err
	}
	err = h.conn.Write([]byte(message))
	if err != nil {
		return err
	}
	h.msgID++
	log.Infof("action: write | result: success | message_id: %v | message: %v | message_length: %v", h.msgID, message, messageLength)
	return nil
}

func (h *ConnectionHandler) Read() (string, error) {
	messageLenghtBytes, err := h.conn.Read(4)
	if err != nil {
		return "", err
	}
	messageLenght := binary.BigEndian.Uint32(messageLenghtBytes)
	messageBytes, err := h.conn.Read(int(messageLenght))
	if err != nil {
		return "", err
	}
	log.Infof("action read | result: success  message: %v | message_lenght: %v", string(messageBytes), messageLenght)
	return string(messageBytes), nil
}
