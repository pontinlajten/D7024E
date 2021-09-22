package d7024e

const ( // RPC Operations
	PING      = "PING"
	FIND_NODE = "FIND_NODE"
	FIND_DATA = "FIND_DATA"
	STORE     = "STORE"
)

type Ping struct {
	Id      string // Kadmelia ID.
	RPC     string // RPC operation.
	Address string // IP Adress.
}

type FindNode struct {
	Id      string
	RPC     string
	Address string
}

type FindValue struct {
}

type Store struct {
}
