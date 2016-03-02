package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/protobuf/opseeproto/types"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestTestCheckRequest(t *testing.T) {
	var (
		assert = assert.New(t)
		server = setup()
	)

	// check req with no token
	recorder := server.testCheckRequest(t, "https", "www.reddit.com", "/r/pepe", int32(443), nil)
	assert.Equal(401, recorder.Code, "requests without token are unauthorized")

	// check req with bad token
	recorder = server.testCheckRequest(t, "https", "www.reddit.com", "/r/pepe", int32(443), map[string]string{
		"Authorization": "Bearer blibbleblibble",
	})
	assert.Equal(401, recorder.Code, "requests with unknown token are unauthorized")

	// token req
	tokenresp := &tokenResponse{}
	recorder = server.testTokenRequest(t)
	assert.Equal(200, recorder.Code, "was able to fetch an auth token")

	err := json.NewDecoder(recorder.Body).Decode(tokenresp)
	if err != nil {
		t.Fatal(err)
	}

	// check req with good token
	recorder = server.testCheckRequest(t, "https", "www.reddit.com", "/r/pepe", int32(443), map[string]string{
		"Authorization": fmt.Sprintf("Bearer %s", tokenresp.Token),
	})
	assert.Equal(200, recorder.Code, "requests with a valid token are authorized")

	testCheckResponse := &opsee.TestCheckResponse{}
	err = jsonpb.Unmarshal(recorder.Body, testCheckResponse)
	if err != nil {
		t.Fatal(err)
	}

	assert.True(strings.Contains(testCheckResponse.Responses[0].GetHttpResponse().Body, "reddit"))
}

func setup() *service {
	log.SetLevel(log.FatalLevel)
	viper.SetEnvPrefix("stinkbait")
	viper.AutomaticEnv()

	return New(viper.GetStringSlice("memcached_nodes"))
}

func (s *service) testTokenRequest(t *testing.T) *httptest.ResponseRecorder {
	request, err := http.NewRequest("POST", "https://foo/token", nil)
	if err != nil {
		t.Fatal(err)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, request)

	return w
}

func (s *service) testCheckRequest(t *testing.T, proto, host, path string, port int32, headers map[string]string) *httptest.ResponseRecorder {
	ts := &types.Timestamp{}
	ts.Scan(time.Now())

	req := &opsee.TestCheckRequest{
		Deadline: ts,
		Check: &schema.Check{
			Id: host,
			Target: &schema.Target{
				Name:    host,
				Id:      host,
				Type:    "url",
				Address: host,
			},
			Name: host,
			Spec: &schema.Check_HttpCheck{
				&schema.HttpCheck{
					Path:     path,
					Protocol: proto,
					Port:     port,
					Verb:     "GET",
				},
			},
		},
	}

	marshaler := &jsonpb.Marshaler{}
	reqJson, err := marshaler.MarshalToString(req)
	if err != nil {
		t.Fatal(err)
	}
	// t.Log("request body: ", reqJson)

	request, err := http.NewRequest("POST", "https://foo/check", bytes.NewBufferString(reqJson))
	if err != nil {
		t.Fatal(err)
	}

	for k, v := range headers {
		request.Header.Add(k, v)
	}

	w := httptest.NewRecorder()
	s.router.ServeHTTP(w, request)

	return w
}
