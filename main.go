package main

/* import{
	"fmt"
	kad "main/d7024e"
} */

import (
	"fmt"
	"log"
	"net"
	"unicode/utf8"

	kad "main/d7024e"
)

const (
	port = "1000"
)

func main() {

	nodeIp := GetOutboundIP()

	bsIP := GenerateBootstrap(nodeIp.String()) + ":" + port

	localIP := nodeIp.String() + ":" + port

	fmt.Println("BootStrap ip:", bsIP)
	fmt.Println("New ip:", localIP)

	bsID := kad.NewKademliaID(kad.HashIt(bsIP))
	bsContact := kad.NewContact(bsID, bsIP)

	me := kad.NewKademlia(localIP)
	me.InitRt(&bsContact)

	network := kad.CreateNetwork(&me)

	if localIP != bsIP {
		newContact := kad.NewContact(kad.NewKademliaID(kad.HashIt(bsIP)), bsIP)
		me.InitRt(&newContact)
		fmt.Printf("\nRoutingtable: %x\n", me.Rt.FindClosestContacts(me.Me.ID, 4))
	}

	go network.Listen(port)

	cli := kad.NewCli(&network, bsIP)
	cli.Run()
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GenerateBootstrap(str string) string {
	for len(str) > 0 {
		_, size := utf8.DecodeLastRuneInString(str)
		return str[:len(str)-size] + "2"
	}

	return str
}
