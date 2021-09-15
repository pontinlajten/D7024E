package d7024e

import (
	"fmt"
	"net"
	"sync"
)

type Network struct {
	me       *Contact
	mutex    *sync.Mutex
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
	network.mutex = &sync.Mutex{}
	network.table = table
	network.kademlia = kademlia
	return network
}

// IN-PROGRESS
func (network *Network) Listen(me Contact, port int) { // Listen(ip string, port int) original.
	raddr, err := net.ResolveUDPAddr(CONN_TYPE, ":8080") // ResolveUDPAddr(str, str). me.Address
	conn, err2 := net.ListenUDP(CONN_TYPE, raddr)
	if (err != nil) || (err2 != nil) {
		fmt.Println("Error udp: ", err, "    ", err2)
	}

	defer conn.Close() // defer: Close last, after all functions execution below is done.
	channel := make(chan []byte)
	buffer := make([]byte, MAX_BUFFER_SIZE) // Recieve ASCII, byte representation.

	for {
		n, addr, err := conn.ReadFromUDP(buffer)
		if err != nil {
			fmt.Println("Error ReadFromUDP", err)
		}
		fmt.Printf("packet-received: bytes=%d from=%s\n", n, addr.String())
	}
}

func (network *Network) SendPingMessage(contact *Contact) {

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
