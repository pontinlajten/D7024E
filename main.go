package main

/* import{
	"fmt"
	kad "main/d7024e"
} */

import (
	"fmt"
	"log"
	"net"
)

const (
	port = "10000"
)

func main() {
	//fmt.Println("hello world")
	//d7024e.HashIt("192.158.1.38")
	//d7024e.NewKademlia("192.158.1.38")

	fmt.Println(GetOutboundIP())
	fmt.Println("123456")

	// Boostrap node ip = 172.20.0.2
	//boostrapIp = "172.20.0.2"

	// network := &kad.Network{}

	// go network.Listen()

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
