package d7024e

import (
	"encoding/json"
	"fmt"
	"net"
	"strconv"
	"sync"

	"github.com/pkg/errors"
)

type Network struct {
	Mutex    *sync.Mutex
	Kademlia *Kademlia // Used in Listen.
}

const (
	CONN_TYPE       = "udp"
	MAX_BUFFER_SIZE = 4096
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
func (network *Network) Listen() { // Listen(ip string, port int) original.
	server := GetUDPAddrFromContact(&network.Kademlia.Me)
	conn, _ := net.ListenUDP(CONN_TYPE, &server)
	defer conn.Close() // defer: Close last, after all functions execution below is done.

	//channel := make(chan []byte)
	buffer := make([]byte, MAX_BUFFER_SIZE) // Recieve ASCII, byte representation.

	for {
		n, addr, _ := conn.ReadFromUDP(buffer)

		decoded := unmarshall(buffer[:n])

		network.Kademlia.Rt.AddContact(*decoded.Source)

		msg := network.MsgHandler(decoded)
		replyEncoded := marshall(msg)

		sendResponse(replyEncoded, addr, conn)
	}
}

func (network *Network) MsgHandler(decoded Message) Message {
	reply := Message{}
	body := MsgBody{}

	// HANDLING REQUEST->NODE
	if decoded.RPC == FIND_NODE {
		reply.Source = &network.Kademlia.Me
		reply.RPC = FIND_NODE_REPLY
		body.Nodes = network.FindNodeHandler(decoded)

		body.OriginalSource = decoded.Source
		reply.Body = body
	} else if decoded.RPC == PING {
		reply.Source = &network.Kademlia.Me
		reply.RPC = PING_REPLY

		body.OriginalSource = decoded.Source
		reply.Body = body
	} else if decoded.RPC == FIND_DATA {
		reply = network.FindValueHandler(decoded)

		body.OriginalSource = decoded.Source
		reply.Body = body
	} else if decoded.RPC == STORE {
		reply.Source = &network.Kademlia.Me
		reply.RPC = STORE_REPLY

		body.Key = network.StoreHandler(decoded)

		body.OriginalSource = decoded.Source
		reply.Body = body
	}

	return reply
}

func (network *Network) FindNodeHandler(msg Message) []Contact {
	contacts := network.Kademlia.Rt.FindClosestContacts(msg.Body.TargetId, bucketSize)
	return contacts
}

func (network *Network) FindValueHandler(msg Message) Message {
	keyVal := network.Kademlia.LookupDataHash(network.Kademlia.HashIt(msg.Body.Key))

	if keyVal != nil {
		return Message{Source: &network.Kademlia.Me, RPC: FIND_DATA_REPLY, Body: MsgBody{Key: keyVal.Key, Value: keyVal.Value}}
	} else {
		id := NewKademliaID(network.Kademlia.HashIt(msg.Body.Key))
		newContacts := network.Kademlia.Rt.FindClosestContacts(id, bucketSize)

		return Message{Source: &network.Kademlia.Me, RPC: FIND_DATA_REPLY, Body: MsgBody{Nodes: newContacts}}
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
		return Message{}, errors.Wrap(err, "Client: Failed to open connection to "+address.IP.String())
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
	Client.Write(sendMsg)

	buf := make([]byte, MAX_BUFFER_SIZE)

	n, _, _ := Client.ReadFrom([]byte(buf))
	response := unmarshall(buf[:n])

	if err != nil {
		fmt.Printf("failed to %s to %s error: %s", rpcMsg, address.String(), err)
	}

	/*
		if network.Validate(msg, response) && response.RPC == FIND_NODE_REPLY {
			for item := range response.Body.Nodes {
				network.Kademlia.Rt.AddContact(response.Body.Nodes[item])
			}
		}
	*/

	if network.Validate(msg, response) {
		network.Kademlia.Rt.AddContact(*response.Source) // Updates routing table if recieves succesful respond from target node.
	}

	return response, nil
}

func (network *Network) SendPingMessage(contact *Contact) error {
	msg := Message{Source: &network.Kademlia.Me, RPC: PING}
	_, err := network.SendData(msg, contact)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}
	return nil
}

func (network *Network) SendFindContactMessage(contact *Contact, targetId *KademliaID) ([]Contact, error) {
	msg := Message{Source: &network.Kademlia.Me, RPC: FIND_NODE, Body: MsgBody{TargetId: targetId}}
	res, err := network.SendData(msg, contact)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}

	return res.Body.Nodes, nil
}

func (network *Network) SendFindDataMessage(hash string, contact *Contact) (Message, error) {
	msg := Message{Source: &network.Kademlia.Me, RPC: FIND_DATA, Body: MsgBody{Key: hash}}
	return network.SendData(msg, contact)
}

func (network *Network) SendStoreMessage(value string, contact *Contact) (Message, error) {
	msg := Message{Source: &network.Kademlia.Me, RPC: STORE, Body: MsgBody{Value: value}}
	res, err := network.SendData(msg, contact)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}
	return res, nil
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

func (network *Network) Validate(msg Message, res Message) bool {
	if (msg.RPC+"_REPLY" == res.RPC) && (network.Kademlia.Me.ID.String() == res.Body.OriginalSource.ID.String()) { // Check if message is original sender.
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
