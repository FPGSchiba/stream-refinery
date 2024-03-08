package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"streamref/src/config"
	"streamref/src/node"
	"streamref/src/util"
)

func handleArgs(args []string) (node.Node, error) {
	var node node.Node
	if len(args) > 1 {
		switch len(args) {
		case 2:
			node.NodeType = config.Master
		case 3:
			if args[1] == config.Sub || args[1] == config.Refinery || args[1] == config.Receiver {
				node.NodeType = args[1]
			} else {
				message := fmt.Sprintf("NodeType: %s is not known.", args[1])
				return node, errors.New(message)
			}
			regexIP, _ := regexp.Compile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
			regexHostName, _ := regexp.Compile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$")

			if regexIP.MatchString(args[2]) || regexHostName.MatchString(args[2]) {
				node.MasterHost = args[2]
			} else {
				message := fmt.Sprintf("Hostname: %v is not a valid hostname.", args[1])
				return node, errors.New(message)
			}
		}
	} else {
		return node, errors.New("at least one argument expected")
	}
	return node, nil
}

func main() {
	argsWithProg := os.Args
	logger := util.Logger{
		LogLevel: util.LevelDebug,
		FilePath: "/Users/schiba/Projects/stream-refinery/log.txt",
		LogType:  util.LogTypeConsole,
	}
	fmt.Println(fmt.Sprintf("Configured Logger: %v", logger))
	resultNode, err := handleArgs(argsWithProg)
	if err != nil {
		logger.Log(err.Error(), util.LevelError)
		os.Exit(util.ArgumentErrorCode)
	}
	result, err := config.Construct(resultNode)
	if err != nil {
		logger.Log(err.Error(), util.LevelError)
		os.Exit(util.NodeTypeError)
	}
	switch result.NodeType {
	case config.Master:
		masterNode := result.ResultNodeMaster.NodeRes
		masterNode.Start(logger)
	case config.Refinery:
		refineryNode := result.ResultNodeRefinery.NodeRes
		refineryNode.Start(logger)
		/*
			case config.Sub:
				subNode := result.ResultNodeSubMaster.NodeRes
				// TODO: Start goroutines
			case config.Receiver:
				receiverNode := result.ResultNodeReceiver.NodeRes
				// TODO: Start goroutines

		*/
	}
}
