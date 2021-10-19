package d7024e

import (
	"log"
	"net"
	"testing"
)

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

func TestNewNetwork(t *testing.T) {
	testIP := "162.20.0.0:1000"
	testID := NewKademliaID("62fa764de089aa2fcc265a3fe57991aa53af2a94")
	testKad := NewKademlia(testIP)
	testNet := CreateNetwork(&testKad)
	if !testNet.Kademlia.Id.Equals(testID) {
		t.Errorf("Excpected diffrent ID")
	}
}

func TestSendPing(t *testing.T) {
	nodeIp := GetOutboundIP()
	port := "1000"
	localIP := nodeIp.String() + ":" + port

	network := Network{}
	kademlia := NewKademlia(localIP)
	network.Kademlia = &kademlia
	go network.Listen()

	err := network.SendPingMessage(&network.Kademlia.Me)
	if err != nil {
		t.Fail()
	}
}

func TestSendStore(t *testing.T) {
	nodeIp := GetOutboundIP()
	port := "1001"
	localIP := nodeIp.String() + ":" + port

	network := Network{}
	kademlia := NewKademlia(localIP)
	network.Kademlia = &kademlia
	go network.Listen()

	value := "1337"
	response, _ := network.SendStoreMessage(value, &network.Kademlia.Me)

	if response.Body.Key == "" {
		t.Fail()
	}
}

func TestSendFindData(t *testing.T) {
	nodeIp := GetOutboundIP()
	port := "1002"
	localIP := nodeIp.String() + ":" + port

	network := Network{}
	kademlia := NewKademlia(localIP)
	network.Kademlia = &kademlia
	go network.Listen()

	value := "1337"
	storeResponse, _ := network.SendStoreMessage(value, &network.Kademlia.Me)
	hash := storeResponse.Body.Key
	response, _ := network.SendFindDataMessage(hash, &network.Kademlia.Me)

	if response.Body.Value != value {
		t.Fail()
	}
}
