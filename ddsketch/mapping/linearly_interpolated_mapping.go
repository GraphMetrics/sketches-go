// Unless explicitly stated otherwise all files in this repository are licensed
// under the Apache License 2.0.
// Copyright 2020 Datadog, Inc. for original work
// Copyright 2021 GraphMetrics for modifications

package mapping

import (
	"bytes"
	"errors"
	"fmt"
	"math"
)

// A fast IndexMapping that approximates the memory-optimal LogarithmicMapping by extracting the floor value
// of the logarithm to the base 2 from the binary representations of floating-point values and linearly
// interpolating the logarithm in-between.
type LinearlyInterpolatedMapping struct {
	relativeAccuracy      float64
	multiplier            float64
	normalizedIndexOffset float64
}

func NewLinearlyInterpolatedMapping(relativeAccuracy float64) (*LinearlyInterpolatedMapping, error) {
	if relativeAccuracy <= 0 || relativeAccuracy >= 1 {
		return nil, errors.New("The relative accuracy must be between 0 and 1.")
	}
	return &LinearlyInterpolatedMapping{
		relativeAccuracy: relativeAccuracy,
		multiplier:       1.0 / math.Log1p(2*relativeAccuracy/(1-relativeAccuracy)),
	}, nil
}

func NewLinearlyInterpolatedMappingWithGamma(gamma, indexOffset float64) (*LinearlyInterpolatedMapping, error) {
	if gamma <= 1 {
		return nil, errors.New("Gamma must be greater than 1.")
	}
	m := LinearlyInterpolatedMapping{
		relativeAccuracy: 1 - 2/(1+math.Exp(math.Log2(gamma))),
		multiplier:       1 / math.Log2(gamma),
	}
	m.normalizedIndexOffset = indexOffset - m.approximateLog(1)*m.multiplier
	return &m, nil
}

func (m *LinearlyInterpolatedMapping) Equals(other IndexMapping) bool {
	o, ok := other.(*LinearlyInterpolatedMapping)
	if !ok {
		return false
	}
	tol := 1e-12
	return withinTolerance(m.multiplier, o.multiplier, tol) && withinTolerance(m.normalizedIndexOffset, o.normalizedIndexOffset, tol)
}

func (m *LinearlyInterpolatedMapping) Index(value float64) int {
	index := m.approximateLog(value)*m.multiplier + m.normalizedIndexOffset
	if index >= 0 {
		return int(index)
	} else {
		return int(index) - 1
	}
}

func (m *LinearlyInterpolatedMapping) Value(index int) float64 {
	return m.approximateInverseLog((float64(index)-m.normalizedIndexOffset)/m.multiplier) * (1 + m.relativeAccuracy)
}

// Return an approximation of log(1) + Math.log(x) / Math.log(2)}
func (m *LinearlyInterpolatedMapping) approximateLog(x float64) float64 {
	bits := math.Float64bits(x)
	return getExponent(bits) + getSignificandPlusOne(bits)
}

// The exact inverse of approximateLog.
func (m *LinearlyInterpolatedMapping) approximateInverseLog(x float64) float64 {
	exponent := math.Floor(x - 1)
	significandPlusOne := x - exponent
	return buildFloat64(int(exponent), significandPlusOne)
}

func (m *LinearlyInterpolatedMapping) MinIndexableValue() float64 {
	return math.Max(
		math.Exp2((math.MinInt16-m.normalizedIndexOffset)/m.multiplier-m.approximateLog(1)+1), // so that index >= MinInt16
		minNormalFloat64*(1+m.relativeAccuracy)/(1-m.relativeAccuracy),
	)
}

func (m *LinearlyInterpolatedMapping) MaxIndexableValue() float64 {
	return math.Min(
		math.Exp2((math.MaxInt16-m.normalizedIndexOffset)/m.multiplier-m.approximateLog(float64(1))-1), // so that index <= MaxInt16
		math.Exp(expOverflow)/(1+m.relativeAccuracy),                                                   // so that math.Exp does not overflow
	)
}

func (m *LinearlyInterpolatedMapping) RelativeAccuracy() float64 {
	return m.relativeAccuracy
}

func (m *LinearlyInterpolatedMapping) string() string {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf("relativeAccuracy: %v, multiplier: %v, normalizedIndexOffset: %v\n", m.relativeAccuracy, m.multiplier, m.normalizedIndexOffset))
	return buffer.String()
}

func withinTolerance(x, y, tolerance float64) bool {
	if x == 0 || y == 0 {
		return math.Abs(x) <= tolerance && math.Abs(y) <= tolerance
	} else {
		return math.Abs(x-y) <= tolerance*math.Max(math.Abs(x), math.Abs(y))
	}
}
