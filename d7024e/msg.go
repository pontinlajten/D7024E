package d7024e

const ( // RPC Operations
	PING      = "PING"
	FIND_NODE = "FIND_NODE"
	FIND_DATA = "FIND_DATA"
	STORE     = "STORE"

	PONG            = "PONG"
	FIND_NODE_REPLY = "FIND_NODE_REPLY"
	FIND_DATA_REPLY = "FIND_DATA_REPLY"
	STORE_REPLY     = "STORE_REPLY"
)

type Message struct {
	Id      string // Kadmelia ID.
	RPC     string // RPC operation.
	Address string // IP Adress.
	Data    Data
}

type Data struct {
	Nodes []Contact
	Key   string
	Value string
	Msg   string
}
