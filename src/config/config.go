package config

import (
	"errors"
	"fmt"
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
		NodeRes NodeMaster
	}
	ResultNodeRefinery struct {
		NodeRes NodeRefinery
	}
	ResultNodeSubMaster struct {
		NodeRes NodeSubMaster
	}
	ResultNodeReceiver struct {
		NodeRes NodeReceiver
	}
}

func Construct(node Node) (ReturnNodeType, error) {
	switch node.NodeType {
	case Master:
		var masterNode NodeMaster
		masterNode.NodeType = node.NodeType
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeMaster.NodeRes = masterNode
		return result, nil
	case Sub:
		var masterNode NodeSubMaster
		masterNode.NodeType = node.NodeType
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeSubMaster.NodeRes = masterNode
		return result, nil
	case Refinery:
		var masterNode NodeRefinery
		masterNode.NodeType = node.NodeType
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeRefinery.NodeRes = masterNode
		return result, nil
	case Receiver:
		var masterNode NodeReceiver
		masterNode.NodeType = node.NodeType
		var result ReturnNodeType
		result.NodeType = node.NodeType
		result.ResultNodeReceiver.NodeRes = masterNode
		return result, nil
	}
	return ReturnNodeType{}, errors.New(fmt.Sprintf("Node Type: %s was not found.", node.NodeType))
}
