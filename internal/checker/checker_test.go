package checker_test

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"gocheck-net/internal/checker"
)

func TestParseTarget(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedProtocol checker.Protocol
		expectedAddress  string
	}{
		{
			name:             "HTTP URL with scheme",
			input:            "http://example.com/status",
			expectedProtocol: checker.ProtocolHTTP,
			expectedAddress:  "http://example.com/status",
		},
		{
			name:             "HTTPS URL with scheme",
			input:            "https://google.com",
			expectedProtocol: checker.ProtocolHTTP,
			expectedAddress:  "https://google.com",
		},
		{
			name:             "TCP host:port target",
			input:            "localhost:5432",
			expectedProtocol: checker.ProtocolTCP,
			expectedAddress:  "localhost:5432",
		},
		{
			name:             "Domain name without scheme defaults to HTTPS",
			input:            "github.com",
			expectedProtocol: checker.ProtocolHTTP,
			expectedAddress:  "https://github.com",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res := checker.ParseTarget(tt.input)
			if res.Protocol != tt.expectedProtocol {
				t.Errorf("ParseTarget(%q).Protocol = %v, want %v", tt.input, res.Protocol, tt.expectedProtocol)
			}
			if res.Address != tt.expectedAddress {
				t.Errorf("ParseTarget(%q).Address = %q, want %q", tt.input, res.Address, tt.expectedAddress)
			}
		})
	}
}

func TestCheckHTTP(t *testing.T) {
	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Server", "TestMockServer/1.0")
		switch r.URL.Path {
		case "/200":
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("OK"))
		case "/500":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Error"))
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer server.Close()

	tests := []struct {
		name         string
		url          string
		expectedUp   bool
		expectedCode int
	}{
		{
			name:         "HTTP 200 OK is UP",
			url:          server.URL + "/200",
			expectedUp:   true,
			expectedCode: http.StatusOK,
		},
		{
			name:         "HTTP 404 is UP (reachable)",
			url:          server.URL + "/404",
			expectedUp:   true,
			expectedCode: http.StatusNotFound,
		},
		{
			name:         "HTTP 500 is DOWN (server error)",
			url:          server.URL + "/500",
			expectedUp:   false,
			expectedCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			target := checker.ParseTarget(tt.url)
			res := checker.CheckHTTP(context.Background(), target, 2*time.Second)

			if res.Up != tt.expectedUp {
				t.Errorf("CheckHTTP(%s) Up = %v, want %v", tt.url, res.Up, tt.expectedUp)
			}
			if res.StatusCode != tt.expectedCode {
				t.Errorf("CheckHTTP(%s) StatusCode = %d, want %d", tt.url, res.StatusCode, tt.expectedCode)
			}
			if res.Headers["Server"] != "TestMockServer/1.0" {
				t.Errorf("CheckHTTP expected Server header, got %q", res.Headers["Server"])
			}
		})
	}
}

func TestCheckHTTP_TLS(t *testing.T) {
	tlsServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("secure"))
	}))
	defer tlsServer.Close()

	target := checker.ParseTarget(tlsServer.URL)
	res := checker.CheckHTTP(context.Background(), target, 2*time.Second)

	if !res.Up {
		t.Fatalf("Expected TLS server to be UP, got error: %v", res.Err)
	}
	if res.Cert == nil {
		t.Fatalf("Expected CertInfo to be populated for HTTPS server")
	}
	if res.Cert.TLSVersion == "" {
		t.Errorf("Expected non-empty TLS version")
	}
	if res.Cert.DaysLeft <= 0 && res.Cert.Expired {
		t.Errorf("Expected newly generated test cert not to be expired")
	}
}

func TestCheckTCP(t *testing.T) {
	// Start a local TCP listener
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("Failed to start test TCP listener: %v", err)
	}
	defer listener.Close()

	// Accept connections in background
	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				return
			}
			conn.Close()
		}
	}()

	t.Run("Open TCP port reports UP", func(t *testing.T) {
		target := checker.ParseTarget(listener.Addr().String())
		res := checker.CheckTCP(context.Background(), target, 1*time.Second)
		if !res.Up {
			t.Errorf("CheckTCP on %s: expected Up=true, got Up=false (%v)", listener.Addr().String(), res.Err)
		}
	})

	t.Run("Closed TCP port reports DOWN", func(t *testing.T) {
		// Use a port that should not have anything listening
		target := checker.ParseTarget("127.0.0.1:59999")
		res := checker.CheckTCP(context.Background(), target, 200*time.Millisecond)
		if res.Up {
			t.Errorf("CheckTCP on closed port: expected Up=false, got Up=true")
		}
	})
}
