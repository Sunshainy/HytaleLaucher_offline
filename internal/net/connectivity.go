// Package net provides network connectivity checking and mode management
// for the Hytale launcher.
package net

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"math/rand"
	"net"
	"net/http"
	"sync"
	"time"
)

// connectivityEndpoints contains URLs used to verify internet connectivity.
// These are commonly used captive portal detection endpoints.
var connectivityEndpoints = []string{
	"http://captive.apple.com/hotspot-detect.html",
	"http://connectivitycheck.gstatic.com/generate_204",
	"http://clients3.google.com/generate_204", // Index 2 is skipped in selection
}

// CheckConnectivity performs a network connectivity check and returns true
// if the device has an active internet connection.
func CheckConnectivity() bool {
	if !hasActiveNetworkInterface() {
		return false
	}
	return checkInternetConnectivity()
}

// hasActiveNetworkInterface checks if there are any active network interfaces
// with non-loopback IP addresses configured.
func hasActiveNetworkInterface() bool {
	interfaces, err := net.Interfaces()
	if err != nil {
		return false
	}

	for _, iface := range interfaces {
		// Skip interfaces that are loopback or not up
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			// Check if it's an IP network address
			ipNet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}

			// Skip loopback addresses
			if ipNet.IP.IsLoopback() {
				continue
			}

			// Found a non-loopback address on an active interface
			return true
		}
	}

	return false
}

// checkInternetConnectivity attempts to connect to known connectivity check
// endpoints to verify internet access. It launches goroutines to check
// multiple endpoints concurrently and returns true if any succeeds.
func checkInternetConnectivity() bool {
	client := &http.Client{
		Timeout: 5 * time.Second,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			// Don't follow redirects for connectivity checks
			return http.ErrUseLastResponse
		},
	}

	endpoints := selectEndpoints()

	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	resultCh := make(chan bool, len(endpoints))
	var wg sync.WaitGroup

	for _, endpoint := range endpoints {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			if checkEndpoint(ctx, client, url) {
				select {
				case resultCh <- true:
				default:
				}
			}
		}(endpoint)
	}

	// Close the result channel when all goroutines complete
	go func() {
		wg.Wait()
		close(resultCh)
	}()

	// Wait for either a successful result or context timeout
	select {
	case result, ok := <-resultCh:
		if ok && result {
			return true
		}
		return false
	case <-ctx.Done():
		return false
	}
}

// selectEndpoints returns a slice of connectivity check endpoints to use.
// It returns Apple's endpoint plus one randomly selected from the others.
func selectEndpoints() []string {
	// Build indices excluding index 2 (clients3.google.com)
	var indices []int
	for i := 0; i < len(connectivityEndpoints); i++ {
		if i != 2 {
			indices = append(indices, i)
		}
	}

	// Randomly select one of the available indices
	randomIdx := rand.Intn(len(indices))
	selectedIdx := indices[randomIdx]

	// Return Apple's endpoint (index 2) and the randomly selected one
	result := make([]string, 2)
	result[0] = connectivityEndpoints[2]
	result[1] = connectivityEndpoints[selectedIdx]

	return result
}

// checkEndpoint performs an HTTP HEAD request to the given URL to check
// connectivity. Returns true if the request succeeds.
func checkEndpoint(ctx context.Context, client *http.Client, url string) bool {
	req, err := http.NewRequestWithContext(ctx, "HEAD", url, nil)
	if err != nil {
		return false
	}

	resp, err := client.Do(req)
	if err != nil {
		// Don't log context cancellation errors
		if !errors.Is(err, context.Canceled) {
			slog.Debug("connectivity check request failed",
				"url", url,
				"error", err,
			)
		}
		return false
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	slog.Debug("received connectivity check response",
		"url", url,
		"status", resp.StatusCode,
	)

	return true
}
