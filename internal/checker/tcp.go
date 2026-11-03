package checker

import (
	"context"
	"crypto/tls"
	"net"
	"time"
)

// CheckTCP attempts to establish a TCP socket connection to host:port.
// If the target is on port 443, it also extracts SSL/TLS certificate metadata.
func CheckTCP(ctx context.Context, target Target, timeout time.Duration) Result {
	dialer := &net.Dialer{
		Timeout: timeout,
	}

	host, port, _ := net.SplitHostPort(target.Address)
	var resolvedIPs []string
	if host != "" {
		ips, _ := net.LookupHost(host)
		resolvedIPs = ips
	}

	start := time.Now()
	conn, err := dialer.DialContext(ctx, "tcp", target.Address)
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
	defer conn.Close()

	var certInfo *CertInfo

	// If connecting to standard TLS port 443, attempt to read certificate chain
	if port == "443" {
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName: host,
			// Allow inspection without hard failing the TCP probe if cert fails verify
			InsecureSkipVerify: true,
		})
		_ = tlsConn.SetDeadline(time.Now().Add(timeout))
		if handshakeErr := tlsConn.HandshakeContext(ctx); handshakeErr == nil {
			state := tlsConn.ConnectionState()
			if len(state.PeerCertificates) > 0 {
				cert := state.PeerCertificates[0]
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
					TLSVersion: tlsVersionString(state.Version),
				}
			}
		}
	}

	return Result{
		Target:      target,
		Up:          true,
		Latency:     latency,
		Err:         nil,
		ResolvedIPs: resolvedIPs,
		Cert:        certInfo,
	}
}
