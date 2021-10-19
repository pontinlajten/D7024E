package d7024e

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type Cli struct {
	Network *Network
}

func NewCli(network *Network) *Cli {
	cli := &Cli{network}
	return cli
}

func (cli *Cli) Run() {
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

				contacts := cli.Network.Kademlia.Store(upload)
				fmt.Println("---------------------")
				fmt.Println()
				fmt.Print("Data stored at: ")
				fmt.Println(contacts)
				fmt.Println("---------------------")

			} else {
				fmt.Println("Invalid arguments for PUT...")
			}
		case "GET":
			if len(inputSplit) == 2 {
				find := inputSplit[1]
				b, _ := cli.Network.Kademlia.LookupData(find)

				fmt.Println("---------------------")
				fmt.Println("Value returned: " + string(b))
				fmt.Println("----------------------")
			} else {
				fmt.Println("Invalid arguments for GET...")
			}
		default:
			fmt.Println("INVALID COMMAND")
			fmt.Println("EXIT, PUT <arg1>, GET <arg1> <arg2> ...")

			fmt.Println("")
			fmt.Println(cli.Network.Kademlia.Rt.FindClosestContacts(cli.Network.Kademlia.Me.ID, bucketSize))
		}
	}
}
