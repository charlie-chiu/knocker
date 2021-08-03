package knocker

import (
	"context"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"net/url"
	"strings"
	"time"
)

type Door struct {
	URL       string
	Host      string
	IgnoreSSL bool
}

func Knock(d Door) (statusCode int, err error) {
	u, err := url.Parse(d.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to Parse URL: %v", err)
	}

	// replace ip if specified
	if d.Host != "" {
		http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if addr == u.Hostname()+":"+u.Port() {
				addr = d.Host + ":" + u.Port()
			}
			return (&net.Dialer{
				Timeout:   30 * time.Second,
				KeepAlive: 30 * time.Second,
			}).DialContext(ctx, network, addr)
		}
	}

	if d.IgnoreSSL {
		http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	request, err := http.NewRequest(http.MethodGet, u.String(), nil)
	if err != nil {
		return 0, fmt.Errorf("failed to http.NewRequest: %v", err)
	}

	t := &transport{}
	trace := &httptrace.ClientTrace{
		DNSDone: func(dnsInfo httptrace.DNSDoneInfo) {
			printDebugHeader("DNS resoling...")
			fmt.Printf("DNS resolve result: %+v\n", dnsInfo.Addrs)
		},
		GotConn: t.printConnInfo,
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))

	client := &http.Client{Transport: t}
	response, err := client.Do(request)
	if err != nil {
		return 0, fmt.Errorf("failed do a request: %v", err)
	}

	printDebugHeader("request completed")

	return response.StatusCode, nil
}

type transport struct {
	current *http.Request
}

func (t *transport) printConnInfo(info httptrace.GotConnInfo) {
	// we can use request info here
	printDebugHeader("GotConn")
	fmt.Printf("request %s @ %s...\n", t.current.URL, info.Conn.RemoteAddr())
}

func (t *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	t.current = request
	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		log.Fatalf("transport.RoundTrip err: %v", err)
	}

	printDebugHeader("RoundTrip")
	printStatusCode(resp.StatusCode)
	printHeader(resp.Header)

	return resp, err
}
func printStatusCode(code int) {
	fmt.Printf("Status Code: %d\n", code)
}

func printDebugHeader(msg string) {
	fmt.Printf("\n### %s ###\n", msg)
}

func printHeader(header http.Header) {
	for k, value := range header {
		fmt.Printf("%s: %s\n", k, strings.Join(value, " "))
	}
}
