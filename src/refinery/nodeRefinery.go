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
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", util.DefaultStreamPort))
	conn, _ := ln.Accept()
	go streamer.HandleRefineryConnection(conn)
}

func (n NodeRefinery) Start(logger util.Logger) {
	n.Node.Start(logger)
	// go n.startSteamService()
	n.startClusterClient()
}
