package d7024e

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

const (
	b_size = 20
)

type cli struct {
	Network *Network
}

func NewCli(network *Network) *cli {
	cli := &cli{network}
	return cli
}

func (cli *cli) Run() {
	for {
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
			if inputSplit[0] == "PUT" {
				upload := inputSplit[1]

				hash := cli.Network.Kademlia.Store(upload)
				fmt.Printf("Successfully uploaded value! Hash: %v", hash)


			} else {
				fmt.Println("Invalid arguments for PUT...")
			}
		case "GET":
			if inputSplit[0] == "GET" {
				hash := inputSplit[1]
				fmt.Println(hash)

				value, nodeId := cli.Network.Kademlia.LookupData(hash)
				fmt.Printf("Succefully return value: %v from node: %v", value, nodeId)
			} else {
				fmt.Println("Invalid arguments for GET...")
			}
		default:
			fmt.Println("INVALID COMMAND")
			fmt.Println("EXIT, PUT <arg1>, GET <arg1> <arg2> ...")
		}

		fmt.Println("")
		fmt.Println(cli.Network.Kademlia.Rt.FindClosestContacts(cli.Network.Kademlia.Me.ID, b_size))
	}
}
