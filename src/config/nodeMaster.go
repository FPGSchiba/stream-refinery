package config

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"streamref/src/cluster"
	"streamref/src/config/handler"
	"streamref/src/streamer"
	"streamref/src/util"
)

type NodeMaster struct {
	Node
	logger util.Logger
}

func (n NodeMaster) log(message string, level int) {
	n.logger.Log(fmt.Sprintf("(Master) - %s", message), level)
}

func (n NodeMaster) startHTTPService() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler.GetRoot)

	n.log(fmt.Sprintf("HTTP Service listening on: 0.0.0.0:%d", util.DefaultConfigPort), util.LevelDebug)
	err := http.ListenAndServe(fmt.Sprintf(":%d", util.DefaultConfigPort), mux)
	return err
}

func (n NodeMaster) startStreamService() {
	n.log(fmt.Sprintf("Stream Service listening on: 0.0.0.0:%d", util.DefaultStreamPort), util.LevelDebug)
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", util.DefaultStreamPort))
	conn, _ := ln.Accept()
	go streamer.HandleMasterConnection(conn)
}

func (n NodeMaster) startClusterService() {
	n.log(fmt.Sprintf("Cluster Service listening on: 0.0.0.0:%d", util.DefaultClusterPort), util.LevelDebug)
	ln, _ := net.Listen("tcp", fmt.Sprintf(":%d", util.DefaultClusterPort))
	conn, _ := ln.Accept()
	go cluster.HandleMasterConnection(conn)
}

func (n NodeMaster) Start(logger util.Logger) {
	n.logger = logger
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
