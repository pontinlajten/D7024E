package d7024e

import (
	"sync"
)

type Lookup struct {
	Con   []Item
	Mutex sync.Mutex
}

type List struct {
	Cons  []Item
	Mutex sync.Mutex
}

type Item struct {
	Con  Contact
	Seen bool
}

// creates New list with x closest
func (kademlia *Kademlia) FindXClosest(target *Contact, x int) []Contact {
	xClosest := kademlia.Rt.FindClosestContacts(target.ID, x)
	return xClosest
}

func (kademlia *Kademlia) NewList(target *Contact) (list *List) {
	list = &List{}
	klist := kademlia.FindXClosest(target, bucketSize)

	for _, item := range klist {
		listitem := &Item{item, false}
		list.Cons = append(list.Cons, *listitem)
	}
	return
}

func (list *List) UpdateList(ID KademliaID, ch chan []Contact, net Network) {

}

func RecieverResponse() {

}
