package master

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"streamref/src/config/handler"
	"streamref/src/node"
	"streamref/src/streamer"
	"streamref/src/util"
)

type NodeMaster struct {
	node.Node
}

const (
	Version = "0.0.1"
)

func (n NodeMaster) startHTTPService() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.GetRoot)

	n.Log(fmt.Sprintf("HTTP Service listening on: 0.0.0.0:%d", util.DefaultConfigPort), util.LevelDebug)
	err := http.ListenAndServe(fmt.Sprintf(":%d", util.DefaultConfigPort), mux)
	return err
}

func (n NodeMaster) startStreamService() {
	n.Log(fmt.Sprintf("Stream Service listening on: 0.0.0.0:%d", util.DefaultStreamPort), util.LevelDebug)
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", util.DefaultStreamPort))
	conn, _ := ln.Accept()
	go streamer.HandleMasterConnection(conn)
}

func (n NodeMaster) startClusterService() {
	var clusterService = ClusterServiceMaster{}
	clusterService.Start(n)
}

func (n NodeMaster) Start(logger util.Logger) {
	n.Node.Start(logger)
	go func() {
		err := n.startHTTPService()
		if err != nil {
			logger.Log(err.Error(), util.LevelError)
			os.Exit(util.HTTPServeError)
		}
	}()
	go n.startStreamService()
	n.startClusterService() // Master Service
}
