// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License 2.0.
// Copyright 2020 Datadog, Inc. for original work
// Copyright 2021 GraphMetrics for modifications

package mapping

const (
	expOverflow      = 7.094361393031e+02      // The value at which math.Exp overflows
	minNormalFloat64 = 2.2250738585072014e-308 //2^(-1022)
)

type IndexMapping interface {
	Equals(other IndexMapping) bool
	Index(value float64) int
	Value(index int) float64
	RelativeAccuracy() float64
	MinIndexableValue() float64
	MaxIndexableValue() float64
}
