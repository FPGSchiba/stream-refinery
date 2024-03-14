package master

import (
	"fmt"
	"net"
	"slices"
	"streamref/src/cluster"
	"streamref/src/util"
)

type NodeConnection struct {
	conn     net.Conn
	nodeType string
	nodeID   string
}

type ClusterServiceMaster struct {
	cluster.ClusterService
	clients []NodeConnection
	node    NodeMaster
}

func (cs ClusterServiceMaster) Start(node NodeMaster) {
	cs.node = node
	cs.node.Log(fmt.Sprintf("Cluster Service listening on: 0.0.0.0:%d", util.DefaultClusterPort), util.LevelDebug)
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", util.DefaultClusterPort)) // TODO: Error Handling
	for {
		conn, _ := ln.Accept() // TODO: Error handling
		connection := NodeConnection{conn: conn}
		cs.clients = append(cs.clients, connection)
		go cs.HandleConnection(conn)
	}
}

func (cs ClusterServiceMaster) HandleConnection(conn net.Conn) {
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			cs.node.Log(err.Error(), util.LevelError)
		}
	}(conn)

	nodeID, nodeType, err := establish(conn)
	if err != nil {
		cs.node.Log(err.Error(), util.LevelError)
		return
	}

	idx := slices.IndexFunc(cs.clients, func(c NodeConnection) bool { return c.conn.RemoteAddr() == conn.RemoteAddr() })
	cs.clients[idx].nodeType = nodeType
	cs.clients[idx].nodeID = nodeID

	err = authenticate(conn)
	if err != nil {
		cs.node.Log(err.Error(), util.LevelError)
		return
	}

	for {
		// Read incoming data
		buf := make([]byte, 1024)
		_, err := conn.Read(buf)
		if err != nil {
			fmt.Println(err)
			break
		}

		// Print the incoming data
		fmt.Printf("Received: %s", buf)
	}
}
