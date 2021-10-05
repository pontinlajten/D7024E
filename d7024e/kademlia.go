package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

const (
	// fanns redan en bucketSize i rt //k int = 20 // num of cont in bucket
	ALPHA = 3 //(alpha) degree of parallelism in network calls
)

type Kademlia struct {
	Id        *KademliaID
	Me        Contact
	Rt        *RoutingTable
	KeyValues []KeyValue
}

type KeyValue struct {
	Key   string
	Value string
}

func NewKademlia(ip string) (kademlia Kademlia) {
	kademlia.Id = NewKademliaID(kademlia.HashIt(ip))
	kademlia.Me = NewContact(kademlia.Id, ip)
	kademlia.Rt = NewRoutingTable(kademlia.Me)
	return
}

func (kademlia *Kademlia) LookupContact(target *Contact) {
	//kClosest := kademlia.rt.FindClosestContacts(target.ID, k)
	kademlia.FindKClosest(target, bucketSize)

	return
}

func (kademlia *Kademlia) FindKClosest(target *Contact, k int) []Contact {
	Kclosest := kademlia.Rt.FindClosestContacts(target.ID, k)
	return Kclosest
}

func (kademlia *Kademlia) LookupData(hash string) *KeyValue {
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			return &keyVal
		}
	}
	return nil
}

func (kademlia *Kademlia) Store(data []byte) {
	// TODO
}

/*
		kClosest := kademlia.rt.FindClosestContacts(target.ID, 3)
		for i, c := range kClosest {
			if c.ID.Equals(target.ID) {

		}
	}
*/

func (kademlia *Kademlia) InitRt(known *Contact) {
	kademlia.Rt.AddContact(*known)
	kademlia.LookupContact(&kademlia.Me)
	fmt.Printf("Kademlia node joining network")
}

//help function that hash data
func (kademlia *Kademlia) HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))
	//fmt.Println(hash)
	return hash

}

func HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))
	//fmt.Println(hash)
	return hash
}
