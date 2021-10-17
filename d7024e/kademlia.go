package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"log"
	"os"
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
	fmt.Println("MY IDIDIDIDIDID")
	fmt.Println(kademlia.Id)
	kademlia.Me = NewContact(kademlia.Id, ip)
	fmt.Println(kademlia.Me)
	kademlia.Rt = NewRoutingTable(kademlia.Me)
	kademlia.Me.Address = ip

	file, err := os.OpenFile("node_log.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	kademlia.Log = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	kademlia.Log.Printf("Node %s created on address %s \n", kademlia.Me.ID.String(), kademlia.Me.Address)

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

	fmt.Println("The length is asffasdfdsa")
	fmt.Println(shortlist.Len())
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
	fmt.Println("asyncfunctionnnn")
	fmt.Println(response)
	if err != nil {
		fmt.Println(err)
	}
	channel <- response
}

//---------------------------------------------------------//

func (kademlia *Kademlia) LookupData(value string) *KeyValue {
	ifExist := HashIt(value)
	for _, keyVal := range kademlia.KeyValues {
		if ifExist == keyVal.Value {
			return &keyVal
		}
	}
	return nil
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

	for _, target := range k_desitnations { // Checks shortlist for k-nearest.
		net.SendStoreMessage(upload, &target)
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
