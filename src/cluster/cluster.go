package cluster

import (
	"encoding/json"
	"fmt"
	"net"
	"regexp"
	"streamref/src/util"
	"strings"
)

type ClusterService struct{}

type Message struct {
	Code string
	Data map[string]interface{}
}

// Protocol internal
const (
	endOfLine     = "<EOF>"
	codeDelimiter = ";;"
	packageSize   = 1024
)

// Command Codes
const (
	// Conn
	ConnStartAuth = "conn:startAuth"

	// Auth
	AuthDec = "auth:dec"
)

// Request Codes
const (
	// Conn
	ConnEstablish = "conn:establish"
	AuthStart     = "auth:start"
)

// Shared Codes
const (
	// Conn
	ConnClose = "conn:close"
	ConnAlive = "conn:alive"

	// Auth
	AuthAck = "auth:ack"
)

func constructPacket(code string, payload map[string]interface{}) []byte {
	var resString string
	payloadLength := len(payload)

	if payloadLength > 0 {
		payloadStr, err := json.Marshal(payload)
		if err != nil {
			logger := util.Logger{LogType: util.LogTypeConsole}
			logger.Log(fmt.Sprintf("Cluster Protocol Error with handling JSON: %s", err.Error()), util.LevelError)
		}
		resString = fmt.Sprintf("%s%s%s%s", code, codeDelimiter, payloadStr, endOfLine)
	} else {
		resString = fmt.Sprintf("%s%s%s", code, codeDelimiter, endOfLine)
	}

	return []byte(resString)
}

func jsonToMap(jsonStr string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	err := json.Unmarshal([]byte(jsonStr), &result)
	if err != nil {
		return nil, err
	}
	return result, nil
}

func deconstructPacket(message string) (string, map[string]interface{}, error) {
	messageStr := message
	packetParts := strings.Split(messageStr, codeDelimiter)
	code := packetParts[0]
	payloadStr := strings.ReplaceAll(packetParts[1], endOfLine, "")
	if payloadStr != "" {
		payload, err := jsonToMap(payloadStr)
		if err != nil {
			return "", nil, err
		}
		return code, payload, nil
	}
	return code, nil, nil
}

func allZero(s []byte) bool {
	for _, v := range s {
		if v != 0 {
			return false
		}
	}
	return true
}

func hasEndOfLine(data []byte) bool {
	if allZero(data) {
		return false
	}
	if strings.Contains(string(data), endOfLine) {
		return true
	}
	return false
}

func ReadNextMessage(conn net.Conn) (Message, error) {
	buffer := make([]byte, packageSize)
	var content string
	for !hasEndOfLine(buffer) {
		_, err := conn.Read(buffer)

		// Find filled bytes
		finalValue := 0
		for i := range buffer {
			if buffer[i] != 0 {
				finalValue = i + 1
			}
		}

		if err != nil {
			return Message{}, err
		}

		currentContent := string(buffer[:finalValue])
		endOfRegex, err := regexp.Compile("<EOF>.*")
		if err != nil {
			return Message{}, err
		}

		// Replacing not needed Data after end of Packet
		if hasEndOfLine(buffer) && !strings.HasSuffix(currentContent, endOfLine) {
			currentContent = endOfRegex.ReplaceAllString(currentContent, "")
		}
		content += currentContent
	}

	code, payload, err := deconstructPacket(content)
	if err != nil {
		return Message{}, err
	}
	return Message{Code: code, Data: payload}, nil
}

func SendMessage(conn net.Conn, code string, payload map[string]interface{}) error {
	message := constructPacket(code, payload)
	_, err := conn.Write(message)
	if err != nil {
		return err
	}
	return nil
}
