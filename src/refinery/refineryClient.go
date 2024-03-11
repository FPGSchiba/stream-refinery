package refinery

import (
	"fmt"
	"net"
	"os"
	"streamref/src/cluster"
	"streamref/src/util"
	"time"
)

type ClusterServiceRefinery struct {
	cluster.ClusterService
	upstream string
	node     NodeRefinery
}

func (cs ClusterServiceRefinery) Start(node NodeRefinery) {
	cs.node = node
	cs.upstream = fmt.Sprintf("%s:%d", node.MasterHost, util.DefaultClusterPort)
	retries := 5
	failed := true
	for i := 0; i <= retries; i++ {
		cs.node.Log(fmt.Sprintf("Connecting to master at: %s", cs.upstream), util.LevelDebug)
		conn, err := net.Dial("tcp", cs.upstream)
		if err == nil && conn != nil {
			connectionError := cs.HandleConnection(conn)
			if connectionError != nil {
				cs.node.Log(err.Error(), util.LevelError)
			} else {
				failed = false
			}
			break
		} else {
			cs.node.Log(err.Error(), util.LevelError)
		}
		time.Sleep(5 * time.Second)
	}
	if failed {
		os.Exit(util.ClusterServiceError)
	}
}

func (cs ClusterServiceRefinery) HandleConnection(conn net.Conn) error {
	defer conn.Close()
	for {
		err := Authenticate(conn, cs.node.NodeID)
		if err != nil {
			return err
		}
		return nil
	}
}
