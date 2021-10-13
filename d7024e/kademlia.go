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
	//TimeStamp int
}

func NewKademlia(ip string) (kademlia Kademlia) {
	kademlia.Id = NewKademliaID(kademlia.HashIt(ip))
	kademlia.Me = NewContact(kademlia.Id, ip)
	kademlia.Rt = NewRoutingTable(kademlia.Me)
	return
}

//---------------------------------------------------------//

func (kademlia *Kademlia) LookupContact(targetID *KademliaID) (resultlist []Contact) {
	//ch := make(chan []Contact)
	net := &Network{}
	net.Kademlia = kademlia
	channel := make(chan []Contact)

	// shortlist of k-closest nodes
	shortlist := kademlia.NewList(targetID)

	// if LookupContact on JoinNetwork
	if shortlist.Len() < ALPHA {
		go reciverResponse(shortlist.Cons[0].Con, *net, channel)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < ALPHA; i++ {
			go reciverResponse(shortlist.Cons[i].Con, *net, channel)
		}
	}

	shortlist.UpdateList(*targetID, channel, *net)

	// creating the result list
	for _, insItem := range shortlist.Cons {
		resultlist = append(resultlist, insItem.Con)
	}
	return
}

func reciverResponse(reciver Contact, net Network, channel chan []Contact) {
	response, _ := net.SendFindContactMessage(&reciver)
	channel <- response
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


func (kademlia *Kademlia) Store (upload string) {
	network := &Network{}
	destContacts := kademlia.LookupContact(&kademlia.Me)
	for _, destContact := range destContacts {
		network.SendStoreMessage(upload, &destContact)
	}
}


//---------------------------------------------------------//
<<<<<<< HEAD
func (kademlia *Kademlia) Store(upload string) {
	network := &Network{}
	destContacts := kademlia.LookupContact(&kademlia.Me)
	for _, destContact := range destContacts {
		network.SendStoreMessage(upload, &destContact)
	}
}

=======
>>>>>>> 3c79d2bc29abd1f7aa4b4fb519570365fc81ff7a
func (kademlia *Kademlia) StoreKeyValue(value string) {
	hash := HashIt(value)
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			//keyVal.TimeStamp = REBUPLISH
			fmt.Printf("Value is already existing")
			return
		}
	}
	var newKeyValue KeyValue
	newKeyValue.Key = hash
	newKeyValue.Value = value
	//newKeyValue.TimeStamp = 24
	kademlia.KeyValues = append(kademlia.KeyValues, newKeyValue)
}

//---------------------------------------------------------//
func (kademlia *Kademlia) InitRt(known *Contact) {
	kademlia.Rt.AddContact(*known)
	kademlia.LookupContact(kademlia.Me.ID)
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
