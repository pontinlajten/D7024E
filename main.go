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
	port        = "10000"
	bootstrapIp = "172.20.0.2"
)

func main() {

	//fmt.Println("hello world")
	//d7024e.HashIt("192.158.1.38")
	//kad.NewKademlia("192.158.1.38")

	//fmt.Println(GetOutboundIP())
	//fmt.Println("123456")

	nodeIp := GetOutboundIP()

	bsNode := bootstrapIp + ":" + port

	node := nodeIp.String() + ":" + port

	fmt.Println("BootStrap ip:", bsNode)
	fmt.Println("New ip:", node)

	// network := &kad.Network{}

	// go network.Listen()
	newNode := kad.NewKademlia(node)
	fmt.Println(newNode)
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
