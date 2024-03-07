package streamer

import (
	"bufio"
	"fmt"
	"net"
)

func HandleMasterConnection(conn net.Conn) {
	for {
		message, _ := bufio.NewReader(conn).ReadString('\n')
		fmt.Println(message)
	}
}
