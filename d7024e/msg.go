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
	Source *Contact // Contact constructing message. (In listener or send)
	RPC    string   // RPC operation.
	Body   MsgBody
}

type MsgBody struct {
	Nodes          []Contact // All nodes from recieve.
	Key            string    // Primarly store and find_val rpc.
	Value          string
	OriginalSource *Contact    // From original sender, used in validator.
	TargetId       *KademliaID // TargetID when checking if exists and so on.
}
