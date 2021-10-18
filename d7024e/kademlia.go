package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"sync"
)

const alpha int = 3

const (
	// fanns redan en bucketSize i rt //k int = 20 // num of cont in bucket
	K = 20 // num of cont in bucket
)

type Kademlia struct {
	Id        *KademliaID
	Me        Contact
	Rt        *RoutingTable
	KeyValues []KeyValue
	Log       *log.Logger
}

type KeyValue struct {
	Key   string
	Value string
	//TimeStamp int
}

func NewKademlia(ip string) (kademlia Kademlia) {
	kademlia.Id = NewKademliaID(kademlia.HashIt(ip))
	kademlia.Me = NewContact(kademlia.Id, ip)
	kademlia.Rt = NewRoutingTable(kademlia.Me)
	kademlia.Me.Address = ip

	return kademlia
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
	if shortlist.Len() < alpha {
		go AsyncFindContact(shortlist.Cons[0].Con, *targetID, *net, channel)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go AsyncFindContact(shortlist.Cons[i].Con, *targetID, *net, channel)
		}
	}

	shortlist.UpdateList(*targetID, channel, *net)

	// creating the result list
	for _, insItem := range shortlist.Cons {
		resultlist = append(resultlist, insItem.Con)
	}

	//kademlia.Log.Printf("Looking up contact %s and found closest %s.", targetID.String(), resultlist)
	return
}

func AsyncFindContact(reciver Contact, targetID KademliaID, net Network, channel chan []Contact) {
	response, err := net.SendFindContactMessage(&reciver, &targetID)
	if err != nil {
		fmt.Println(err)
	}
	channel <- response
}

//---------------------------------------------------------//

func (kademlia *Kademlia) LookupDataHash(hash string) *KeyValue {
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			return &keyVal
		}
	}
	return nil
}

func (kademlia *Kademlia) LookupData(hash string) ([]byte, Contact) {
	net := &Network{}
	net.Kademlia = kademlia
	var wg sync.WaitGroup // gorutine waiting pool

	hashID := NewKademliaID(hash) // create kademlia ID from the hashed data
	/*
		shortlist (below) is a LookupList which both contains the contacts
		that need to be traversed in order to find the data as well
		as data itself.
	*/

	shortlist := kademlia.NewList(hashID)

	ch := make(chan []Contact)          // channel -> returns contacts
	targetData := make(chan []byte)     // channel -> when the data is found it is communicated through this channel
	dataContactCh := make(chan Contact) // channel that only takes the contact that returned the data

	if shortlist.Len() < alpha {
		go asyncLookupData(hash, shortlist.Cons[0].Con, *net, ch, targetData, dataContactCh)
	} else {
		// sending RPCs to the alpha nodes async
		for i := 0; i < alpha; i++ {
			go asyncLookupData(hash, shortlist.Cons[i].Con, *net, ch, targetData, dataContactCh)
		}
	}

	data, con := shortlist.updateLookupData(hash, ch, targetData, dataContactCh, *net, wg)

	// creating the resultdata, con :=shortlist.updateLook list
	return data, con
}

func asyncLookupData(hash string, receiver Contact, net Network, ch chan []Contact, target chan []byte, dataContactCh chan Contact) {
	response, _ := net.SendFindDataMessage(&receiver, hash)
	ch <- response.Body.Nodes
	target <- targetData
	dataContactCh <- dataContact
}

func (kademlia *Kademlia) StoreKeyValue(value string) string {
	hash := HashIt(value)
	hashID := NewKademliaID(hash).String()
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			//keyVal.TimeStamp = REBUPLISH
			fmt.Printf("Value is already existing")
			return keyVal.Key
		}
	}
	var newKeyValue KeyValue

	newKeyValue.Key = hashID
	newKeyValue.Value = hash
	//newKeyValue.TimeStamp = 24
	kademlia.KeyValues = append(kademlia.KeyValues, newKeyValue)

	return newKeyValue.Key
}

func (kademlia *Kademlia) Store(upload string) []Contact {
	net := &Network{}
	net.Kademlia = kademlia
	hash := HashIt(upload)
	hashID := NewKademliaID(hash)

	k_desitnations := kademlia.LookupContact(hashID)
	fmt.Println("STORING AT: ")
	fmt.Println(k_desitnations)

	var hashList []string

	for _, target := range k_desitnations { // Checks shortlist for k-nearest.
		response, _ := net.SendStoreMessage(upload, &target)
		hash := response.Body.Key
		hashList = append(hashList, hash)
	}
	// resp, _ = net.SendStoreMessageIP(upload, ip)
	return k_desitnations
}

//---------------------------------------------------------//

func (kademlia *Kademlia) InitNetwork(known *Contact) []Contact {
	kademlia.Rt.AddContact(*known) // Add bootstrap conctact
	contacts := kademlia.LookupContact(kademlia.Me.ID)

	//fmt.Printf("Joining network via %s", known.String())
	return contacts
}

//help function that hash data
func (kademlia *Kademlia) HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))

	return hash
}

func HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))

	return hash
}

//---------------------------------------------------------//
