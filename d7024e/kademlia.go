package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
)

const (
<<<<<<< HEAD
	k int = 20 // num of cont in bucket
	a int = 3  //(alpha) degree of parallelism in network calls
=======
	// fanns redan en bucketSize i rt //k int = 20 // num of cont in bucket
	a int = 3 //(alpha) degree of parallelism in network calls
>>>>>>> origin/peder
)

type Kademlia struct {
	id *KademliaID
	me Contact
	rt *RoutingTable
	//nt *Network
}

type LookedAt struct {
}

<<<<<<< HEAD
/* func NewKademlia(ip string) (kadNode Kademlia) {
=======
func NewKademlia(ip string) (kadNode Kademlia) {
>>>>>>> origin/peder
	kadNode.id = NewKademliaID(HashIt(ip))
	kadNode.me = NewContact(kadNode.id, ip)
	kadNode.rt = NewRoutingTable(kadNode.me)

<<<<<<< HEAD
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

=======
	return
}

//help function that hash data
func HashIt(str string) (data string) {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	data = hex.EncodeToString(hashStr.Sum(nil))
	return
>>>>>>> origin/peder
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//kClosest := kademlia.rt.FindClosestContacts(target.ID, k)
<<<<<<< HEAD
=======
	//findkclosest(target, bucketSize) []Contact
>>>>>>> origin/peder

	/*
		kClosest := kademlia.rt.FindClosestContacts(target.ID, 3)
		for i, c := range kClosest {
			if c.ID.Equals(target.ID) {

			}
		}
	*/
}

func (Kademlia *Kademlia) findkclosest(target *Contact, k int) []Contact {
	Kclosest := Kademlia.rt.FindClosestContacts(target.ID, k)
	return Kclosest
}

func (kademlia *Kademlia) LookupData(hash string) {
	// TODO
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}
