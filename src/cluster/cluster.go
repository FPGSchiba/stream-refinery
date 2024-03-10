package cluster

import "fmt"

type ClusterService struct {
}

const (
	endOfLine     = "<EOF>"
	codeDelimiter = ";;"
)

func ConstructPacket(code string, payload map[string]interface{}) []byte {
	var payloadStr string
	var resString string
	payloadLength := len(payload)

	if payloadLength > 0 {
		var count = 0
		for key, value := range payload {
			if payloadLength-1 == count {
				payloadStr += fmt.Sprintf("\"%s\":\"%s\"", key, value)
			} else {
				payloadStr += fmt.Sprintf("\"%s\":\"%s\",", key, value)
			}
			fmt.Println(key, value)
			count++
		}
		resString = fmt.Sprintf("%s%s%s%s", code, codeDelimiter, payloadStr, endOfLine)
	} else {
		resString = fmt.Sprintf("%s%s%s", code, codeDelimiter, endOfLine)
	}

	return []byte(resString)
}
