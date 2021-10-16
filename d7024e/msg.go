package d7024e

const ( // RPC Operations
	PING      = "PING"
	FIND_NODE = "FIND_NODE"
	FIND_DATA = "FIND_DATA"
	STORE     = "STORE"

	PING_REPLY      = "PING_REPLY"
	FIND_NODE_REPLY = "FIND_NODE_REPLY"
	FIND_DATA_REPLY = "FIND_DATA_REPLY"
	STORE_REPLY     = "STORE_REPLY"
)

type Message struct {
	Source *Contact
	RPC    string // RPC operation.
	Body   Data
}

type Data struct {
	Nodes          []Contact
	Key            string
	Value          string
	OriginalSource *Contact    // From original sender, used in validator.
	TargetId       *KademliaID // TargetID when checking if exists and so on.
}
