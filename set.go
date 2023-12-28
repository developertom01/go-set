package set

type (
	//Concurrent implementation of set
	Set interface {
		//Add new element to the set
		Add(any)

		//Checks if an element is in set
		Has(any) bool

		//Checks  if all members in another set is in parent  set
		Contains(Set) bool

		//Removes element from a set
		Remove(any)

		//Length of set
		Len() int

		//Checks if set is empty
		IsEmpty() bool

		// -- Operations

		//Returns the union of the set with another set
		Union(Set) Set

		//Return new set of intersection with another set
		Intersection(Set) Set

		//New set of element in parent set but not in given set
		Complement(Set) Set

		// Converts set into array
		ToSlice() []any

		Iterator() Iterator
	}

	Iterator interface {
		HasNext() bool
		Next() any
		Close()
	}
)
