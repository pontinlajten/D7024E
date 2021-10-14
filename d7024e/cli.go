package d7024e

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type cli struct {
	Network  *Network
	TargetIP string
}

func NewCli(network *Network, ip string) *cli {
	cli := &cli{network, ip}
	return cli
}

func (cli *cli) Run() {
	fmt.Println("<CMD> ")
	input := bufio.NewScanner(os.Stdin)
	input.Scan()
	inputText := input.Text()

	space := regexp.MustCompile(` `)
	inputSplit := space.Split(inputText, 10)

	switch strings.ToUpper(inputSplit[0]) {
	case "EXIT":
		fmt.Println("EXIT ENTERED.")
		return
	case "PUT":
		if len(inputSplit) == 2 {
			upload := inputSplit[1]
			fmt.Println(upload)

			// resp := cli.Network.Kademlia.StoreIP(upload, cli.TargetIP)
			// fmt.Printf("Successfully uploaded value! Hash: %v", resp)
		} else {
			fmt.Println("Invalid arguments for PUT...")
		}
	case "GET":
		if len(inputSplit) > 2 {

		} else {
			fmt.Println("Invalid arguments for GET...")
		}
	default:
		fmt.Println("INVALID COMMAND")
		fmt.Println("EXIT, PUT <arg1>, GET <arg1> <arg2> ...")
	}

	fmt.Println("")
	fmt.Println(cli.Network.Kademlia.Rt.FindClosestContacts(cli.Network.Kademlia.Me.ID, 4))
	cli.Run()
}
