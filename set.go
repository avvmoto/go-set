// Copyright 2019 avvmoto. All rights reserved.

// package set implements list which can delete item, specially intended to be able to get all items fast.
package set

import "fmt"

// ItemIterator allows callers of All() to iterate items.
// When this function returns false, iteration will stop and
// the associated All() function will immediately return.
type Iterator func(item Item) bool

type Interface interface {
	// All iterate all items in the set
	All(Iterator)

	// Delete delete item from the set
	Delete(item Item) Item

	// Append append item to the set
	Append(item Item)
}

// Set implements list which can delete item, specially intended to be able to get all items fast.
// All is faster, but Delete and Append may slower than Set who simply use map as internal data structure.
// See SetMap in set_test.go
type Set struct {
	indexOf map[int]int
	items   []Item
	deleted []bool
	len     int
}

// NewSet create Set with given capacity c.
func NewSet(c int) *Set {
	return &Set{
		indexOf: make(map[int]int, c),
		items:   make([]Item, 0, c),
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
func (s *Set) Delete(item Item) Item {
	idx, ok := s.indexOf[item.Key()]
	if !ok {
		return nil
	}

	s.deleted[idx] = true
	s.len--

	return s.items[idx]
}

// Append append item to the set.
func (s *Set) Append(item Item) {
	_, ok := s.indexOf[item.Key()]
	if ok {
		panic(fmt.Sprintf("duplicate item:%#v", item))
	}

	s.indexOf[item.Key()] = len(s.items)

	s.items = append(s.items, item)
	s.deleted = append(s.deleted, false)
	s.len++
}

// Clear removes all items from the set.
func (s *Set) Clear() {
	s.items = s.items[:0]
	s.deleted = s.deleted[:0]
	s.len = 0
}

// Len returns the number of items currently in the set.
func (s *Set) Len() int {
	return s.len
}

// Item represents a single object in the set.
type Item interface {

	// Key represents Item uniqueness. Key must be able to be used as map key.
	Key() int
}

// Int implements the Item interface for integers.
type Int int

// Key returns key for map.
func (i Int) Key() int {
	return int(i)
}
