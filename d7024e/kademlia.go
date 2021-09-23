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

type LookedAt struct {
}

/* func NewKademlia(ip string) (kadNode Kademlia) {
	kadNode.id = NewKademliaID(HashIt(ip))
	kadNode.me = NewContact(kadNode.id, ip)
	kadNode.rt = NewRoutingTable(kadNode.me)

	fmt.Println("")
	fmt.Println(kadNode.id)
	fmt.Println("")
	fmt.Println(kadNode.me)
	fmt.Println("")
	fmt.Println(kadNode.rt)
	return
} */

//help function that hash data
func (kademlia *Kademlia) HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))
	//fmt.Println(hash)
	return hash

}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//kClosest := kademlia.rt.FindClosestContacts(target.ID, k)

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
