package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"
)

type Network struct {
	Mutex    *sync.Mutex
	Kademlia *Kademlia // Used in Listen.
}

const (
	CONN_TYPE       = "udp"
	MAX_BUFFER_SIZE = 1024
)

// Template for init. an network.
func CreateNetwork(kademlia *Kademlia) Network {
	network := Network{} // Create from Network struct
	network.Mutex = &sync.Mutex{}
	network.Kademlia = kademlia
	return network
}

/////////////////////////////// RESPONSE /////////////////////////////////////////

/*
	TODO: Maybe implement channel model instead of mutex. RESEARCH.
*/
// IN-PROGRESS
func (network *Network) Listen(port string) { // Listen(ip string, port int) original.
	raddr, err := net.ResolveUDPAddr(CONN_TYPE, ":"+port) // ResolveUDPAddr(str, str). me.Address
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
		msg, ifSend := network.MsgHandler(buffer[:n])
		replyEncoded := marshall(msg)
		if ifSend {
			sendResponse(replyEncoded, addr, conn)
		}
	}
}

func (network *Network) MsgHandler(data []byte) (Message, bool) {
	decoded := unmarshall(data)
	reply := Message{}
	body := Data{}
	fmt.Println("RPC: " + decoded.RPC)

	var ifSend bool

	// HANDLING REQUEST->NODE
	if decoded.RPC == FIND_NODE {
		reply.Id = network.Kademlia.Me.ID.String()
		reply.RPC = FIND_NODE_REPLY
		reply.Address = network.Kademlia.Me.Address

		body.Nodes = network.FindNodeHandler(decoded)

		ifSend = true
	} else if decoded.RPC == PING {
		reply.Id = network.Kademlia.Me.ID.String()
		reply.RPC = PING_REPLY
		reply.Address = network.Kademlia.Me.Address

		ifSend = true
	} else if decoded.RPC == FIND_DATA {
		reply = network.FindValueHandler(decoded)

		ifSend = true
	} else if decoded.RPC == STORE {
		reply.Id = network.Kademlia.Me.ID.String()
		reply.RPC = STORE_REPLY
		reply.Address = network.Kademlia.Me.Address

		body.Key = network.StoreHandler(decoded)

		ifSend = true
	}

	body.RequestId = decoded.Id
	reply.Body = body

	return reply, ifSend
}

func (network *Network) FindNodeHandler(msg Message) []Contact {
	contacts := network.Kademlia.Rt.FindClosestContacts(msg.Body.TargetId, bucketSize)

	return contacts
}

func (network *Network) FindValueHandler(msg Message) Message {
	keyVal := network.Kademlia.LookupData(msg.Body.Key)
	if keyVal != nil {
		return Message{Id: network.Kademlia.Me.ID.String(), RPC: FIND_DATA_REPLY, Address: network.Kademlia.Me.Address, Body: Data{Key: keyVal.Key, Value: keyVal.Value}}
	} else {
		id := NewKademliaID(msg.Body.Key)
		newContacts := network.Kademlia.Rt.FindClosestContacts(id, ALPHA)

		return Message{Id: network.Kademlia.Me.ID.String(), RPC: FIND_DATA_REPLY, Address: network.Kademlia.Me.Address, Body: Data{Nodes: newContacts}}
	}
}

func (network *Network) StoreHandler(msg Message) string {
	key := network.Kademlia.StoreKeyValue(msg.Body.Value)
	return key
}

func sendResponse(responseMsg []byte, addr *net.UDPAddr, conn *net.UDPConn) {
	_, err := conn.WriteToUDP([]byte(responseMsg), addr)
	if err != nil {
		fmt.Printf("Could'nt send response %v", err)
	}
}

//////////////////////////////////////////////////////////////////////////////////

