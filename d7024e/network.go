package d7024e

import (
	"fmt"
	"net"

	"github.com/pkg/errors"
)

type Network struct {
	Kademlia *Kademlia
}

const (
	CONN_TYPE       = "udp"
	MAX_BUFFER_SIZE = 8000
)

// Template for init. an network.
func CreateNetwork(kademlia *Kademlia) Network {
	network := Network{} // Create from Network struct
	network.Kademlia = kademlia
	return network
}

/////////////////////////////// RESPONSE /////////////////////////////////////////

func (network *Network) Listen() { // Listen(ip string, port int) original.
	server := GetUDPAddrFromContact(&network.Kademlia.Me) // Help function
	conn, _ := net.ListenUDP(CONN_TYPE, &server)
	defer conn.Close() // defer: Close last, after all functions execution below is done.

	buffer := make([]byte, MAX_BUFFER_SIZE) // Recieve ASCII, byte representation.

	for {
		n, addr, _ := conn.ReadFromUDP(buffer)

		decoded := unmarshall(buffer[0:n])

		network.Kademlia.Rt.AddContact(*decoded.Source)

		msg := network.ListenHandler(decoded)
		replyEncoded := marshall(msg)

		sendResponse(replyEncoded, addr, conn)
	}
}

func (network *Network) ListenHandler(decoded Message) Message {
	reply := Message{}

	reply.Source = &network.Kademlia.Me

	if decoded.RPC == FIND_NODE {
		reply.RPC = FIND_NODE_REPLY
		reply.Body.Nodes = network.FindNodeHandler(decoded)

	} else if decoded.RPC == PING {
		reply.RPC = PING_REPLY

	} else if decoded.RPC == FIND_DATA {
		reply.Body = network.FindValueHandler(decoded)

	} else if decoded.RPC == STORE {
		reply.RPC = STORE_REPLY
		reply.Body.Key = network.StoreHandler(decoded)

	}

	reply.Body.OriginalSource = decoded.Source

	return reply
}

func (network *Network) FindNodeHandler(msg Message) []Contact {
	contacts := network.Kademlia.Rt.FindClosestContacts(msg.Body.TargetId, 20)
	return contacts
}

func (network *Network) FindValueHandler(msg Message) MsgBody {
	keyVal := network.Kademlia.LookupDataHash(msg.Body.Key)

	id := NewKademliaID(msg.Body.Key)

	if keyVal != nil {
		return MsgBody{Key: keyVal.Key, Value: keyVal.Value}
	} else {
		newContacts := network.Kademlia.Rt.FindClosestContacts(id, 20)
		return MsgBody{Nodes: newContacts}
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

//////////////////////////////// REQUEST /////////////////////////////////////

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

	// Wait for respond from target node.

	buf := make([]byte, MAX_BUFFER_SIZE)

	n, _, _ := Client.ReadFrom([]byte(buf))
	response := unmarshall(buf[0:n])

	if err != nil {
		fmt.Printf("failed to %s to %s error: %s", rpcMsg, address.String(), err)
	}

	if network.Validate(msg, response) {
		network.Kademlia.Rt.AddContact(*response.Source) // Updates routing table if recieves succesful respond from target node.
	}

	return response, nil
}

func (network *Network) PingMessage(contact *Contact) error {
	msg := Message{Source: &network.Kademlia.Me, RPC: PING}
	_, err := network.SendData(msg, contact)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}
	return nil
}

func (network *Network) FindContactMessage(contact *Contact, targetId *KademliaID) ([]Contact, error) {
	msg := Message{Source: &network.Kademlia.Me, RPC: FIND_NODE, Body: MsgBody{TargetId: targetId}}
	res, err := network.SendData(msg, contact)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}

	return res.Body.Nodes, nil
}

func (network *Network) FindDataMessage(hash string, contact *Contact) (Message, error) {
	msg := Message{Source: &network.Kademlia.Me, RPC: FIND_DATA, Body: MsgBody{Key: hash}}
	return network.SendData(msg, contact)
}

func (network *Network) StoreMessage(value string, contact *Contact) (Message, error) {
	msg := Message{Source: &network.Kademlia.Me, RPC: STORE, Body: MsgBody{Value: value}}
	res, err := network.SendData(msg, contact)
	if err != nil {
		errors.Wrap(err, "Something went wrong")
	}
	return res, nil
}
