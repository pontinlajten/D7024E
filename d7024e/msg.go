package d7024e

const (
	PING      = "PING"
	FIND_NODE = "FIND_NODE"
	FIND_DATA = "FIND_DATA"
	STORE     = "STORE"
)

type Ping struct {
}

type FindNode struct {
}

type FindValue struct {
}

type Store struct {
}

func MsgHandler(channel chan []byte, me Contact, network Network) {

}
