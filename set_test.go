package set

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func getAllAsInt(fn func(Iterator)) (all []int) {
	it := func(item Item) bool {
		all = append(all, int(item.(Int)))
		return true
	}

	fn(it)

	return
}

func TestSet(t *testing.T) {
	all := []int{0, 1, 2, 3, 4, 5}
	toDelete := []int{0, 2, 4}
	want := []int{1, 3, 5}

	cases := []struct {
		set Interface
	}{
		{
			NewSet(10),
		},
		{
			NewSetMap(10),
		},
	}

	for _, c := range cases {
		for _, item := range all {
			c.set.Append(Int(item))
		}
		for _, item := range toDelete {
			c.set.Delete(Int(item))
		}

		got := getAllAsInt(c.set.All)

		opt := cmpopts.SortSlices(func(x, y int) bool {
			return x < y
		})

		if d := cmp.Diff(want, got, opt); d != "" {
			t.Error(d)
		}

	}

}

func BenchmarkSet(b *testing.B) {

	for _, benchmarkSetSize := range []int64{1000, 10000 * 100} {
		cases := []struct {
			desc   string
			newSet func() Interface
		}{
			{
				desc: "Set",
				newSet: func() Interface {
					return NewSet(int(benchmarkSetSize))
				},
			},
			{
				desc: "SetMap",
				newSet: func() Interface {
					return NewSetMap(int(benchmarkSetSize))
				},
			},
		}

		b.Run(fmt.Sprintf("%d", benchmarkSetSize), func(b *testing.B) {
			for _, c := range cases {
				b.Run(c.desc, func(b *testing.B) {

					b.Run("Append", func(b *testing.B) {
						set := c.newSet()
						var i int64
						b.ResetTimer()
						for i = 0; i < benchmarkSetSize; i++ {
							set.Append(Int(i))
						}
					})
					b.Run("Delete", func(b *testing.B) {
						set := c.newSet()
						var i int64
						for i = 0; i < benchmarkSetSize; i++ {
							set.Append(Int(i))
						}
						b.ResetTimer()
						for i = 0; i < benchmarkSetSize; i++ {
							set.Delete(Int(i))
						}
					})
					b.Run("All", func(b *testing.B) {
						set := c.newSet()
						var i int64
						for i = 0; i < benchmarkSetSize; i++ {
							set.Append(Int(i))
						}
						b.ResetTimer()

						set.All(func(item Item) bool {
							i += int64(item.(Int))
							return true
						})
					})
				})
			}
		})
	}
}

// SetMap provie list which satisfy Interface interface.
// This list simply use map as internal data structure.
type SetMap struct {
	items map[interface{}]Item
}

func NewSetMap(c int) *SetMap {
	return &SetMap{
		items: make(map[interface{}]Item, c),
	}

}

func (s *SetMap) All(fn Iterator) {
	for _, item := range s.items {

		if !fn(item) {
			break
		}
	}
}

func (s *SetMap) Delete(item Item) Item {
	delete(s.items, item.Key())
	return nil // dummy
}

func (s *SetMap) Append(item Item) {
	s.items[item.Key()] = item
}
