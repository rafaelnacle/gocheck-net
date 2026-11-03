package checker

import (
	"context"
	"net"
	"strings"
	"time"
)

// Protocol indicates the network protocol used for the check.
type Protocol string

const (
	ProtocolHTTP Protocol = "HTTP"
	ProtocolTCP  Protocol = "TCP"
)

// Target represents a network destination to probe.
type Target struct {
	Raw      string   // original input, e.g. "https://google.com" or "github.com:443"
	Protocol Protocol // HTTP or TCP
	Address  string   // parsed URL or host:port
}

// CertInfo holds SSL/TLS certificate metadata.
type CertInfo struct {
	Issuer     string    // e.g. "Let's Encrypt Authority X3"
	CommonName string    // Subject CN
	DNSNames   []string  // Subject Alternative Names (SANs)
	NotBefore  time.Time // Validity start
	NotAfter   time.Time // Expiration timestamp
	DaysLeft   int       // Days until expiration (negative if expired)
	Expired    bool      // True if NotAfter is in the past
	TLSVersion string    // e.g. "TLS 1.3", "TLS 1.2"
}

// Result contains the outcome and metrics of a probe.
type Result struct {
	Target      Target
	Up          bool
	StatusCode  int               // HTTP response status (e.g. 200, 404, 500)
	Latency     time.Duration     // Time taken for probe
	Err         error             // Detailed error if probe failed
	ResolvedIPs []string          // DNS resolved IP addresses
	Cert        *CertInfo         // SSL/TLS certificate details (if HTTPS / TLS)
	Headers     map[string]string // Useful HTTP headers (Server, Content-Type, etc.)
}

// ParseTarget normalizes a raw user input string into a Target struct.
// - "http://..." or "https://..." -> ProtocolHTTP
// - "host:port" -> ProtocolTCP
// - "host" without scheme or port -> defaults to "https://host" (ProtocolHTTP)
func ParseTarget(raw string) Target {
	trimmed := strings.TrimSpace(raw)

	if strings.HasPrefix(trimmed, "http://") || strings.HasPrefix(trimmed, "https://") {
		return Target{
			Raw:      raw,
			Protocol: ProtocolHTTP,
			Address:  trimmed,
		}
	}

	// Check if input has a port: e.g. "example.com:8080" or "127.0.0.1:5432"
	if _, _, err := net.SplitHostPort(trimmed); err == nil {
		return Target{
			Raw:      raw,
			Protocol: ProtocolTCP,
			Address:  trimmed,
		}
	}

	// If no scheme and no port, default to HTTPS
	return Target{
		Raw:      raw,
		Protocol: ProtocolHTTP,
		Address:  "https://" + trimmed,
	}
}

// Check probes a target using its respective protocol within the given timeout.
func Check(ctx context.Context, target Target, timeout time.Duration) Result {
	switch target.Protocol {
	case ProtocolTCP:
		return CheckTCP(ctx, target, timeout)
	case ProtocolHTTP:
		fallthrough
	default:
		return CheckHTTP(ctx, target, timeout)
	}
}
