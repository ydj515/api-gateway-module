package common

import (
	"time"

	"github.com/sony/gobreaker/v2"
)

var CB *gobreaker.CircuitBreaker[[]byte]

func init() {
	var st gobreaker.Settings
	st.Name = "circuit"
	st.MaxRequests = 10
	st.Interval = time.Second

	st.ReadyToTrip = func(counts gobreaker.Counts) bool {
		ratio := float64(counts.TotalFailures) / float64(counts.TotalSuccesses)
		return counts.Requests >= 3 && ratio >= 0.6
	}

	CB = gobreaker.NewCircuitBreaker[[]byte](st)
}
