package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

const (
	// fanns redan en bucketSize i rt //k int = 20 // num of cont in bucket
	ALPHA     = 3 //(alpha) degree of parallelism in network calls
	REBUPLISH = 24
	K         = 20 // num of cont in bucket
)

type Kademlia struct {
	Id        *KademliaID
	Me        Contact
	Rt        *RoutingTable
	KeyValues []KeyValue
}

type KeyValue struct {
	Key       string
	Value     string
	TimeStamp int
}

func NewKademlia(ip string) (kademlia Kademlia) {
	kademlia.Id = NewKademliaID(kademlia.HashIt(ip))
	kademlia.Me = NewContact(kademlia.Id, ip)
	kademlia.Rt = NewRoutingTable(kademlia.Me)
	return
}

//---------------------------------------------------------//

func (kademlia *Kademlia) LookupContact(target *Contact) []Contact {
	//ch := make(chan []Contact)
	net := &Network{}
	net.Kademlia = kademlia
	List := kademlia.FindXClosest(target, K)
	aClosest := kademlia.FindXClosest(target, ALPHA)
	if len(List) < ALPHA {

	}

	return aClosest
}

//---------------------------------------------------------//

func (kademlia *Kademlia) LookupData(hash string) *KeyValue {
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			return &keyVal
		}
	}
	return nil
}

//---------------------------------------------------------//
func (kademlia *Kademlia) Store(value string) {
	hash := HashIt(value)
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			keyVal.TimeStamp = REBUPLISH
			return
		}
	}
	var newKeyValue KeyValue
	newKeyValue.Key = hash
	newKeyValue.Value = value
	newKeyValue.TimeStamp = 24
	kademlia.KeyValues = append(kademlia.KeyValues, newKeyValue)
}

//---------------------------------------------------------//
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

//---------------------------------------------------------//
