// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License 2.0.
// Copyright 2020 Datadog, Inc. for original work
// Copyright 2021 GraphMetrics for modifications

package store

import (
	"errors"
)

type Bin struct {
	index int
	count int32
}

func NewBin(index int, count int32) (*Bin, error) {
	if count < 0 {
		return nil, errors.New("count cannot be negative")
	}
	return &Bin{index: index, count: count}, nil
}

func (b *Bin) Index() int {
	return b.index
}

func (b *Bin) Count() int32 {
	return b.count
}
