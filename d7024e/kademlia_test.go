package d7024e

import (
	"testing"
)

func TestHashIt(t *testing.T) {
	testIP := "162.20.0.0:1000"
	testHash := "62fa764de089aa2fcc265a3fe57991aa53af2a94"
	hash := HashIt(testIP)

	if hash != testHash {
		t.Errorf("wrong hash")
	}

}
func TestNewKademlia(t *testing.T) {
	testIP := "162.20.0.0:1000"
	testHash := "62fa764de089aa2fcc265a3fe57991aa53af2a94"
	testID := NewKademliaID(testHash)
	newKad := NewKademlia(testIP)
	newContact := newKad.Id

	if !newContact.Equals(testID) {
		t.Errorf("Wrong id")
	}

}
