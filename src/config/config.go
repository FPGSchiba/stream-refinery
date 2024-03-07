package config

const (
	Master   = "master"
	Sub      = "submaster"
	Refinery = "refinery"
	Receiver = "receiver"
)

type Node struct {
	NodeType   string
	MasterHost string
}
