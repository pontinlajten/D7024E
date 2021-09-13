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
	CONN_TYPE       = "udp"
	MAX_BUFFER_SIZE = 1024
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

	defer conn.Close()                      // defer: Close last, after all functions execution below is done.
	buffer := make([]byte, MAX_BUFFER_SIZE) // Recieve ASCII, byte representation.

	for {
		n, addr, err := conn.ReadFrom(buffer)
		if err != nil {
			fmt.Println("Error ReadFromUDP", err)
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())

		// n, err = pc.WriteTo(buffer[:n], addr)
	}

	// buffer <- conn

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
