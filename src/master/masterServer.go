package master

import (
	"fmt"
	"net"
	"streamref/src/cluster"
	"streamref/src/util"
)

type ClusterServiceMaster struct {
	cluster.ClusterService
	clients []net.Conn
	node    NodeMaster
}

func (cs ClusterServiceMaster) Start(node NodeMaster) {
	cs.node = node
	cs.node.Log(fmt.Sprintf("Cluster Service listening on: 0.0.0.0:%d", util.DefaultClusterPort), util.LevelDebug)
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", util.DefaultClusterPort)) // TODO: Error Handling
	for {
		conn, _ := ln.Accept() // TODO: Error handling
		cs.clients = append(cs.clients, conn)
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
