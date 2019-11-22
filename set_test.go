package set

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func getAll(fn func(Iterator)) (all []interface{}) {
	it := func(item interface{}) bool {
		all = append(all, item)
		return true
	}

	fn(it)

	return
}

func TestSet(t *testing.T) {
	all := []interface{}{0, 1, 2, 3, 4, 5}
	toDelete := []interface{}{0, 2, 4}
	want := []interface{}{1, 3, 5}

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

		opt := cmpopts.SortSlices(func(x, y interface{}) bool {
			return x.(int) < y.(int)
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

						set.All(func(item interface{}) bool {
							i += item.(int64)
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
	items map[interface{}]struct{}
}

func NewSetMap(c int) *SetMap {
	return &SetMap{
		items: make(map[interface{}]struct{}, c),
	}

}

func (s *SetMap) All(fn Iterator) {
	for item, _ := range s.items {

		if !fn(item) {
			break
		}
	}
}

func (s *SetMap) Delete(item interface{}) {
	delete(s.items, item)
}

func (s *SetMap) Append(item interface{}) {
	s.items[item] = struct{}{}
}
