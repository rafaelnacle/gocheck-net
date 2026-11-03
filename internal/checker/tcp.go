package checker

import (
	"context"
	"net"
	"time"
)

// CheckTCP attempts to establish a raw TCP socket connection to host:port.
func CheckTCP(ctx context.Context, target Target, timeout time.Duration) Result {
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", target.Address)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Target:  target,
			Up:      false,
			Latency: latency,
			Err:     err,
		}
	}
	defer conn.Close()

	return Result{
		Target:  target,
		Up:      true,
		Latency: latency,
		Err:     nil,
	}
}
