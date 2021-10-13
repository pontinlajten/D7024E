package main

/* import{
	"fmt"
	kad "main/d7024e"
} */

import (
	"fmt"
	"log"
	"net"

	kad "main/d7024e"
)

const (
	port        = "1000"
	bootstrapIp = "172.20.0.2"
)

func main() {

	//fmt.Println("hello world")
	//d7024e.HashIt("192.158.1.38")
	//kad.NewKademlia("192.158.1.38")

	//fmt.Println(GetOutboundIP())
	//fmt.Println("123456")

	nodeIp := GetOutboundIP()

	bsIP := bootstrapIp + ":" + port

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

	cli := kad.NewCli(&network)
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
