package main

import (
	"fmt"
	"os"
	"regexp"
	"streamref/src/config"
	"streamref/src/util"
)

func main() {
	argsWithProg := os.Args
	logger := util.Logger{
		LogLevel: util.LevelDebug,
		FilePath: "./test.txt",
		LogType:  util.LogTypeConsole,
	}
	var node config.Node
	if len(argsWithProg) > 1 {
		switch len(argsWithProg) {
		case 2:
			node.NodeType = config.Master
		case 3:
			if argsWithProg[1] == config.Sub || argsWithProg[1] == config.Refinery || argsWithProg[1] == config.Receiver {
				node.NodeType = argsWithProg[1]
			} else {

				message := fmt.Sprintf("NodeType: %v is not known.", argsWithProg[1])
				logger.Log(message, util.LevelError)
				os.Exit(util.ArgumentErrorCode)
			}
			regexIP, _ := regexp.Compile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
			regexHostName, _ := regexp.Compile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$")

			if regexIP.MatchString(argsWithProg[2]) || regexHostName.MatchString(argsWithProg[2]) {
				node.MasterHost = argsWithProg[2]
			} else {
				message := fmt.Sprintf("Hostname: %v is not a valid hostname.", argsWithProg[1])
				logger.Log(message, util.LevelError)
				os.Exit(util.ArgumentErrorCode)
			}
		}
	} else {
		os.Exit(util.ArgumentErrorCode)
	}
	fmt.Printf("%+v\n", node)
}
