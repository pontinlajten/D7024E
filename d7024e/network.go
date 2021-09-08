package d7024e

import (
	"fmt"
	"net"
)

type Network struct {
	me       *Contact
	table    *RoutingTable
	kademlia *Kademlia
}

const (
	CONN_TYPE = "udp"
)

// Template for init. an network.
func createNetwork(me *Contact, table *RoutingTable, kademlia *Kademlia) Network {
	network := Network{} // Create from Network struct
	network.me = me
	network.table = table
	network.kademlia = kademlia
	return network
}

// IN-PROGRESS
func (network *Network) Listen(me Contact, port int) {
	raddr, err := net.ResolveUDPAddr(CONN_TYPE, me.Address) // ResolveUDPAddr(str, str)
	conn, err2 := net.ListenUDP(CONN_TYPE, raddr)
	if (err != nil) || (err2 != nil) {
		fmt.Println("Error udp: ", err, "    ", err2)
	}

	defer conn.Close()

	ch := make()
	buffer := make([]byte, 1024) // Recieve ASCII, byte representation.

}

func (network *Network) SendPingMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindContactMessage(contact *Contact) {
	// TODO
}

func (network *Network) SendFindDataMessage(hash string) {
	// TODO
}

func (network *Network) SendStoreMessage(data []byte) {
	// TODO
}
