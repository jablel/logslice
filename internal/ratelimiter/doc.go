// Package ratelimiter implements a sliding-window rate limiter for log line
// throughput control.
//
// Usage:
//
//	rl, err := ratelimiter.New(1000) // allow 1000 lines/sec
//	if err != nil {
//		log.Fatal(err)
//	}
//	for _, line := range lines {
//		rl.Wait()          // throttle if needed
//		fmt.Println(line)
//	}
//
// A rate of 0 disables throttling entirely, making Wait a no-op.
// Negative rates are rejected with ErrInvalidRate.
package ratelimiter