func (network *Network) SendData(msg Message, contact *Contact) (Message, error) {
	var rpcMsg string
	sendMsg := marshall(msg)

	address := GetUDPAddrFromContact(contact)
	Client, err := net.DialUDP("udp", nil, &address)
	if err != nil {
		fmt.Printf("failed to dial %s error: %s", address.String(), err)
	}
	defer Client.Close()

	fmt.Println(msg.RPC + " SEND MESSAGE")

	switch msg.RPC {
	case PING:
		fmt.Println("PING SEND MESSAGE")
		rpcMsg = PING

	case FIND_NODE:
		rpcMsg = FIND_NODE

	case FIND_DATA:
		rpcMsg = FIND_DATA

	case STORE:
		rpcMsg = STORE
	}

	Client.Write([]byte(sendMsg))

	buf := make([]byte, MAX_BUFFER_SIZE)

	n, _, _ := Client.ReadFromUDP(buf)
	response := unmarshall(buf[0:n])

	if err != nil {
		fmt.Printf("failed to %s to %s error: %s", rpcMsg, address.String(), err)
	}

	if Validate(msg, response) {
		network.Kademlia.Rt.AddContact(NewContact(NewKademliaID(response.Id), response.Address)) // Updates routing table if recieves succesful respond from target node.
	}

	return response, nil
}

func SendDataIP(msg Message, ip string) (Message, error) {
	var rpcMsg string
	sendMsg := marshall(msg)

	Client, err := net.Dial("udp", ip)
	if err != nil {
		fmt.Printf("failed to dial %s error: %s", ip, err)
	}

	fmt.Println(msg.RPC + " SEND MESSAGE")

	switch msg.RPC {
	case PING:
		fmt.Println("PING SEND MESSAGE")
		rpcMsg = PING

	case FIND_NODE:
		rpcMsg = FIND_NODE

	case FIND_DATA:
		rpcMsg = FIND_DATA

	case STORE:
		rpcMsg = STORE
	}
	defer Client.Close()

	Client.Write(sendMsg)

	buf := make([]byte, 2048)

	n, _ := Client.Read(buf)
	response := unmarshall(buf[0:n])

	if err != nil {
		fmt.Printf("failed to %s to %s error: %s", rpcMsg, ip, err)
	}

	return response, nil
}

func (network *Network) SendPingMessage(contact *Contact) (Message, error) {
	msg := Message{Id: network.Kademlia.Me.ID.String(), RPC: PING, Address: network.Kademlia.Me.Address}
	return network.SendData(msg, contact)
}

func (network *Network) SendFindContactMessage(contact *Contact, targetId *KademliaID) (Message, error) {
	msg := Message{Id: network.Kademlia.Me.ID.String(), RPC: FIND_NODE, Address: network.Kademlia.Me.Address, Body: Data{TargetId: targetId}}
	return network.SendData(msg, contact)
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) (Message, error) {
	msg := Message{Id: network.Kademlia.Me.ID.String(), Address: network.Kademlia.Me.Address, RPC: FIND_DATA, Body: Data{Key: hash}}
	return network.SendData(msg, contact)
}

func (network *Network) SendStoreMessage(value string, contact *Contact) (Message, error) {
	msg := Message{Address: network.Kademlia.Me.Address, RPC: STORE, Body: Data{Value: value}}
	return network.SendData(msg, contact)
}

func (network *Network) SendStoreMessageIP(value string, ip string) (Message, error) {
	msg := Message{Address: network.Kademlia.Me.Address, RPC: STORE, Body: Data{Value: value}}
	return SendDataIP(msg, ip)
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

func Validate(msg Message, res Message) bool {
	if (msg.RPC == res.RPC+"_REPLY") && (msg.Id == res.Body.RequestId) {
		return true
	} else {
		return false
	}
}

func GetUDPAddrFromContact(contact *Contact) net.UDPAddr {
	addr, port, _ := net.SplitHostPort(contact.Address)
	netAddr := net.ParseIP(addr)
	intPort, _ := strconv.Atoi(port)
	netAddress := net.UDPAddr{
		IP:   netAddr,
		Port: intPort,
	}
	return netAddress
}
