package knocker

import (
	"context"
	"crypto/tls"
	"fmt"
	"io/ioutil"
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
	Error      error
}

func Knock(d Door) (results []Result, err error) {
	u, err := url.Parse(d.URL)
	if err != nil {
		return nil, fmt.Errorf("failed to Parse URL: %v", err)
	}

	// replace addr if specified
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

	s := &sniffer{}

	request, err := http.NewRequest(http.MethodGet, u.String(), http.NoBody)
	if err != nil {
		return s.results, fmt.Errorf("failed to http.NewRequest: %v", err)
	}
	request = request.WithContext(httptrace.WithClientTrace(request.Context(), &httptrace.ClientTrace{
		DNSDone: s.sniffDNSDoneInfo,
		GotConn: s.sniffConnInfo,
	}))

	client := &http.Client{Transport: s}
	_, err = client.Do(request)
	if err != nil {
		return s.results, fmt.Errorf("failed do a request: %v", err)
	}

	return s.results, nil
}

type sniffer struct {
	request *http.Request
	result  Result
	results []Result
}

func (s *sniffer) sniffDNSDoneInfo(info httptrace.DNSDoneInfo) {
	var temp []string
	for _, addr := range info.Addrs {
		temp = append(temp, addr.String())
	}
	s.result.DNS = temp
}

func (s *sniffer) sniffConnInfo(info httptrace.GotConnInfo) {
	s.result.URL = s.request.URL.String()
	s.result.Host = info.Conn.RemoteAddr().String()
}

func (s *sniffer) RoundTrip(request *http.Request) (*http.Response, error) {
	s.request = request
	resp, err := http.DefaultTransport.RoundTrip(request)
	if err != nil {
		s.result.Error = err
		s.finishRequest()
		return nil, err
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		s.result.Error = err
		s.finishRequest()
		return nil, err
	} else {
		s.result.Body = bytes
		resp.Body.Close()
	}

	s.result.StatusCode = resp.StatusCode
	s.result.Header = resp.Header
	s.finishRequest()
	return resp, nil
}

func (s *sniffer) finishRequest() {
	s.results = append(s.results, s.result)
	s.result = Result{}
}

func PrintResults(results []Result, withBody bool) {
	for i, r := range results {
		fmt.Printf("\nREQUEST %d\n", i+1)
		if r.Error != nil {
			fmt.Printf("\n### ERROR ##\n")
			fmt.Printf("%v", r.Error)
			fmt.Printf("\n### ERROR ##\n")
		}
		fmt.Printf("DNS result: %s\n", strings.Join(r.DNS, ", "))
		fmt.Printf("Request URL: %s\n", r.URL)
		fmt.Printf("Request Host: %s\n", r.Host)
		fmt.Printf("Status: %d\n", r.StatusCode)
		printRespHeader(r.Header)
		if withBody {
			fmt.Printf("\nResponse Body:\n%s\n", r.Body)
		}
		fmt.Printf("\nEND of REQUEST %d\n", i+1)
	}
}

func printRespHeader(header http.Header) {
	fmt.Println("Header:")
	for k, value := range header {
		fmt.Printf("%s: %s\n", k, strings.Join(value, " "))
	}
}
