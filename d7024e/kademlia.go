package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
)

const (
	k int = 20 // num of cont in bucket
	a int = 3  //(alpha) degree of parallelism in network calls
)

type Kademlia struct {
	id *KademliaID
	me Contact
	rt *RoutingTable
	//nt *Network
}

func NewKademlia(ip string) (kadNode Kademlia) {
	//kademlia := Kademlia{}
	kadNode.id = NewKademliaID(ip)
	kadNode.me = NewContact(kadNode.id, ip)
	kadNode.rt = NewRoutingTable(kadNode.me)
	//kademlia.nt = createNetwork(me, rt, kademlia)
	return
}

//help function that has data
func hashIt(str string) (hash string) {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash = hex.EncodeToString(hashStr.Sum(nil))
	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	// TODO
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
