package d7024e

import (
	"sync"
)

//lookup list
type Lookup struct {
	Cons  []Item
	Mutex sync.Mutex
}

type List struct {
	Cons  []Item
	Mutex sync.Mutex
}

type Item struct {
	Con  Contact
	Seen bool
}

// creates New list with x closest
func (kademlia *Kademlia) FindXClosest(target *Contact, x int) []Contact {
	xClosest := kademlia.Rt.FindClosestContacts(target.ID, x)
	return xClosest
}

func (kademlia *Kademlia) NewList(target *Contact) (list *List) {
	list = &List{}
	klist := kademlia.FindXClosest(target, bucketSize)

	for _, item := range klist {
		listitem := &Item{item, false}
		list.Cons = append(list.Cons, *listitem)
	}
	return
}

func (list *List) UpdateList(ID KademliaID, ch chan []Contact, net Network) {
	for {

		copyOfList := list.Cons //templist2
		responeList := List{}   //templist

		cons := <-ch

		for _, con := range cons {
			item := Item{con, false}
			responeList.Cons = append(responeList.Cons, item)
		}

		SortedList := list.SortIt(copyOfList, responeList.Cons)

		if len(SortedList.Cons) >= K {
			list.Cons = SortedList.GetContacts(K)
		} else {
			list.Cons = SortedList.GetContacts(len(SortedList.Cons))

		}
		//more to do
	}
}

func (list *List) SortIt(list1 []Item, list2 []Item) Lookup {
	sorted := Lookup{}
	sorted.Append(list1)
	sorted.Append(list2)
	sorted.Sort()
	return sorted
}

/*
func RecieverResponse(reciver Contact, nt Network, ch chan []Contact) {
	response, _ := nt.SendFindContactMessage(&reciver)
	ch <- response
}
*/
