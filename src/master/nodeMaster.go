package master

import (
	"crypto/rsa"
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
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
}

const (
	Version  = "0.0.1"
	certSize = 4096
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

func (n NodeMaster) generateCertificate() {
	bitSize := certSize

	privateKey, err := node.LoadRsaPrivateKey(n.KeyPath)
	if err != nil {
		privateKey, err := node.GeneratePrivateKey(bitSize)
		n.publicKey = &privateKey.PublicKey
		n.privateKey = privateKey
		if err != nil {
			n.Logger.Log(fmt.Sprintf("Could not generate Private Key: %s", err.Error()), util.LevelError)
			os.Exit(util.CertificateError)
		}

		publicKeyBytes, err := node.GeneratePublicKey(&privateKey.PublicKey)
		if err != nil {
			n.Logger.Log(fmt.Sprintf("Could not generate Public Key: %s", err.Error()), util.LevelError)
			os.Exit(util.CertificateError)
		}

		privateKeyBytes := node.EncodePrivateKeyToPEM(privateKey)

		err = node.WriteKeyToFile(privateKeyBytes, n.KeyPath)
		if err != nil {
			n.Logger.Log(fmt.Sprintf("Could not save Private Key: %s", err.Error()), util.LevelError)
			os.Exit(util.CertificateError)
		}

		err = node.WriteKeyToFile(publicKeyBytes, n.CertificatePath)
		if err != nil {
			n.Logger.Log(fmt.Sprintf("Could not save Public Key: %s", err.Error()), util.LevelError)
			os.Exit(util.CertificateError)
		}
	} else {
		n.publicKey = &privateKey.PublicKey
		n.privateKey = privateKey
	}
}

func (n NodeMaster) Start(logger util.Logger) {
	// Super Starting Node
	n.Node.Start(logger)
	// Generate or fetch Certificate
	n.generateCertificate()
	// Starting Config Web Server
	go func() {
		err := n.startHTTPService()
		if err != nil {
			logger.Log(err.Error(), util.LevelError)
			os.Exit(util.HTTPServeError)
		}
	}()
	go n.startStreamService() // Starting Stream listener
	n.startClusterService()   // Master Service
}
