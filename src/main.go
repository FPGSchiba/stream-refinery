package main

import (
	"errors"
	"fmt"
	"github.com/akamensky/argparse"
	"os"
	"path/filepath"
	"regexp"
	"streamref/src/config"
	"streamref/src/node"
	"streamref/src/util"
)

func main() {
	parser := argparse.NewParser("stream-refinery", "To start a node of the Stream-Refinery Cluster following arguments can be specified :D Please support this project on: https://github.com/FPGSchiba/stream-refinery.          Thanks for using Stream-Refinery")
	nodeType := parser.String("t", "type", &argparse.Options{Required: true, Help: "The type of node you want to run", Validate: func(args []string) error {
		if args[0] == config.Master || args[0] == config.Refinery || args[0] == config.Receiver || args[0] == config.Sub {
			return nil
		}
		return errors.New(fmt.Sprintf("NodeType %s not know. Any of: (`master`, `submaster`, `refinery`, `receiver`) accepted", args[0]))
	}})
	privateKey := parser.File("k", "key", os.O_RDONLY, 0600, &argparse.Options{Required: false, Help: "The path of the private key. To read as master node and generate if not exists", Default: "./id_rsa_test", Validate: func(args []string) error {
		if *nodeType != config.Master {
			return errors.New("only Master node is allowed a PrivateKey. No need for another node to have one")
		}
		return nil
	}})
	certificate := parser.File("c", "cert", os.O_RDONLY, 0600, &argparse.Options{Required: false, Help: "The path to the public certificate to connect to a master node", Default: "./id_rsa_test.pub"})
	logFile := parser.File("l", "log-file", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600, &argparse.Options{Required: false, Help: "The path to the log file", Default: "./log.txt"})
	logLevel := parser.Int("L", "log-level", &argparse.Options{Required: false, Help: fmt.Sprintf("The Log-Level for the logger used. Valid Levels:\n                   Debug: [%d]\n                   Info: [%d]\n                   Error: [%d]", util.LevelDebug, util.LevelInfo, util.LevelError), Default: util.LevelDebug, Validate: func(args []string) error {
		if args[0] == fmt.Sprintf("%d", util.LevelError) || args[0] == fmt.Sprintf("%d", util.LevelDebug) || args[0] == fmt.Sprintf("%d", util.LevelInfo) {
			return nil
		} else {
			return errors.New(fmt.Sprintf("Unkown Log Level `%s`. \nKnown Levels: \n Debug: [%d]\n Info: [%d]\n Error: [%d]", args[0], util.LevelDebug, util.LevelInfo, util.LevelError))
		}
	}})
	logType := parser.String("T", "log-type", &argparse.Options{Required: false, Help: fmt.Sprintf("The Log-Type of the logger, so how the logger will log. Known Log Types:\n                   Console Only: `%s`\n                   File Only: `%s`\n                   File and Console: `%s`", util.LogTypeConsole, util.LogTypeFile, util.LogTypeConFile), Default: util.LogTypeConsole, Validate: func(args []string) error {
		if args[0] == util.LogTypeFile || args[0] == util.LogTypeConsole || args[0] == util.LogTypeConFile {
			return nil
		} else {
			return errors.New(fmt.Sprintf("Unkown Log Type `%s`.\nKnown Log Types: \n Console Only: `%s`\n File Only: `%s`\n File and Console: `%s`", args[0], util.LogTypeConsole, util.LogTypeFile, util.LogTypeConFile))
		}
	}})
	masterHost := parser.String("H", "host", &argparse.Options{Required: false, Help: "The hostname of the master node you want to connect to", Validate: func(args []string) error {
		if *nodeType == config.Master {
			return errors.New("master does not need an other master node to connect to")
		} else {
			regexIP, _ := regexp.Compile("^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$")
			regexHostName, _ := regexp.Compile("^(([a-zA-Z0-9]|[a-zA-Z0-9][a-zA-Z0-9\\-]*[a-zA-Z0-9])\\.)*([A-Za-z0-9]|[A-Za-z0-9][A-Za-z0-9\\-]*[A-Za-z0-9])$")
			if regexIP.MatchString(args[0]) || regexHostName.MatchString(args[0]) {
				return nil
			} else {
				message := fmt.Sprintf("Hostname: %v is not a valid hostname", args[0])
				return errors.New(message)
			}
		}
	}})
	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(util.ArgumentErrorCode)
	}

	defer func(logFile *os.File) {
		err := logFile.Close()
		if err != nil {
			fmt.Println("Closing log file has resulted in an error")
		}
	}(logFile)

	logger := util.Logger{
		LogLevel: *logLevel,
		LogFile:  logFile,
		LogType:  *logType,
	}

	certPath, err := filepath.Abs(certificate.Name())
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to read Abs Path: %s", err.Error()), util.LevelError)
		os.Exit(util.ArgumentErrorCode)
	}
	keyPath, err := filepath.Abs(privateKey.Name())
	if err != nil {
		logger.Log(fmt.Sprintf("Failed to read Abs Path: %s", err.Error()), util.LevelError)
		os.Exit(util.ArgumentErrorCode)
	}

	err = certificate.Close()
	if err != nil {
		logger.Log(fmt.Sprintf("%s", err.Error()), util.LevelError)
		os.Exit(util.ArgumentErrorCode)
	}
	err = privateKey.Close()
	if err != nil {
		logger.Log(fmt.Sprintf("%s", err.Error()), util.LevelError)
		os.Exit(util.ArgumentErrorCode)
	}

	currentNode := node.NewNode(true)
	currentNode.NodeType = *nodeType
	currentNode.MasterHost = *masterHost
	currentNode.CertificatePath = certPath
	currentNode.KeyPath = keyPath

	if err != nil {
		logger.Log(err.Error(), util.LevelError)
		os.Exit(util.ArgumentErrorCode)
	}
	result, err := config.Construct(currentNode)
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
