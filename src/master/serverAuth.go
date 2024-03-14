package master

import (
	"errors"
	"fmt"
	"net"
	"streamref/src/cluster"
)

func closeConnectionPlanned(conn net.Conn) error {
	message := cluster.ConstructPacket(cluster.ConnClose, nil)
	_, err := conn.Write(message)
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
	buf := make([]byte, 1024)
	_, err := conn.Read(buf)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Failed to read from node: %s", err.Error()))
	}
	code, payload, err := cluster.DeconstructPacket(buf)
	if err != nil {
		return "", "", errors.New(fmt.Sprintf("Failed to deconstruct Package: %s", err.Error()))
	}
	switch code {
	case cluster.ConnEstablish:
		if Version == payload["version"].(string) {
			return payload["id"].(string), payload["type"].(string), nil
		} else {
			err := closeConnectionPlanned(conn)
			if err != nil {
				return "", "", err
			}
			return "", "", errors.New(fmt.Sprintf("Node with different Version tried to connect. Expected Version: %s, given Version: %s", Version, payload["version"].(string)))
		}
	default:
		return "", "", errors.New(fmt.Sprintf("Did not expect code: %s here.", code))
	}
}

func authenticate(conn net.Conn) error {
	message := cluster.ConstructPacket(cluster.ConnStartAuth, nil)
	_, err := conn.Write(message)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to send Start Auth Package: %s", err.Error()))
	}
	buf := make([]byte, 1024)
	_, err = conn.Read(buf)
	if err != nil {
		return errors.New(fmt.Sprintf("Failed to read Start Auth Package: %s", err.Error()))
	}
	code, payload, err := cluster.DeconstructPacket(buf)
	switch code {
	case cluster.AuthStart:
		cert := payload["cert"]
		fmt.Println(fmt.Sprintf("Certificate: %s", cert))
		return nil
	default:
		return errors.New(fmt.Sprintf("Code not expected from Node: %s", code))
	}
}
