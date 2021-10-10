package d7024e

import (
	"sort"
)

// Append an array of Contacts to the ContactCandidates
func (candidates *Lookup) Append(contacts []Item) {
	candidates.Cons = append(candidates.Cons, contacts...)
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
