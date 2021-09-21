package d7024e

const (
	PING      = "PING"
	FIND_NODE = "FIND_NODE"
	FIND_DATA = "FIND_DATA"
	STORE     = "STORE"
)

type Ping struct {
	Id      string
	Address string
}

type FindNode struct {
	Id      string
	Address string
}

type FindValue struct {
}

type Store struct {
}

func MsgHandler(channel chan []byte, me Contact, network Network) {

}
