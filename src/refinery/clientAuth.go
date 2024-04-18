package refinery

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"streamref/src/cluster"
	"streamref/src/node"
	_ "streamref/src/node"
)

func authenticate(conn net.Conn, nodeID string, key *rsa.PublicKey) error {
	err := cluster.SendMessage(conn, cluster.ConnEstablish, map[string]interface{}{"id": nodeID, "version": Version, "type": "refinery"})
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send establish packet: %s", err.Error()))
	}
	packet, err := cluster.ReadNextMessage(conn)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read response from master: %s", err.Error()))
	}
	if packet.Code == cluster.ConnStartAuth {
		fmt.Println("Ready to start auth")
		keyBytes := node.EncodePublicKeyToPEM(key)
		err := cluster.SendMessage(conn, cluster.AuthStart, map[string]interface{}{"cert": keyBytes})
		if err != nil {
			return errors.New(fmt.Sprintf("Failed to send establish packet: %s", err.Error()))
		}
		// TODO: Handle Auth Ack or Dec
		return nil
	} else if packet.Code == cluster.ConnClose {
		return errors.New("version of master did not match version of Node. Please update the master or the Node")
	}
	return errors.New("master did not respond with correct Code for Authentication")
}
