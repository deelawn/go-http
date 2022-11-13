package backoff

import (
	"math"
	"time"
)

type Strategy interface {
	IntervalForRetry(retryNum uint) time.Duration
}

type StrategyConstant struct {
	Interval time.Duration
}

func NewStrategyConstant(interval time.Duration) StrategyConstant {
	return StrategyConstant{Interval: interval}
}

func (s StrategyConstant) IntervalForRetry(retryNum uint) time.Duration {
	return s.Interval * time.Duration(retryNum)
}

type StrategyExponential struct {
	Interval time.Duration
	Base     uint
}

func NewStrategyExponential(interval time.Duration, base uint) StrategyExponential {
	return StrategyExponential{Interval: interval, Base: base}
}

func (s StrategyExponential) IntervalForRetry(retryNum uint) time.Duration {
	return s.Interval * time.Duration(math.Pow(float64(s.Base), float64(s.Interval)))
}
