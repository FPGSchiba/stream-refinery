package node

import (
	"fmt"
	"github.com/google/uuid"
	"streamref/src/util"
)

type Node struct {
	NodeType        string
	MasterHost      string
	Logger          util.Logger
	Development     bool
	NodeID          string
	CertificatePath string
	KeyPath         string
}

func NewNode(development bool) Node {
	n := Node{Development: development, NodeID: uuid.New().String()}
	return n
}

func (n Node) Log(message string, level int) {
	n.Logger.Log(fmt.Sprintf("(%s) - %s", n.NodeType, message), level)
}

func (n Node) Start(logger util.Logger) {
	n.Logger = logger
	if n.NodeType != "master" {
		n.Log(fmt.Sprintf("Configured Node: `%s` with upstream master: %s", n.NodeType, n.MasterHost), util.LevelInfo)
	} else {
		n.Log(fmt.Sprintf("Configured Node: `%s`", n.NodeType), util.LevelInfo)
	}

}
