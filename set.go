// Copyright 2019 avvmoto. All rights reserved.

// package set implements list which can delete item, specially intended to be able to get all items fast.
package set

// ItemIterator allows callers of All() to iterate items.
// When this function returns false, iteration will stop and
// the associated All() function will immediately return.
type Iterator func(item interface{}) bool

type Interface interface {
	// All iterate all items in the set
	All(Iterator)

	// Delete delete item from the set
	Delete(item interface{})

	// Append append item to the set
	Append(item interface{})
}

// Set implements list which can delete item, specially intended to be able to get all items fast.
// All is faster, but Delete and Append may slower than Set who simply use map as internal data structure.
// See SetMap in set_test.go
type Set struct {
	indexOf map[interface{}]int
	items   []interface{}
	deleted []bool
}

// NewSet create Set with given capacity c.
func NewSet(c int) *Set {
	return &Set{
		indexOf: make(map[interface{}]int, c),
		items:   make([]interface{}, 0, c),
		deleted: make([]bool, 0, c),
	}

}

// All iterate all items in the set.
func (s *Set) All(fn Iterator) {
	for i, item := range s.items {
		if s.deleted[i] {
			continue
		}

		if !fn(item) {
			break
		}
	}
}

// Delete delete item from the set.
func (s *Set) Delete(item interface{}) {
	s.deleted[s.indexOf[item]] = true
}

// Append append item to the set.
func (s *Set) Append(item interface{}) {
	_, ok := s.indexOf[item]
	if ok {
		panic("duplicate item")
	}

	s.indexOf[item] = len(s.items)

	s.items = append(s.items, item)
	s.deleted = append(s.deleted, false)
}
