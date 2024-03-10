package refinery

import (
	"fmt"
	"net"
	"streamref/src/node"
	"streamref/src/streamer"
	"streamref/src/util"
)

type NodeRefinery struct {
	node.Node
}

func (n NodeRefinery) startClusterClient() {
	var clusterService = ClusterServiceRefinery{}
	clusterService.Start(n)
}

func (n NodeRefinery) startSteamService() {
	n.Log(fmt.Sprintf("Stream Service listening on: 0.0.0.0:%d", util.DefaultStreamPort), util.LevelDebug)
	var port int
	if port = util.DefaultConfigPort; n.Development {
		port = util.DefaultConfigPort + 2
	}
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", port))
	conn, _ := ln.Accept()
	go streamer.HandleRefineryConnection(conn)
}

func (n NodeRefinery) Start(logger util.Logger, development bool) {
	n.Node.Start(logger, development)
	// go n.startSteamService()
	n.startClusterClient()
}
