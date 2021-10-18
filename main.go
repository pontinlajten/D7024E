package main

/* import{
	"fmt"
	kad "main/d7024e"
} */

import (
	"fmt"

	project "main/d7024e"
)

const (
	port        = "1000"
	bs_template = "xxx.x.xx.x"
)

func main() {
	ip := "162.20.0.0:1000"
	hash := project.HashIt(ip)
	fmt.Println(hash)
	/*
		nodeIp := GetOutboundIP()

		bsIP := GenerateBootstrap(nodeIp.String(), bs_template) + ":" + port

		localIP := nodeIp.String() + ":" + port

		fmt.Println("BootStrap ip:", bsIP)
		fmt.Println("My ip:", localIP)

		bsID := project.NewKademliaID(project.HashIt(bsIP))
		bsContact := project.NewContact(bsID, bsIP)

		me := project.NewKademlia(localIP)

		network := project.CreateNetwork(&me)

		if localIP != bsIP {
			// newContact := kad.NewContact(kad.NewKademliaID(kad.HashIt(bsIP)), bsIP)
			me.InitNetwork(&bsContact)
			//fmt.Printf("\nRoutingtable: %x\n", me.Rt.FindClosestContacts(me.Me.ID, 4))
		}

		go network.Listen()

		cli := project.NewCli(&network)
		cli.Run()
	*/
}

/*
func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func GenerateBootstrap(str string, bp string) string {
	dif := len(str) - len(bp)

	for len(str) > 0 {
		_, size := utf8.DecodeLastRuneInString(str)
		return str[:len(str)-size-dif] + "2"
	}

	return str
}
*/
