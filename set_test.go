package set

// package set implements list which can delete item, specially intended to be able to get all items fast.

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func getAll(fn func(Iterator)) (all []int64) {
	it := func(item int64) bool {
		all = append(all, item)
		return true
	}

	fn(it)

	return
}

func TestSet(t *testing.T) {
	all := []int64{0, 1, 2, 3, 4, 5}
	toDelete := []int64{0, 2, 4}
	want := []int64{1, 3, 5}

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
			c.set.Append(item)
		}
		for _, item := range toDelete {
			c.set.Delete(item)
		}

		got := getAll(c.set.All)

		opt := cmpopts.SortSlices(func(x, y int64) bool {
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
							set.Append(i)
						}
					})
					b.Run("Delete", func(b *testing.B) {
						set := c.newSet()
						var i int64
						for i = 0; i < benchmarkSetSize; i++ {
							set.Append(i)
						}
						b.ResetTimer()
						for i = 0; i < benchmarkSetSize; i++ {
							set.Delete(i)
						}
					})
					b.Run("All", func(b *testing.B) {
						set := c.newSet()
						var i int64
						for i = 0; i < benchmarkSetSize; i++ {
							set.Append(i)
						}
						b.ResetTimer()

						set.All(func(item int64) bool {
							i += item
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
	items map[int64]struct{}
}

func NewSetMap(c int) *SetMap {
	return &SetMap{
		items: make(map[int64]struct{}, c),
	}

}

func (s *SetMap) All(fn Iterator) {
	for item, _ := range s.items {

		if !fn(item) {
			break
		}
	}
}

func (s *SetMap) Delete(item int64) {
	delete(s.items, item)
}

func (s *SetMap) Append(item int64) {
	s.items[item] = struct{}{}
}
