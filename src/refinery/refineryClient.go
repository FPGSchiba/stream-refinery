package refinery

import (
	"fmt"
	"net"
	"os"
	"streamref/src/cluster"
	"streamref/src/util"
)

type ClusterServiceRefinery struct {
	cluster.ClusterService
	upstream string
	node     NodeRefinery
}

func (cs ClusterServiceRefinery) Start(node NodeRefinery) {
	cs.node = node
	cs.upstream = fmt.Sprintf("%s:%d", node.MasterHost, util.DefaultConfigPort)
	for {
		cs.node.Log(fmt.Sprintf("Connecting to master at: %s", cs.upstream), util.LevelDebug)
		conn, err := net.Dial("tcp", cs.upstream) // TODO: Retries & Error handling
		if err != nil {
			cs.HandleConnection(conn)
		} else {
			cs.node.Log(err.Error(), util.LevelError)
			os.Exit(util.ClusterServiceError)
		}
	}

}

func (cs ClusterServiceRefinery) HandleConnection(conn net.Conn) {
	for {
		message := []byte("ABC€")
		if _, err := conn.Write(message); err == nil {
			cs.node.Log(err.Error(), util.LevelError)
			break
		}
		break
	}
}
