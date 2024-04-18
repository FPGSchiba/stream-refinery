package refinery

import (
	"crypto/rsa"
	"fmt"
	"net"
	"os"
	"streamref/src/node"
	"streamref/src/streamer"
	"streamref/src/util"
)

const (
	Version = "0.0.1"
)

type NodeRefinery struct {
	node.Node
	publicKey *rsa.PublicKey
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

func (n NodeRefinery) readPublicKey() *rsa.PublicKey {
	publicKey, err := node.LoadRsaPublicKey(n.CertificatePath)
	if err != nil {
		n.Logger.Log(fmt.Sprintf("Could not load Public Key: %s", err.Error()), util.LevelError)
		os.Exit(util.CertificateError)
	}
	return publicKey
}

func (n NodeRefinery) Start(logger util.Logger) {
	n.Node.Start(logger)
	n.publicKey = n.readPublicKey()
	// go n.startSteamService()
	n.startClusterClient()
}
