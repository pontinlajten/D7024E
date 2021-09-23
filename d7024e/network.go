package d7024e

import "C"
import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Network struct {
	me       *Contact
	mutex    *sync.Mutex
	rt       *RoutingTable
	kademlia *Kademlia
}

const (
	CONN_TYPE       = "udp"
	MAX_BUFFER_SIZE = 1024
)


// Template for init. an network.
func createNetwork(me *Contact, rt *RoutingTable, kademlia *Kademlia) Network {
	network := Network{} // Create from Network struct
	network.me = me
	network.mutex = &sync.Mutex{}
	network.rt = rt
	network.kademlia = kademlia
	return network
}

func MsgHandler(channel chan []byte, me Contact, network Network) {

}


func (network *Network) PingHandler(msg Message) {
	newId := NewKademliaID(msg.Id)
	newContact := NewContact(newId, msg.Address)
	network.rt.AddContact(newContact)
}

// IN-PROGRESS
func (network *Network) Listen(me Contact, port int) { // Listen(ip string, port int) original.
	raddr, err := net.ResolveUDPAddr(CONN_TYPE, ":8080") // ResolveUDPAddr(str, str). me.Address
	conn, err2 := net.ListenUDP(CONN_TYPE, raddr)
	if (err != nil) || (err2 != nil) {
		fmt.Println("Error udp: ", err, "    ", err2)
	}

	defer conn.Close() // defer: Close last, after all functions execution below is done.
	//channel := make(chan []byte)
	buffer := make([]byte, MAX_BUFFER_SIZE) // Recieve ASCII, byte representation.

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error ReadFromUDP", err)
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
	}
}

func SendData (msg Message, contact *Contact) {
	var rpcMsg string
	sendMsg, err := json.Marshal(msg)
	Client, err := net.Dial("udp", contact.Address)
	if err != nil {
		fmt.Printf("failed to dial %s error: %s", contact.Address, err)
	}
	switch msg.RPC {
	case PING:
		rpcMsg = PING

	case FIND_NODE:
		rpcMsg = FIND_NODE

	case FIND_DATA:
		rpcMsg = FIND_DATA

	case STORE:
		rpcMsg = STORE
	}
	defer Client.Close()
	_, err = Client.Write(sendMsg)
	if err != nil {
		fmt.Printf("failed to %s to %s error: %s",rpcMsg ,contact.Address, err)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	msg := Message{Id: network.me.ID.String(),RPC: PING, Address: network.me.Address}
	SendData(msg,contact)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	msg := Message{Id: network.me.ID.String(), RPC: FIND_NODE, Address: network.me.Address}
	SendData(msg,contact)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) {
	msg := Message{Address: network.me.Address,RPC: FIND_DATA, data: Data{Key: hash}}
	SendData(msg,contact)
}

func (network *Network) SendStoreMessage(data string, contact *Contact) {
	hash := network.kademlia.HashIt(data)
	msg := Message{Address: network.me.Address,RPC: STORE, data: Data{Key: hash, Value: data}}
	SendData(msg,contact)
}
