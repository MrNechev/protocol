package protocol

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"io"
)

const messageHeaderSize = 4

type Request struct {
	Command string `json:"command,omitempty"`
	Domain  string `json:"domain"`
}

type Response struct {
	Domain       string   `json:"domain"`
	Registrar    string   `json:"registrar"`
	Registration string   `json:"registration date"`
	Expiration   string   `json:"expiration date"`
	NameServers  []string `json:"nameservers"`
	ErrorMessage string   `json:"error,omitempty"`
}

func EncodeMessage(message interface{}, writer io.Writer) error {
	messageBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	messageHeader := make([]byte, messageHeaderSize)
	binary.BigEndian.PutUint32(messageHeader, uint32(len(messageBytes)))
	if _, err := writer.Write(messageHeader); err != nil {
		return err
	}
	if _, err := writer.Write(messageBytes); err != nil {
		return err
	}
	return nil
}

func DecodeMessage(reader io.Reader) (interface{}, error) {
	messageHeader := make([]byte, messageHeaderSize)
	if _, err := io.ReadFull(reader, messageHeader); err != nil {
		return nil, err
	}
	messageSize := binary.BigEndian.Uint32(messageHeader)
	messageBytes := make([]byte, messageSize)
	if _, err := io.ReadFull(reader, messageBytes); err != nil {
		return nil, err
	}
	var message interface{}
	if err := json.Unmarshal(messageBytes, &message); err != nil {
		return nil, err
	}
	return message, nil
}

func ParseRequest(message interface{}) (*Request, error) {
	request, ok := message.(map[string]interface{})
	if !ok {
		return nil, errors.New("invalid request message format")
	}
	domain, ok := request["domain"].(string)
	if !ok {
		return nil, errors.New("missing or invalid domain in request message")
	}
	command, _ := request["command"].(string)
	return &Request{Command: command, Domain: domain}, nil
}

func FormatResponse(response *Response) (interface{}, error) {
	if response.ErrorMessage != "" {
		return map[string]string{"error": response.ErrorMessage}, nil
	}
	return response, nil
}
