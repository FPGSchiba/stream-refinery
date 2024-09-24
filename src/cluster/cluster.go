package cluster

import (
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"regexp"
	"strconv"
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
		resString = fmt.Sprintf("%s%s%d%s%s%s", code, codeDelimiter, len(payloadStr), codeDelimiter, payloadStr, endOfLine)
	} else {
		resString = fmt.Sprintf("%s%s0%s%s", code, codeDelimiter, codeDelimiter, endOfLine)
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
	packetParts := strings.Split(message, codeDelimiter)
	code := packetParts[0]
	payloadStr := strings.ReplaceAll(packetParts[2], endOfLine, "")
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

func parsePackageLength(data []byte, finalValue int) int {
	dataStr := string(data[:finalValue])
	headerRegex, err := regexp.Compile(".+;;")
	if err != nil {
		return 0
	}
	header := headerRegex.Find(data)
	packageParts := strings.Split(dataStr, codeDelimiter)
	payloadLength, err := strconv.Atoi(packageParts[1])
	if err != nil {
		return 0
	}
	return len(header) + payloadLength + len(endOfLine)
}

func ReadNextMessage(conn net.Conn) (Message, error) {
	buffer := make([]byte, packageSize)
	var content string
	var readLength = 0
	var packageLength = 0
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

		if readLength == 0 {
			packageLength = parsePackageLength(buffer, finalValue)
		}
		if packageLength == 0 {
			return Message{}, errors.New("failed to parse package length")
		}
		if packageLength == finalValue {
			content = string(buffer[:packageLength])
			break
		}
		if packageLength > (finalValue + (readLength * packageSize)) {
			content += string(buffer[:finalValue])
		}
		if packageLength < (finalValue + (readLength * packageSize)) {
			content += string(buffer[:packageLength-(readLength*packageSize)])
		}
		readLength++
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
