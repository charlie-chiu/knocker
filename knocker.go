package knocker

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
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

type Result struct {
	DNS        []string
	URL        string
	Host       string
	Header     http.Header
	Body       []byte
	StatusCode int
}

func Knock(d Door) (results []Result, err error) {
	u, err := url.Parse(d.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to Parse URL: %v", err)
	}

	// replace ip if specified
	if d.Host != "" {
		http.DefaultTransport.(*http.Transport).DialContext = func(ctx context.Context, network, addr string) (net.Conn, error) {
			if strings.HasPrefix(addr, u.Hostname()) {
				addr = d.Host + strings.TrimPrefix(addr, u.Hostname())
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
		return nil, fmt.Errorf("failed to http.NewRequest: %v", err)
	}

	t := &transport{}
	trace := &httptrace.ClientTrace{
		DNSDone: t.sniffDNSDoneInfo,
		GotConn: t.sniffConnInfo,
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), trace))

	client := &http.Client{Transport: t}
	_, err = client.Do(request)
	if err != nil {
		return nil, fmt.Errorf("failed do a request: %v", err)
	}

	return t.results, nil
}

type transport struct {
	currentRequest *http.Request
	currentLog     Result
	results        []Result
}

func (t *transport) sniffDNSDoneInfo(info httptrace.DNSDoneInfo) {
	var temp []string
	for _, addr := range info.Addrs {
		temp = append(temp, addr.String())
	}
	t.currentLog.DNS = temp
}

func (t *transport) sniffConnInfo(info httptrace.GotConnInfo) {
	t.currentLog.URL = t.currentRequest.URL.String()
	t.currentLog.Host = info.Conn.RemoteAddr().String()
}

func (t *transport) RoundTrip(request *http.Request) (*http.Response, error) {
	t.currentRequest = request
	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		log.Fatalf("transport.RoundTrip err: %v", err)
	}

	t.currentLog.Header = resp.Header
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error on read response body: %v", err)
	}
	t.currentLog.Body = bytes
	t.currentLog.StatusCode = resp.StatusCode

	t.results = append(t.results, t.currentLog)
	t.currentLog = Result{}
	return resp, err
}

func PrintResults(results []Result, withBody bool) {
	for i, r := range results {
		fmt.Printf("\nREQUEST %d\n", i+1)
		fmt.Printf("DNS result: %s\n", r.DNS)
		fmt.Printf("Request URL: %s\n", r.URL)
		fmt.Printf("Request Host: %s\n", r.Host)
		fmt.Printf("Status: %d\n", r.StatusCode)
		printRespHeader(r.Header)
		if withBody {
			fmt.Printf("\nRequest Body:\n%s\n", r.Body)
		}
	}
}

func printRespHeader(header http.Header) {
	fmt.Println("Header:")
	for k, value := range header {
		fmt.Printf("%s: %s\n", k, strings.Join(value, " "))
	}
}
