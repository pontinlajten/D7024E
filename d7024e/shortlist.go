package d7024e

import (
	"sort"
	"sync"
)

//lookup list
type Lookup struct { // List for shorlist
	Cons []Item
}

type List struct { // Temp for shortlist
	Cons []Item
}

type Item struct {
	Con  Contact
	Seen bool // IF VISITED.
}

//return a list of the k closest kademlianodes from kademlias routingtable
func (kademlia *Kademlia) NewList(targetID *KademliaID) (list *List) {

	closestK := kademlia.Rt.FindClosestContacts(targetID, bucketSize)

	list = &List{}

	for _, item := range closestK {
		listitem := &Item{item, false}
		list.Cons = append(list.Cons, *listitem)
	}
	return
}

func (list *List) Update(cons []Contact) (Contact, bool) {
	copyOfList := list.Cons
	responeList := List{}

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
	nextContact, Finished := list.findContact()
	return nextContact, Finished
}

func (list *List) findContact() (Contact, bool) {
	var newContact Contact
	Finished := true
	for i, item := range list.Cons {
		if !item.Seen {
			list.Cons[i].Seen = true
			Finished = false
		}
	}
	return newContact, Finished
}

func (list *List) UpdateList(ID KademliaID, ch chan []Contact, net Network) {
	for {
		contacts := <-ch
		nextContact, Finished := list.Update(contacts)

		if Finished {
			return
		} else {
			go AsyncFindContact(nextContact, ID, net, ch)
		}
	}
}

func (list *List) updateFindData(hash string, ch chan []Contact, target chan []byte, dataConCh chan Contact, net Network, wg sync.WaitGroup) ([]byte, Contact) {
	for {
		cons := <-ch
		targetData := <-target
		dataCon := <-dataConCh

		if targetData != nil {
			return targetData, dataCon
		}

		nextContact, Done := list.Update(cons)
		if Done {
			return nil, Contact{}
		} else {
			go AsyncFindData(hash, nextContact, net, ch, target, dataConCh)
		}
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
	Modified version of contact.go append. Instead applied to the shortlist.
*/
func (lookup *Lookup) Append(contacts []Item) {
	for _, nextCon := range contacts {
		approved := true

		for _, con := range lookup.Cons {

			if con.Con.ID.Equals(nextCon.Con.ID) {
				approved = false
				break
			}
		}
		if approved {
			lookup.Cons = append(lookup.Cons, nextCon)
		}
	}
}

/*
	Everything below is gathered from contact.go
*/

func (lookup *Lookup) GetContacts(count int) []Item {
	return lookup.Cons[:count]
}

func (lookup *Lookup) Sort() {
	sort.Sort(lookup)
}

func (lookup *Lookup) Len() int {
	return len(lookup.Cons)
}

func (list *List) Len() int {
	return len(list.Cons)
}

func (lookup *Lookup) Swap(i, j int) {
	lookup.Cons[i], lookup.Cons[j] = lookup.Cons[j], lookup.Cons[i]
}

func (lookup *Lookup) Less(i, j int) bool {
	return lookup.Cons[i].Con.Less(&lookup.Cons[j].Con)
}
