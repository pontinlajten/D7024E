package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"sync"
)

type Network struct {
	Me       *Contact
	Mutex    *sync.Mutex
	Rt       *RoutingTable
	Kademlia *Kademlia // Used in Listen.
}

const (
	CONN_TYPE       = "udp"
	MAX_BUFFER_SIZE = 1024
)

// Template for init. an network.
func createNetwork(me *Contact, rt *RoutingTable, kademlia *Kademlia) Network {
	network := Network{} // Create from Network struct
	network.Me = me
	network.Mutex = &sync.Mutex{}
	network.Rt = rt
	network.Kademlia = kademlia
	return network
}

/////////////////////////////// RESPONSE /////////////////////////////////////////

/*
	TODO: Maybe implement channel model instead of mutex. RESEARCH.
*/
// IN-PROGRESS
func (network *Network) Listen(port string) { // Listen(ip string, port int) original.
	raddr, err := net.ResolveUDPAddr(CONN_TYPE, port) // ResolveUDPAddr(str, str). me.Address
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
		msg := network.MsgHandler(buffer[:n])
		replyEncoded := marshall(msg)
		sendResponse(replyEncoded, addr, conn)
	}
}

func (network *Network) MsgHandler(data []byte) Message {
	decoded := unmarshall(data)
	reply := Message{}
	fmt.Println("RPC: " + decoded.RPC)

	if decoded.RPC == FIND_NODE {
		reply.Id = network.Me.ID.String()
		reply.RPC = FIND_NODE
		reply.Address = network.Me.Address
		reply.Data.Nodes = network.FindNodeHandler(decoded)

	} else if decoded.RPC == PING {
		reply.Id = network.Me.ID.String()
		reply.RPC = PONG
		reply.Address = network.Me.Address
		network.PingHandler(decoded)

	} else if decoded.RPC == FIND_DATA {
		reply = network.FindValueHandler(decoded)

	} else if decoded.RPC == STORE {
		network.StoreHandler(decoded)
		reply.Id = network.Me.ID.String()
		reply.RPC = STORE_REPLY
		reply.Address = network.Me.Address

	} else if decoded.RPC == PONG {

	}

	return reply
}

func (network *Network) PingHandler(msg Message) {
	id := NewKademliaID(msg.Id)
	newContact := NewContact(id, msg.Address)
	network.Rt.AddContact(newContact)
}

func (network *Network) FindNodeHandler(msg Message) []Contact {
	contacts := network.Rt.FindClosestContacts(NewKademliaID(msg.Id), ALPHA)
	return contacts
}

func (network *Network) FindValueHandler(msg Message) Message {
	keyVal := network.Kademlia.LookupData(msg.Data.Key)
	if keyVal != nil {
		newId := NewKademliaID(msg.Id)
		newContact := NewContact(newId, msg.Address)
		network.Rt.AddContact(newContact)
		return Message{Id: network.Me.ID.String(), RPC: FIND_DATA_REPLY, Address: network.Me.Address, Data: Data{Key: keyVal.Key, Value: keyVal.Value}}

	} else {
		id := NewKademliaID(msg.Data.Key)
		newContacts := network.Rt.FindClosestContacts(id, ALPHA)

		newId := NewKademliaID(msg.Id)
		newContact := NewContact(newId, msg.Address)
		network.Rt.AddContact(newContact)
		return Message{Id: network.Me.ID.String(), RPC: FIND_DATA_REPLY, Address: network.Me.Address, Data: Data{Nodes: newContacts}}
	}
}

func (network *Network) StoreHandler(msg Message) {
	network.Kademlia.Store(msg.Data.Value)
}

func sendResponse(responseMsg []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	_, err := conn.WriteToUDP([]byte(responseMsg), addr)
	if err != nil {
		fmt.Printf("Could'nt send response %v", err)
	}
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
	msg := Message{Id: network.Me.ID.String(), RPC: PING, Address: network.Me.Address}
	SendData(msg, contact)
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	msg := Message{Id: network.Me.ID.String(), RPC: FIND_NODE, Address: network.Me.Address}
	SendData(msg, contact)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) {
	msg := Message{Id: network.Me.ID.String(), Address: network.Me.Address, RPC: FIND_DATA, Data: Data{Key: hash}}
	SendData(msg, contact)
}

func (network *Network) SendStoreMessage(value string, contact *Contact) {
	msg := Message{Address: network.Me.Address, RPC: STORE, Data: Data{Value: value}}
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
