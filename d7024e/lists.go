package d7024e

import (
	"sort"
	"sync"
)

//lookup list
type Lookup struct { // List
	Cons  []Item
	Mutex sync.Mutex
}

type List struct { // Temp
	Cons  []Item
	Mutex sync.Mutex
}

type Item struct {
	Con  Contact
	Seen bool
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
	copyOfList := list.Cons //2
	responeList := List{}   //1

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

func (list *List) SortIt(list1 []Item, list2 []Item) Lookup {
	sorted := Lookup{}
	sorted.Append(list1)
	sorted.Append(list2)
	sorted.Sort()
	return sorted
}

// Append an array of Contacts to the ContactCandidates
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

// GetContacts returns the first count number of Contacts
func (candidates *Lookup) GetContacts(count int) []Item {
	return candidates.Cons[:count]
}

// Sort the Contacts in ContactCandidates
func (candidates *Lookup) Sort() {
	sort.Sort(candidates)
}

// Len returns the length of the ContactCandidates
func (candidates *Lookup) Len() int {
	return len(candidates.Cons)
}

func (candidates *List) Len() int {
	return len(candidates.Cons)
}

// Swap the position of the Contacts at i and j
// WARNING does not check if either i or j is within range
func (candidates *Lookup) Swap(i, j int) {
	candidates.Cons[i], candidates.Cons[j] = candidates.Cons[j], candidates.Cons[i]
}

// Less returns true if the Contact at index i is smaller than
// the Contact at index j
func (candidates *Lookup) Less(i, j int) bool {
	return candidates.Cons[i].Con.Less(&candidates.Cons[j].Con)
}
