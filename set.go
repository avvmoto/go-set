package set

type Iterator func(item int64) bool

type Interface interface {
	All(Iterator)
	Delete(item int64)
	Append(item int64)
}

type Set struct {
	indexOf map[int64]int
	items   []int64
	deleted []bool
}

func NewSet(c int) *Set {
	return &Set{
		indexOf: make(map[int64]int, c),
		items:   make([]int64, 0, c),
		deleted: make([]bool, 0, c),
	}

}

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

func (s *Set) Delete(item int64) {
	s.deleted[s.indexOf[item]] = true
}

func (s *Set) Append(item int64) {
	_, ok := s.indexOf[item]
	if ok {
		panic("duplicate item")
	}

	s.indexOf[item] = len(s.items)

	s.items = append(s.items, item)
	s.deleted = append(s.deleted, false)
}
