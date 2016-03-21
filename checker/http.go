package checker

import (
	"bufio"
	"bytes"
	"crypto/tls"
	"fmt"
	"github.com/opsee/basic/schema"
	"net"
	"net/http"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const (
	userAgentString = "OpseeBot/0.1-beta (try our monitoring product at https://opsee.com)"
)

// HTTPRequest and HTTPResponse leave their bodies as strings to make life
// easier for now. As soon as we move away from JSON, these should be []byte.

type HTTPRequest struct {
	Method  string           `json:"method"`
	URL     string           `json:"url"`
	Headers []*schema.Header `json:"headers"`
	Body    string           `json:"body"`
}

var (
	// NOTE: http.Client, net.Dialer are safe for concurrent use.
	client *http.Client
)

func init() {
	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig:       &tls.Config{},
			ResponseHeaderTimeout: 30 * time.Second,
			Dial: (&net.Dialer{
				Timeout: 15 * time.Second,
			}).Dial,
		},
	}
}

func NewRequest(address string, check *schema.HttpCheck) *HTTPRequest {
	return &HTTPRequest{
		Method:  "GET", // we only allow get
		URL:     fmt.Sprintf("%s://%s:%d%s", check.Protocol, address, check.Port, check.Path),
		Headers: nil, // we don't allow headers
		Body:    "",  // we don't allow a request body
	}
}

func (r *HTTPRequest) Do() (*schema.HttpResponse, error) {
	req, err := http.NewRequest(r.Method, r.URL, strings.NewReader(r.Body))
	if err != nil {
		return nil, err
	}

	for _, header := range r.Headers {
		key := header.Name

		// we have to special case the host header, since the go client
		// wants that in req.Host
		if strings.ToLower(key) == "host" && len(header.Values) > 0 {
			req.Host = header.Values[0]
		}

		for _, value := range header.Values {
			req.Header.Add(key, value)
		}
	}

	req.Header.Set("User-Agent", userAgentString)

	t0 := time.Now()
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	log.Debug("Attempting to read body of response...")
	// WARNING: You cannot do this.
	//
	// 	body, err := ioutil.ReadAll(resp.Body)
	//
	// We absolutely must limit the size of the body in the response or we will
	// end up using up too much memory. There is no telling how large the bodies
	// could be. If we need to address exceptionally large HTTP bodies, then we
	// can do that in the future.
	//
	// For a breakdown of potential messaging costs, see:
	// https://docs.google.com/a/opsee.co/spreadsheets/d/14Y8DvBkJMhIQoZ11C5_GKeB7NknYyt-fHJaQixkJfKs/edit?usp=sharing

	rdr := bufio.NewReader(resp.Body)
	var contentLength int64

	if resp.ContentLength >= 0 && resp.ContentLength < 131072 {
		contentLength = resp.ContentLength
	} else {
		contentLength = 131072
	}

	body := make([]byte, int64(contentLength))
	if contentLength > 0 {
		rdr.Read(body)
		body = bytes.Trim(body, "\x00")
		body = bytes.Trim(body, "\n")
	}

	httpResponse := &schema.HttpResponse{
		Code: int32(resp.StatusCode),
		Body: string(body),
		Metrics: []*schema.Metric{
			{
				Name:  "request_latency_ms",
				Value: time.Since(t0).Seconds() * 1000,
			},
		},
		Headers: []*schema.Header{},
	}

	for k, v := range resp.Header {
		header := &schema.Header{}
		header.Name = k
		header.Values = v
		httpResponse.Headers = append(httpResponse.Headers, header)
	}

	return httpResponse, nil
}
