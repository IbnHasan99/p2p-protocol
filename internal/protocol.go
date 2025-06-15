package internal

const (
	TaskProtocolID = "/p2p/task/"
	RoleProtocolID = "/p2p/role/"
	DiscoveryTag   = "p2p-network"
)

type TaskMessage struct {
	Instruction string `json:"instruction"`
}
