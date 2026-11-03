package checker

import (
	"context"
	"crypto/tls"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

// CheckHTTP sends an HTTP GET request to the target address and measures latency.
func CheckHTTP(ctx context.Context, target Target, timeout time.Duration) Result {
	// Clone default transport to allow inspecting certificates without failing early on self-signed/expired certs
	tr := http.DefaultTransport.(*http.Transport).Clone()
	tr.TLSClientConfig = &tls.Config{
		InsecureSkipVerify: true,
	}

	client := &http.Client{
		Timeout:   timeout,
		Transport: tr,
		// Do not automatically follow endless redirects, but allow up to 10
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 10 {
				return http.ErrUseLastResponse
			}
			return nil
		},
	}

	// Resolve DNS IPs in background
	var resolvedIPs []string
	if u, err := url.Parse(target.Address); err == nil && u.Hostname() != "" {
		ips, _ := net.LookupHost(u.Hostname())
		resolvedIPs = ips
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, target.Address, nil)
	if err != nil {
		return Result{
			Target:      target,
			Up:          false,
			Err:         err,
			ResolvedIPs: resolvedIPs,
		}
	}
	req.Header.Set("User-Agent", "gocheck-net/1.0")

	start := time.Now()
	resp, err := client.Do(req)
	latency := time.Since(start)

	if err != nil {
		return Result{
			Target:      target,
			Up:          false,
			Latency:     latency,
			Err:         err,
			ResolvedIPs: resolvedIPs,
		}
	}
	defer resp.Body.Close()

	// Drain a small amount of body to allow connection reuse, discarding errors
	_, _ = io.CopyN(io.Discard, resp.Body, 4096)

	// An endpoint is considered UP if we got an HTTP response below 500 (Server Error)
	isUp := resp.StatusCode < http.StatusInternalServerError

	// Extract TLS certificate details if HTTPS
	var certInfo *CertInfo
	if resp.TLS != nil && len(resp.TLS.PeerCertificates) > 0 {
		cert := resp.TLS.PeerCertificates[0]
		issuer := cert.Issuer.CommonName
		if issuer == "" && len(cert.Issuer.Organization) > 0 {
			issuer = cert.Issuer.Organization[0]
		}
		daysLeft := int(time.Until(cert.NotAfter).Hours() / 24)

		certInfo = &CertInfo{
			Issuer:     issuer,
			CommonName: cert.Subject.CommonName,
			DNSNames:   cert.DNSNames,
			NotBefore:  cert.NotBefore,
			NotAfter:   cert.NotAfter,
			DaysLeft:   daysLeft,
			Expired:    time.Now().After(cert.NotAfter),
			TLSVersion: tlsVersionString(resp.TLS.Version),
		}
	}

	// Extract select response headers
	headers := make(map[string]string)
	for _, key := range []string{"Server", "Content-Type", "Content-Length", "Date"} {
		if val := resp.Header.Get(key); val != "" {
			headers[key] = val
		}
	}

	return Result{
		Target:      target,
		Up:          isUp,
		StatusCode:  resp.StatusCode,
		Latency:     latency,
		Err:         nil,
		ResolvedIPs: resolvedIPs,
		Cert:        certInfo,
		Headers:     headers,
	}
}

// tlsVersionString maps standard library uint16 TLS version codes to human-readable names.
func tlsVersionString(ver uint16) string {
	switch ver {
	case tls.VersionTLS13:
		return "TLS 1.3"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS10:
		return "TLS 1.0"
	default:
		return "TLS (unknown)"
	}
}
