package cluster

import (
	"encoding/json"
	"fmt"
	"net"
	"streamref/src/util"
	"strings"
)

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

func ConstructPacket(code string, payload map[string]interface{}) []byte {
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

func DeconstructPacket(message []byte) (string, map[string]interface{}, error) {
	finalValue := 0
	// Find filled bytes
	for i := range message {
		if message[i] != 0 {
			finalValue = i + 1
		}
	}
	// Only convert filled bytes to string
	messageStr := string(message[:finalValue])
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

func ReadNextMessage(conn net.Conn) (Message, error) {
	// TODO: Handle reading of multiple Packages
	// use packageSize here
}

func SendMessage(conn net.Conn, code string, payload map[string]interface{}) error {
	// TODO: Handle sending of multiple Packages
	// use packageSize here
}
