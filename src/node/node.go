package node

import (
	"fmt"
	"streamref/src/util"
)

type Node struct {
	NodeType    string
	MasterHost  string
	Logger      util.Logger
	Development bool
}

func (n Node) Log(message string, level int) {
	n.Logger.Log(fmt.Sprintf("(%s) - %s", n.NodeType, message), level)
}

func (n Node) Start(logger util.Logger, development bool) {
	n.Logger = logger
	n.Development = development
	n.Log(fmt.Sprintf("Configured Node: %s with upstream master: %s", n.NodeType, n.MasterHost), util.LevelInfo)
}
