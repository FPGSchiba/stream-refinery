package refinery

import (
	"errors"
	"net"
	"streamref/src/cluster"
)

func Authenticate(conn net.Conn, nodeID string) error {
	message := cluster.ConstructPacket(cluster.ConnEstablish, map[string]interface{}{"id": nodeID, "version": Version})
	_, err := conn.Write(message)
	if err != nil {
		return err
	}
	response := make([]byte, 1024)
	_, err = conn.Read(response)
	if err != nil {
		return err
	}
	code, _, err := cluster.DeconstructPacket(response)
	if err != nil {
		return err
	}
	if code == cluster.ConnStartAuth {
		// TODO: Start Authentication
	} else if code == cluster.ConnClose {
		return errors.New("version of master did not match version of Node. Please update the master or the Node")
	}
	return errors.New("master did not respond with correct Code for Authentication")
}
