package d7024e

import (
	"testing"
)

type hash struct {
	IP   string
	Hash string
}

var hashes = []hash{
	{"162.20.0.0:1000", "62fa764de089aa2fcc265a3fe57991aa53af2a94"},
	{"162.20.0.1:1000", "b21d1e5991c149cbc6b651d47cfd8d4e7b76cf85"},
	{"xxx.xx.x.x:xxxx", "80d0a2db314dd4c28edf7f4c0202a28227384c66"},
}

func TestHashIt(t *testing.T) {
	for _, test := range hashes {
		testresult := HashIt(test.IP)
		if testresult != test.Hash {
			t.Errorf("wrong hash")
		}
	}
}

/*
func TestHashIt(t *testing.T) {
	testIP := "162.20.0.0:1000"
	testHash := "62fa764de089aa2fcc265a3fe57991aa53af2a94"
	hash := HashIt(testIP)

	if hash != testHash {
		t.Errorf("wrong hash")
	}

}
*/
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
