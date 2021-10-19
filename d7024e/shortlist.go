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

func (kademlia *Kademlia) NewList(targetID *KademliaID) (list *List) {
	list = &List{}

	closestK := kademlia.Rt.FindClosestContacts(targetID, bucketSize)

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

func (lookuplist *List) updateLookupData(hash string, ch chan []Contact, target chan []byte, dataContactCh chan Contact, net Network, wg sync.WaitGroup) ([]byte, Contact) {
	for {
		contacts := <-ch
		targetData := <-target
		dataContact := <-dataContactCh

		if targetData != nil {
			return targetData, dataContact
		}

		nextContact, Done := lookuplist.Update(contacts)
		if Done {
			return nil, Contact{}
		} else {
			go AsyncLookupData(hash, nextContact, net, ch, target, dataContactCh)
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
func (candidates *Lookup) Append(contacts []Item) {
	for _, nextCandidate := range contacts {
		approved := true

		for _, candidate := range candidates.Cons {

			if candidate.Con.ID.Equals(nextCandidate.Con.ID) {
				approved = false
				break
			}
		}
		if approved {
			candidates.Cons = append(candidates.Cons, nextCandidate)
		}
	}
}

/*
	Everything below is gathered from contact.go
*/

func (candidates *Lookup) GetContacts(count int) []Item {
	return candidates.Cons[:count]
}

func (candidates *Lookup) Sort() {
	sort.Sort(candidates)
}

func (candidates *Lookup) Len() int {
	return len(candidates.Cons)
}

func (candidates *List) Len() int {
	return len(candidates.Cons)
}

func (candidates *Lookup) Swap(i, j int) {
	candidates.Cons[i], candidates.Cons[j] = candidates.Cons[j], candidates.Cons[i]
}

func (candidates *Lookup) Less(i, j int) bool {
	return candidates.Cons[i].Con.Less(&candidates.Cons[j].Con)
}
