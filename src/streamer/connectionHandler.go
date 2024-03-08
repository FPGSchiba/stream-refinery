package streamer

import (
	"bufio"
	"fmt"
	"net"
)

type StreamerService struct {
}

// TODO: Remove this
func HandleMasterConnection(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(message)
	}
}

// TODO: Remove this
func HandleRefineryConnection(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(message)
	}
}
