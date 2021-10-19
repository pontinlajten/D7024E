package d7024e

import (
	"testing"
)

func TestNewList(t *testing.T) {
	testKad := NewKademlia("162.20.0.0")

	testKad.Rt.AddContact(NewContact(NewKademliaID("62fa764de089aa2fcc265a3fe57991aa53af2a94"), "1000"))

	for i := 1; i < K; i++ {
		testKad.Rt.AddContact(NewContact(NewRandomKademliaID(), "1000"))

	}
	newList := testKad.NewList(NewKademliaID("cef2cf821221e92087544fc59a023871797051fb"))
	testLen := newList.Len()

	if K != testLen {
		t.Errorf("Wrong length on list")
	}
}
