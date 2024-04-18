package master

import (
	"crypto/rsa"
	"errors"
	"fmt"
	"net"
	"streamref/src/cluster"
	"streamref/src/node"
)

func closeConnectionPlanned(conn net.Conn) error {
	err := cluster.SendMessage(conn, cluster.ConnClose, nil)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send close to node: %s", err.Error()))
	}
	err = conn.Close()
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to close connection: %s", err.Error()))
	}
	return nil
}

func establish(conn net.Conn) (string, string, error) {
	packet, err := cluster.ReadNextMessage(conn)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Failed to deconstruct Package: %s", err.Error()))
	}
	switch packet.Code {
	case cluster.ConnEstablish:
		if Version == packet.Data["version"].(string) {
			return packet.Data["id"].(string), packet.Data["type"].(string), nil
		} else {
			err := closeConnectionPlanned(conn)
			if err != nil {
				return "", "", err
			}
			return "", "", errors.New(fmt.Sprintf("Node with different Version tried to connect. Expected Version: %s, given Version: %s", Version, packet.Data["version"].(string)))
		}
	default:
		return "", "", errors.New(fmt.Sprintf("Did not expect code: %s here.", packet.Code))
	}
}

func authenticate(conn net.Conn) (*rsa.PublicKey, error) {
	err := cluster.SendMessage(conn, cluster.ConnStartAuth, nil)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to send Start Auth Package: %s", err.Error()))
	}
	packet, err := cluster.ReadNextMessage(conn)
	if err != nil {
		return nil, errors.New(fmt.Sprintf("Failed to read Start Auth Package: %s", err.Error()))
	}
	switch packet.Code {
	case cluster.AuthStart:
		cert := packet.Data["cert"].(string)
		pub, err := node.DecodePublicKeyFromPEM([]byte(cert))
		if err != nil {
			return nil, errors.New(fmt.Sprintf("Failed to parse Certificate: %s", err.Error()))
		}
		return pub, nil
	default:
		return nil, errors.New(fmt.Sprintf("Code not expected from Node: %s", packet.Code))
	}
}
