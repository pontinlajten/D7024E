package d7024e

import (
	"testing"
)

func TestNewNetwork(t *testing.T) {
	testIP := "162.20.0.0:1000"
	testID := NewKademliaID("62fa764de089aa2fcc265a3fe57991aa53af2a94")
	testKad := NewKademlia(testIP)
	testNet := CreateNetwork(&testKad)
	if !testNet.Kademlia.Id.Equals(testID) {
		t.Errorf("Excpected diffrent ID")
	}
}
