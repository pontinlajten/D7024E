package d7024e

// creates New list with x closest
func (kademlia *Kademlia) FindXClosest(target *Contact, x int) []Contact {
	xClosest := kademlia.Rt.FindClosestContacts(target.ID, x)
	return xClosest
}

func UpdateList() {

}

func RecieverResponse() {

}
