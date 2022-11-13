package backoff

import (
	"math"
	"time"
)

// Strategy is used to obtain the next duration to wait after the Nth failure of a request.
type Strategy interface {
	IntervalForRetry(retryNum uint) time.Duration
}

// StrategyConstant implements Strategy and returns the same interval duration
// each time IntervalForRetry is called.
type StrategyConstant struct {
	Interval time.Duration
}

// NewStrategyConstant returns a new instance of StrategyConstant with the provided interval.
func NewStrategyConstant(interval time.Duration) StrategyConstant {
	return StrategyConstant{Interval: interval}
}

// IntervalForRetry returns the same interval value each time it is called.
func (s StrategyConstant) IntervalForRetry(_ uint) time.Duration {
	return s.Interval
}

// StrategyExponential implements Stategy and returns exponentially derived durations
// based on the indicated retry sequence.
type StrategyExponential struct {
	Interval time.Duration
	Base     uint
}

// NewStrategyExponential returns a new instance of StrategyExponential.
func NewStrategyExponential(interval time.Duration, base uint) StrategyExponential {
	return StrategyExponential{Interval: interval, Base: base}
}

// IntervalForRetry returns the duration to wait after the provided number of retries.
// It is computed as interval * base^retryNum
func (s StrategyExponential) IntervalForRetry(retryNum uint) time.Duration {
	return s.Interval * time.Duration(math.Pow(float64(s.Base), float64(retryNum)))
}
