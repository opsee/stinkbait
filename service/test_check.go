package service

import (
	"encoding/json"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/stinkbait/checker"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
)

const (
	responseCacheExpirationSeconds = 5 * 60
)

func (s *service) TestCheck(ctx context.Context, req *opsee.TestCheckRequest) (*opsee.TestCheckResponse, error) {
	httpReq := checker.NewRequest(req.Check.Target.Address, req.Check.GetHttpCheck())
	cachedRespItem, err := s.memcacheClient.Get(httpReq.URL)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			log.Infof("cache miss for url: %s", httpReq.URL)
			return s.testCheck(ctx, req.Check.Target, httpReq)
		}

		log.WithError(err).Error("error fetching from memcached")
		return nil, err
	}

	httpResp := &schema.HttpResponse{}
	err = json.Unmarshal(cachedRespItem.Value, httpResp)
	if err != nil {
		log.WithError(err).Error("error unmarshaling http_response")
		return nil, err
	}

	return testCheckResponse(req.Check.Target, httpResp, ""), nil
}

func (s *service) testCheck(ctx context.Context, target *schema.Target, httpReq *checker.HTTPRequest) (*opsee.TestCheckResponse, error) {
	var errstr string

	resp, err := httpReq.Do()
	if err != nil {
		errstr = err.Error()
	} else {
		// save the response in memcached
		jsond, err := json.Marshal(resp)
		if err != nil {
			log.WithError(err).Error("error marshaling http_response")
			return nil, err
		}

		err = s.memcacheClient.Set(&memcache.Item{
			Key:        httpReq.URL,
			Value:      jsond,
			Expiration: int32(responseCacheExpirationSeconds),
		})

		if err != nil {
			log.WithError(err).Error("error saving http_response in memcached")
			return nil, err
		}
	}

	return testCheckResponse(target, resp, errstr), nil
}

func testCheckResponse(target *schema.Target, response *schema.HttpResponse, errstr string) *opsee.TestCheckResponse {
	return &opsee.TestCheckResponse{
		Responses: []*schema.CheckResponse{
			{
				Target: target,
				Error:  errstr,
				Reply: &schema.CheckResponse_HttpResponse{
					HttpResponse: response,
				},
			},
		},
		Error: errstr,
	}
}
