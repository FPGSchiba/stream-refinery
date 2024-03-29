package refinery

import (
	"errors"
	"fmt"
	"net"
	"streamref/src/cluster"
)

func authenticate(conn net.Conn, nodeID string) error {
	message := cluster.ConstructPacket(cluster.ConnEstablish, map[string]interface{}{"id": nodeID, "version": Version, "type": "refinery"})
	_, err := conn.Write(message)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send establish packet: %s", err.Error()))
	}
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read response from master: %s", err.Error()))
	}
	code, _, err := cluster.DeconstructPacket(response)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to parse Package: %s", err.Error()))
	}
	if code == cluster.ConnStartAuth {
		fmt.Println("Ready to start auth")
		message = cluster.ConstructPacket(cluster.AuthStart, map[string]interface{}{"cert": ""})
		_, err := conn.Write(message)
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to send establish packet: %s", err.Error()))
		}
		// TODO: Handle Auth Ack or Dec
		return nil
	} else if code == cluster.ConnClose {
		return errors.New("version of master did not match version of Node. Please update the master or the Node")
	}
	return errors.New("master did not respond with correct Code for Authentication")
}
