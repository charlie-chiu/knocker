package knocker

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

func Knock2(site, ip string, port int, withTrace bool) (statusCode int, err error) {
	url, err := url.Parse(site)
	if err != nil {
		return 0, fmt.Errorf("failed to Parse URL: %v", err)
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	// replace ip only if specified
	if ip != "" {
		http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == url.Host+":80" {
				addr = ip + ":80"
			}
			if addr == url.Host+":443" {
				addr = ip + ":443"
			}
			return dialer.DialContext(ctx, network, addr)
		}
	}

	request, err := http.NewRequest(http.MethodGet, url.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to http.NewRequest: %v", err)
	}

	if withTrace {
		trace := &httptrace.ClientTrace{
			DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
				fmt.Printf("DNS Info: %+v\n", dnsInfo)
			},
			GotConn: func(connInfo httptrace.GotConnInfo) {
				fmt.Printf("Got Conn: %+v\n", connInfo)
			},
		}
		request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, fmt.Errorf("failed do a request: %v", err)
	}

	//printStatusCode(response.StatusCode)
	//printSeparateLine()
	//printHeader(response.Header)

	return response.StatusCode, nil
}

func printStatusCode(code int) {
	fmt.Printf("Status Code: %d\n", code)
}

func printSeparateLine() {
	fmt.Printf("---\n")
}

func printHeader(header http.Header) {
	for k, value := range header {
		fmt.Printf("%s: %s\n", k, strings.Join(value, " "))
	}
}
