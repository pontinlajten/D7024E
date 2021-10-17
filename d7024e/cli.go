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
			if len(inputSplit) == 2 {
				upload := inputSplit[1]
				fmt.Println(upload)

				cli.Network.Kademlia.Store(upload)

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
		fmt.Println(cli.Network.Kademlia.Rt.FindClosestContacts(cli.Network.Kademlia.Me.ID, b_size))
	}
}
