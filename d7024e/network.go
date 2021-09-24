package d7024e

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
	kademlia *Kademlia // Used in Listen.
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

/////////////////////////////// RESPONSE /////////////////////////////////////////

// IN-PROGRESS
func (network *Network) Listen(ip string, port int, node Kademlia) { // Listen(ip string, port int) original.
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
		msg := network.MsgHandler(buffer[:n], conn, node)
		replyEncoded := marshall(msg)
		sendResponse(replyEncoded,addr,conn)
	}
}


func (network *Network) MsgHandler(data []byte, conn *net.UDPConn, node Kademlia) Message {
	decoded := unmarshall(data)
	var reply Message
	fmt.Println("RPC: " + decoded.RPC)

	if decoded.RPC == FIND_NODE {
		contacts := node.rt.FindClosestContacts(NewKademliaID(decoded.Id), ALPHA)
		//msg.Data.Nodes = contacts
	} else if decoded.RPC == PING {
		network.SendPongMessage(decoded)
		reply = Message{Id: network.me.ID.String(), RPC: PONG, Address: network.me.Address}
	}

	return reply
}

func sendResponse(responseMsg []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	_, err := conn.WriteToUDP([]byte(responseMsg), addr)
	if err != nil {
		fmt.Printf("Could'nt send response %v", err)
	}
}

func (network *Network) SendPongMessage(msg Message) {
	id := NewKademliaID(msg.Id)
	newContact := NewContact(id, msg.Address)
	network.rt.AddContact(newContact)
}


//////////////////////////////////////////////////////////////////////////////////

func SendData(msg Message, contact *Contact) {
	var rpcMsg string
	sendMsg := marshall(msg)
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
		fmt.Printf("failed to %s to %s error: %s", rpcMsg, contact.Address, err)
	}
}

func (network *Network) SendPingMessage(contact *Contact) {
	msg := Message{Id: network.me.ID.String(), RPC: PING, Address: network.me.Address}
	SendData(msg, contact)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	msg := Message{Id: network.me.ID.String(), RPC: FIND_NODE, Address: network.me.Address}
	SendData(msg, contact)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) {
	msg := Message{Address: network.me.Address, RPC: FIND_DATA, Data: Data{Key: hash}}
	SendData(msg, contact)
}

func (network *Network) SendStoreMessage(data string, contact *Contact) {
	hash := network.kademlia.HashIt(data)
	msg := Message{Address: network.me.Address, RPC: STORE, Data: Data{Key: hash, Value: data}}
	SendData(msg, contact)
}

/////////////////////// HELP FUNCTIONS //////////////////////////

func marshall(msg Message) []byte {
	encoded, _ := json.Marshal(msg)
	return encoded
}

func unmarshall(data []byte) Message {
	var decoded Message
	json.Unmarshal([]byte(data), &decoded)
	return decoded
}
