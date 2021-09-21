package d7024e

type Kademlia struct {
	me      Contact
	rt      *RoutingTable
	network *Network
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	/*
		kClosest := kademlia.rt.FindClosestContacts(target.ID, 3)
		for i, c := range kClosest {
			if c.ID.Equals(target.ID) {

			}
		}
	*/
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
