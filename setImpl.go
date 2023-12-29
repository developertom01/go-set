package set

import (
	"sync"
)

type (
	setImpl struct {
		mut      sync.RWMutex
		byKey    map[any]int
		byAccess []any
	}
	iteratorImpl struct {
		nextItem *any
		set      *setImpl
	}
)

func newSet() *setImpl {
	return &setImpl{
		byKey:    make(map[any]int),
		byAccess: []any{},
	}
}

func NewSet() Set {
	set := newSet()
	return set
}

func (s *setImpl) addInternal(item any) {
	_, ok := s.byKey[item]
	if ok {
		return
	}
	s.byAccess = append(s.byAccess, item)
	index := len(s.byAccess) - 1
	s.byKey[item] = index
}

func (s *setImpl) Add(item any) {
	s.mut.Lock()
	defer s.mut.Unlock()

	s.addInternal(item)
}

func (s *setImpl) hasInternal(item any) bool {
	//Check if item is in set and the item is exactly what is in the array
	index, ok := s.byKey[item]
	if !ok {
		return false
	}
	itemAccess := s.byAccess[index]

	return itemAccess == item
}

func (s *setImpl) Has(item any) bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	return s.hasInternal(item)
}

// Time complexity is O(len(set2))
func (s *setImpl) Contains(set Set) bool {
	s.mut.Lock()
	defer s.mut.Unlock()

	if set.Len() > s.lenInternal() {
		return false
	}

	for _, item := range set.ToSlice() {
		if exists := s.hasInternal(item); !exists {
			return false
		}
	}

	return true
}

func (s *setImpl) Remove(item any) {
	s.mut.Lock()
	defer s.mut.Unlock()

	index, ok := s.byKey[item]
	if !ok {
		return
	}

	delete(s.byKey, item)
	s.byAccess = append(s.byAccess[:index], s.byAccess[index+1:]...)
}

func (s *setImpl) Len() int {
	s.mut.Lock()
	defer s.mut.Unlock()
	length := s.lenInternal()

	return length
}
func (s *setImpl) lenInternal() int {
	return len(s.byAccess)
}

func (s *setImpl) IsEmpty() bool {
	return len(s.byAccess) == 0
}

func (s *setImpl) toSliceInternal() []any {
	slice := make([]any, s.lenInternal())
	copy(slice, s.byAccess)

	return slice
}

func (s *setImpl) ToSlice() []any {
	s.mut.Lock()
	defer s.mut.Unlock()

	slice := s.toSliceInternal()

	return slice
}

func (s *setImpl) Union(set Set) Set {
	//Build new set
	var (
		//Build new set
		newSet = newSet()
		wg     sync.WaitGroup
	)

	wg.Add(2)

	go func() {
		defer wg.Done()

		for _, item := range s.toSliceInternal() {
			newSet.addInternal(item)
		}
	}()

	go func() {
		defer wg.Done()

		for _, item := range set.ToSlice() {
			newSet.addInternal(item)
		}
	}()

	wg.Wait()

	return newSet
}

func (s *setImpl) Intersection(set Set) Set {
	var (
		results  = newSet()
		setItems = set.ToSlice()
	)

	s.mut.Lock()
	defer s.mut.Unlock()

	for _, item := range setItems {
		if exists := s.hasInternal(item); exists {
			results.addInternal(item)
		}
	}

	return results
}

func (s *setImpl) Complement(set Set) Set {
	var (
		results  = newSet()
		setItems = set.ToSlice()
	)

	s.mut.Lock()
	defer s.mut.Unlock()

	for _, item := range setItems {
		if exists := !s.hasInternal(item); exists {
			results.addInternal(item)
		}
	}

	return results
}

func (s *setImpl) Iterator() Iterator {
	itr := newIterator(s)
	s.mut.Lock()
	itr.prepareNext()

	return itr
}

func newIterator(set *setImpl) *iteratorImpl {
	return &iteratorImpl{
		set: set,
	}
}

func (itr *iteratorImpl) prepareNext() {
	if itr.nextItem == nil && !itr.set.IsEmpty() {
		itr.nextItem = &itr.set.byAccess[0]
	} else {
		prevItem := itr.nextItem
		index := itr.set.byKey[*prevItem]
		if index+1 >= itr.set.lenInternal() {
			itr.nextItem = nil
		} else {
			itr.nextItem = &itr.set.byAccess[index+1]
		}
	}
}

func (itr *iteratorImpl) HasNext() bool {
	return itr.nextItem != nil
}

func (itr *iteratorImpl) Next() any {
	if itr.nextItem == nil {
		panic("Set out of range")
	}

	next := itr.nextItem
	itr.prepareNext()

	return next
}

func (itr *iteratorImpl) Close() {
	itr.set.mut.Unlock()
}
