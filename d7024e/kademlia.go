package d7024e

import (
	"fmt"
)

type Kademlia struct {
	me      Contact
	table   *RoutingTable
	network *Network
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
