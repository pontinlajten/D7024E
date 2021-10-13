package d7024e

import (
	"bufio"
	"fmt"
	"os"
	"regexp"
	"strings"
)

type cli struct {
	Network *Network
}

func NewCli(network *Network) *cli {
	cli := &cli{network}
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

			//cli.Network.Kademlia.Store(upload)


			//h_uploaded := cli.Network.Kademlia.HashIt(upload)

			//response := cli.Network.SendFindDataMessage(h_uploaded, &cli.Network.Kademlia.Me)
			/*
				if response {
					fmt.Println("Uploaded succesfully! Hashed: ")
					fmt.Println(h_uploaded)
				} else {

				}
			*/

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
	cli.Run()
}
