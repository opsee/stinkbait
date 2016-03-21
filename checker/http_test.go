package checker

import (
	"fmt"
	"github.com/opsee/basic/schema"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"
	"strings"
	"testing"
)

func TestRequestDo(t *testing.T) {
	assert := assert.New(t)
	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("X-Pepe", "rare")
		if r.URL.Path == "/huge" {
			fmt.Fprintln(w, strings.Repeat("Getout", 131072))
		} else {
			fmt.Fprintln(w, "Getout")
		}
	}))

	resp, err := newTestRequest(t, testServer.URL).Do()
	if err != nil {
		t.Fatal(err)
	}

	assert.EqualValues(200, resp.Code)
	assert.Equal("Getout", resp.Body)
	assert.Contains(resp.Headers, &schema.Header{Name: "X-Pepe", Values: []string{"rare"}})

	resp, err = newTestRequest(t, testServer.URL+"/huge").Do()
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(131072, len(resp.Body))
}

func newTestRequest(t *testing.T, checkUrl string) *HTTPRequest {
	url, err := url.Parse(checkUrl)
	if err != nil {
		t.Fatal(err)
	}

	var (
		hostPort = strings.Split(url.Host, ":")
		port     = 80
		host     = hostPort[0]
	)

	if url.Scheme == "https" {
		port = 443
	}

	if len(hostPort) > 1 {
		port, err = strconv.Atoi(hostPort[1])
		if err != nil {
			t.Fatal(err)
		}
	}

	return NewRequest(host, &schema.HttpCheck{
		Path:     url.Path,
		Protocol: url.Scheme,
		Port:     int32(port),
		Verb:     "GET",
	})
}
