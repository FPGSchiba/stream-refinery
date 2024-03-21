package config

import (
	"errors"
	"fmt"
	"streamref/src/master"
	"streamref/src/node"
	"streamref/src/refinery"
)

const (
	Master   = "master"
	Sub      = "submaster"
	Refinery = "refinery"
	Receiver = "receiver"
)

type ReturnNodeType struct {
	NodeType         string
	ResultNodeMaster struct {
		NodeRes master.NodeMaster
	}
	ResultNodeRefinery struct {
		NodeRes refinery.NodeRefinery
	}
	ResultNodeSubMaster struct {
		NodeRes NodeSubMaster
	}
	ResultNodeReceiver struct {
		NodeRes NodeReceiver
	}
}

func Construct(node node.Node) (ReturnNodeType, error) {
	switch node.NodeType {
	case Master:
		var masterNode master.NodeMaster
		masterNode.NodeType = node.NodeType
		masterNode.NodeID = node.NodeID
		masterNode.Development = node.Development
		masterNode.KeyPath = node.KeyPath
		masterNode.CertificatePath = node.CertificatePath
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeMaster.NodeRes = masterNode
		return result, nil
	case Sub:
		var subNode NodeSubMaster
		subNode.NodeType = node.NodeType
		subNode.MasterHost = node.MasterHost
		subNode.NodeID = node.NodeID
		subNode.Development = node.Development
		subNode.KeyPath = node.KeyPath
		subNode.CertificatePath = node.CertificatePath
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeSubMaster.NodeRes = subNode
		return result, nil
	case Refinery:
		var refineryNode refinery.NodeRefinery
		refineryNode.NodeType = node.NodeType
		refineryNode.MasterHost = node.MasterHost
		refineryNode.NodeID = node.NodeID
		refineryNode.Development = node.Development
		refineryNode.KeyPath = node.KeyPath
		refineryNode.CertificatePath = node.CertificatePath
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeRefinery.NodeRes = refineryNode
		return result, nil
	case Receiver:
		var nodeReceiver NodeReceiver
		nodeReceiver.NodeType = node.NodeType
		nodeReceiver.MasterHost = node.MasterHost
		nodeReceiver.NodeID = node.NodeID
		nodeReceiver.Development = node.Development
		nodeReceiver.KeyPath = node.KeyPath
		nodeReceiver.CertificatePath = node.CertificatePath
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeReceiver.NodeRes = nodeReceiver
		return result, nil
	}
	return ReturnNodeType{}, errors.New(fmt.Sprintf("Node Type: %s was not found.", node.NodeType))
}
