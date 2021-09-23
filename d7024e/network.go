package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"net/rpc"
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

//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

// IN-PROGRESS
func (network *Network) Listen(ip string, port int) { // Listen(ip string, port int) original.
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
		msg := MsgHandler(buffer[:n], addr, conn)
		marshalledMsg := marshall(msg)
		sendResponse(marshalledMsg, addr, conn)

		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
	}
}

func MsgHandler(data []byte, addr *net.UDPAddr, conn *net.UDPConn) Message {
	decoded := unmarshall(data)

	fmt.Println("RPC: " + decoded.RPC)

	msg := Message{}
	msg.RPC = decoded.RPC // RPC operation
	msg.Id = decoded.Id   // Kademlia id represented as a string
	msg.Body = Data{}     // Body data

	if decoded.RPC == FIND_DATA || decoded.RPC == FIND_NODE {
		contacts := rt.find
	}

	return msg
}

func sendResponse(responseMsg []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	_, err := conn.WriteToUDP([]byte(responseMsg), addr)
	if err != nil {
		fmt.Printf("Could'nt send response %v", err)
	}
}

//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////
//////////////////////////////////////////////////////////////////////////////////

func (network *Network) SendPingMessage(contact *Contact, destination net.UDPAddr) {
	var pingReply Ping
	msg := Ping{Id: contact.ID.String(), Address: contact.Address}
	client, err := rpc.Dial("udp", destination.String())
	if err != nil {
		fmt.Printf("failed to dial %s error: %s", destination.String(), err)
	}

	err = client.Call(PING, msg, pingReply)
	if err != nil {
		fmt.Printf("failed to PING %s error: %s", destination.String(), err)
	} else {
		newId := NewKademliaID(pingReply.Id)
		newContact := NewContact(newId, pingReply.Address)
		network.rt.AddContact(newContact)
	}
}

func (network *Network) SendFindContactMessage(contact *Contact, destination net.UDPAddr) []Contact {
	var findNodeReply FindNode
	msg := FindNode{Id: contact.ID.String(), Address: contact.Address}
	client, err := rpc.Dial("udp", destination.String())
	if err != nil {
		fmt.Printf("failed to dial %s error: %s", destination.String(), err)
		return nil
	}

	err = client.Call(FIND_NODE, msg, findNodeReply)
	if err != nil {
		fmt.Printf("failed to FIND NODE %s error: %s", destination.String(), err)
		return nil
	}

}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
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
