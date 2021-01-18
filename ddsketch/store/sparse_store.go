package store

import (
	"errors"
	"math"
	"sort"
)

type SparseStore struct {
	bins     map[int]int32
	count    int32
	minIndex int
	maxIndex int
}

func NewSparseStore() *SparseStore {
	// TODO: Initialize the bins capacity
	return &SparseStore{minIndex: math.MaxInt32, maxIndex: math.MinInt32}
}

func (s *SparseStore) Add(index int) {
	s.AddWithCount(index, int32(1))
}

func (s *SparseStore) AddBin(bin Bin) {
	index := bin.Index()
	count := bin.Count()
	if count == 0 {
		return
	}
	s.AddWithCount(index, count)
}

func (s *SparseStore) AddWithCount(index int, count int32) {
	if count == 0 {
		return
	}
	if index > s.maxIndex {
		s.maxIndex = index
	}
	if index < s.minIndex {
		s.minIndex = index
	}
	// TODO: have a better growth strategy than double
	s.bins[index] += count
	s.count += count
}

func (s *SparseStore) Bins() <-chan Bin {
	ch := make(chan Bin)
	go func() {
		defer close(ch)
		for k, v := range s.bins {
			ch <- Bin{index: k, count: v}
		}
	}()
	return ch
}

func (s *SparseStore) Copy() Store {
	bins := make(map[int]int32, len(s.bins))
	for k, v := range s.bins {
		bins[k] = v
	}
	return &SparseStore{
		bins:     bins,
		count:    s.count,
		minIndex: s.minIndex,
		maxIndex: s.maxIndex,
	}
}

func (s *SparseStore) IsEmpty() bool {
	return s.count == 0
}

func (s *SparseStore) MaxIndex() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("MaxIndex of empty store is undefined")
	}
	return s.maxIndex, nil
}

func (s *SparseStore) MinIndex() (int, error) {
	if s.IsEmpty() {
		return 0, errors.New("MinIndex of empty store is undefined")
	}
	return s.minIndex, nil
}

func (s *SparseStore) TotalCount() int32 {
	return s.count
}

func (s *SparseStore) KeyAtRank(rank float64) int {
	// map are not ordered in golang
	keys := make([]int, len(s.bins))
	for k, _ := range s.bins {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	var n int32
	for _, k := range keys {
		n += s.bins[k]
		if float64(n) > rank {
			return k
		}
	}
	return s.maxIndex
}

func (s *SparseStore) MergeWith(other Store) {
	if other.IsEmpty() {
		return
	}
	o, ok := other.(*SparseStore)
	if !ok {
		for bin := range other.Bins() {
			s.AddBin(bin)
		}
		return
	}
	// TODO: have a better growth strategy than double
	if o.minIndex < s.minIndex {
		s.minIndex = o.minIndex
	}
	if o.maxIndex > s.maxIndex {
		s.maxIndex = o.maxIndex
	}
	for k, v := range o.bins {
		s.bins[k] += v
	}
	s.count += o.count
}
