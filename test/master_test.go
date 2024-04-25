package test

import (
	"streamref/src/config"
	"streamref/src/master"
	"streamref/src/node"
	"streamref/src/util"
	"testing"
)

//TODO: Create Master with config

func TestMasterStart(t *testing.T) {
	node := node.NewNode(false)
	logger := util.Logger{LogType: util.LogTypeConsole, LogFile: nil, LogLevel: util.LevelDebug}
	node.NodeType = config.Master
	node.CertificatePath = "./id_rsa_test"
	node.KeyPath = "id_rsa_test.pem"
	master := master.NodeMaster{Node: node}
	master.Start(logger)
}
