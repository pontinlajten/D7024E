package d7024e

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"net"
	"strconv"
)

/*
	Used when transmitting data to nodes. (ESSENTIAL)
*/
func marshall(msg Message) []byte {
	encoded, _ := json.Marshal(msg)
	return encoded
}

/*
	Used when transmitting data to nodes. (ESSENTIAL)
*/
func unmarshall(data []byte) Message {
	var decoded Message
	json.Unmarshal([]byte(data), &decoded)
	return decoded
}

/*
	Validator in network inorder to check if response is legit.
*/
func (network *Network) Validate(msg Message, res Message) bool {
	if (msg.RPC+"_REPLY" == res.RPC) && (network.Kademlia.Me.ID.String() == res.Body.OriginalSource.ID.String()) { // Check if message is original sender. And correct RPC.
		return true
	} else {
		return false
	}
}

/*
	Great converter inorder to get correct representation of IP:PORT from contact!
*/
func GetUDPAddrFromContact(contact *Contact) net.UDPAddr {
	addr, port, _ := net.SplitHostPort(contact.Address)
	netAddr := net.ParseIP(addr)
	intPort, _ := strconv.Atoi(port)
	netAddress := net.UDPAddr{
		IP:   netAddr,
		Port: intPort,
	}
	return netAddress
}

/*
	Used when storing hashed data in kademlia and so on.
*/
func (kademlia *Kademlia) HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))

	return hash
}

func HashIt(str string) string {
	hashStr := sha1.New()
	hashStr.Write([]byte(str))
	hash := hex.EncodeToString(hashStr.Sum(nil))

	return hash
}
