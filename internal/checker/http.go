package checker

import (
	"context"
	"io"
	"net/http"
	"time"
)

// CheckHTTP sends an HTTP GET request to the target address and measures latency.
func CheckHTTP(ctx context.Context, target Target, timeout time.Duration) Result {
	client := &http.Client{
		Timeout: timeout,
		// Do not automatically follow endless redirects, but allow up to 10
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.Address, nil)
	if err != nil {
		return Result{
			Target: target,
			Up:     false,
			Err:    err,
		}
	}
	req.Header.Set("User-Agent", "gocheck-net/1.0")

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Target:  target,
			Up:      false,
			Latency: latency,
			Err:     err,
		}
	}
	defer resp.Body.Close()

	// Drain a small amount of body to allow connection reuse, discarding errors
	_, _ = io.CopyN(io.Discard, resp.Body, 4096)

	// An endpoint is considered UP if we got an HTTP response below 500 (Server Error)
	isUp := resp.StatusCode < http.StatusInternalServerError

	return Result{
		Target:     target,
		Up:         isUp,
		StatusCode: resp.StatusCode,
		Latency:    latency,
		Err:        nil,
	}
}
