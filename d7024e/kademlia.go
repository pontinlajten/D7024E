package d7024e

import (
	"fmt"
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
}

type KeyValue struct {
	Key   string
	Value string
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
	net := &Network{}
	net.Kademlia = kademlia
	channel := make(chan []Contact)

	shortlist := kademlia.NewList(targetID)

	if shortlist.Len() < alpha {
		go AsyncFindContact(shortlist.Cons[0].Con, *targetID, *net, channel)
	} else {
		for i := 0; i < alpha; i++ {
			go AsyncFindContact(shortlist.Cons[i].Con, *targetID, *net, channel)
		}
	}

	shortlist.UpdateList(*targetID, channel, *net)

	for _, insItem := range shortlist.Cons {
		resultlist = append(resultlist, insItem.Con)
	}

	return
}

/*
	Use channels inorder to keep data from find_contact "safe". In terms of data write/read safety.
*/
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
	var wg sync.WaitGroup
	s_hash := string(hash)
	hashID := NewKademliaID(s_hash)

	shortlist := kademlia.NewList(hashID)

	ch := make(chan []Contact)
	targetData := make(chan []byte)
	dataContactCh := make(chan Contact)

	if shortlist.Len() < alpha {
		go asyncLookupData(hash, shortlist.Cons[0].Con, *net, ch, targetData, dataContactCh)
	} else {
		for i := 0; i < alpha; i++ {
			go asyncLookupData(hash, shortlist.Cons[i].Con, *net, ch, targetData, dataContactCh)
		}
	}

	data, con := shortlist.updateLookupData(hash, ch, targetData, dataContactCh, *net, wg)

	return data, con
}

/*
	Use channels inorder to keep data from find_value "safe". In terms of data write/read safety.
*/
func asyncLookupData(hash string, receiver Contact, net Network, ch chan []Contact, target chan []byte, dataContactCh chan Contact) {
	response, _ := net.SendFindDataMessage(hash, &receiver)
	ch <- response.Body.Nodes
	target <- []byte(response.Body.Value)
	dataContactCh <- *response.Source
}

/////////////////////////////// STORE RPC /////////////////////////////////////////

func (kademlia *Kademlia) Store(upload string) string {
	net := &Network{}
	net.Kademlia = kademlia
	hash := HashIt(upload)
	hashID := NewKademliaID(hash)

	k_desitnations := kademlia.LookupContact(hashID)
	var hashReturn string

	for _, target := range k_desitnations {
		response, _ := net.SendStoreMessage(upload, &target)
		if response.Body.Key != "" {
			hashReturn = response.Body.Key
		}
	}

	return hashReturn
}

func (kademlia *Kademlia) StoreKeyValue(value string) string {
	hash := HashIt(value)
	hashID := NewKademliaID(hash).String()
	for _, keyVal := range kademlia.KeyValues {
		if hash == keyVal.Key {
			fmt.Printf("Value is already existing")
			return keyVal.Key
		}
	}
	var newKeyValue KeyValue

	newKeyValue.Key = hashID
	newKeyValue.Value = value
	kademlia.KeyValues = append(kademlia.KeyValues, newKeyValue)

	fmt.Println("ALL STORED VALUES IN NODE: ")
	fmt.Println(kademlia.KeyValues)

	return newKeyValue.Key
}

//---------------------------------------------------------//

func (kademlia *Kademlia) InitNetwork(known *Contact) []Contact {
	kademlia.Rt.AddContact(*known) // Add bootstrap conctact
	contacts := kademlia.LookupContact(kademlia.Me.ID)

	//fmt.Printf("Joining network via %s", known.String())
	return contacts
}

//---------------------------------------------------------//
