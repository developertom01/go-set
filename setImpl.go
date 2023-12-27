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
		nextItem any
		set      setImpl
	}
)

func NewSet() *setImpl {
	return &setImpl{
		byKey:    make(map[any]int),
		byAccess: []any{},
	}
}

func (s *setImpl) addInternal(item any) {
	//Check if item is in set
	s.mut.Lock()
	defer s.mut.Unlock()
	_, ok := s.byKey[item]
	if ok {
		return
	}
	s.byAccess = append(s.byAccess, item)
	index := len(s.byAccess) - 1
	s.byKey[item] = index
}

func (s *setImpl) Add(item any) {
	s.addInternal(item)
}

func (s *setImpl) hasInternal(item any) bool {
	//Check if item is in set and the item is exactly what is in the array
	s.mut.Lock()
	defer s.mut.Unlock()
	index, ok := s.byKey[item]
	if !ok {
		return false
	}
	itemAccess := s.byAccess[index]
	return itemAccess == item
}

func (s *setImpl) Has(item any) bool {
	return s.hasInternal(item)
}

// Time complexity is O(len(set2))
func (s *setImpl) Contains(set *setImpl) bool {
	//Check if item is in set and the item is exactly what is in the array
	s.mut.Lock()
	set.mut.Lock()
	defer s.mut.Unlock()
	defer set.mut.Unlock()
	for _, item := range set.byAccess {
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
	return len(s.byAccess)
}

func (s *setImpl) IsEmpty() bool {
	return len(s.byAccess) == 0
}

func (s *setImpl) ToSlice() []any {
	slice := make([]any, 0, s.Len())
	s.mut.Lock()
	defer s.mut.Unlock()
	copy(slice, s.byAccess)
	return slice
}

func (s *setImpl) Union(set *setImpl) *setImpl {
	//Build new set
	var (
		//Build new set
		newSet = NewSet()
	)

	var wg sync.WaitGroup

	wg.Add(2)

	go func() {
		defer wg.Done()
		for _, item := range s.ToSlice() {
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

func (s *setImpl) Intersection(set *setImpl) *setImpl {
	var (
		results  = NewSet()
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

func (s *setImpl) Complement(set *setImpl) *setImpl {
	var (
		results  = NewSet()
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
