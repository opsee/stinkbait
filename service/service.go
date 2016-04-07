package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/bradfitz/gomemcache/memcache"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
	"github.com/julienschmidt/httprouter"
	"github.com/opsee/basic/schema"
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/basic/tp"
	"github.com/opsee/stinkbait/limiter"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

const (
	requestKey = iota
	tokenKey

	rateL
)

var (
	errBadRequest    = errors.New("Bad request.")
	errUnknown       = errors.New("An unknown error happened.")
	errUnauthorized  = errors.New("Unauthorized request.")
	errLimitExceeded = errors.New("request limit exceeded.")
)

type service struct {
	router         *tp.Router
	limiter        *limiter.Limiter
	memcacheClient *memcache.Client
}

type tokenResponse struct {
	Token string `json:"token"`
}

func New(limiter *limiter.Limiter, memcachedList []string) *service {
	s := &service{
		limiter:        limiter,
		memcacheClient: memcache.New(memcachedList...),
	}

	router := tp.NewHTTPRouter(context.Background())

	router.Timeout(1 * time.Minute)
	router.CORS(
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		[]string{`https?://localhost:\d+`, `https?://opsee-ferengi(-staging\d+)?\.s3-website-us-west-2\.amazonaws\.com`, `https://(\w+\.)?(opsy\.co|opsee\.co|opsee\.com)`},
	)

	// get a token
	router.Handle("POST", "/token", []tp.DecodeFunc{}, s.handleToken())

	// check a url
	router.Handle("POST", "/check", []tp.DecodeFunc{
		s.bearerDecodeFunc(),
		s.checkRequestDecodeFunc(),
	}, s.handleCheck())

	// set a custom encoder for proto objects -> json
	router.Encoder("application/json", maybeEncodeProto)

	s.router = router

	log.Info("service initialzed with memcached nodes: %v", memcachedList)

	return s
}

func (s *service) Start(addr, cert, certkey string) error {
	return http.ListenAndServeTLS(addr, cert, certkey, s.router)
}

func (s *service) bearerDecodeFunc() tp.DecodeFunc {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request, p httprouter.Params) (context.Context, int, error) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			log.Error("missing authorization header")
			return nil, http.StatusUnauthorized, errUnauthorized
		}

		var token string
		_, err := fmt.Sscanf(authHeader, "Bearer %s", &token)
		if err != nil {
			log.Error("error decoding bearer token")
			return nil, http.StatusUnauthorized, errUnauthorized
		}

		if waitTime := s.limiter.LimitToken(token); waitTime > 0 {
			log.Info("rate-limited token: %s - wait %d seconds", token, waitTime)
			return nil, http.StatusTooManyRequests, errLimitExceeded
		}

		authed, err := s.authorizeToken(token)
		if err != nil {
			log.WithError(err).Error("error fetching auth token")
			return nil, http.StatusInternalServerError, errUnknown
		}

		if !authed {
			log.Error("unauthorized token")
			return nil, http.StatusUnauthorized, errUnauthorized
		}

		return ctx, 0, nil
	}
}

func (s *service) checkRequestDecodeFunc() tp.DecodeFunc {
	return func(ctx context.Context, rw http.ResponseWriter, r *http.Request, p httprouter.Params) (context.Context, int, error) {
		testCheckRequest := &opsee.TestCheckRequest{}
		err := jsonpb.Unmarshal(r.Body, testCheckRequest)
		if err != nil {
			log.WithError(err).Error("error decoding TestCheckRequest")
			return nil, http.StatusBadRequest, errBadRequest
		}

		err = validateTestCheckRequest(testCheckRequest)
		if err != nil {
			log.WithError(err).Errorf("invalid test check request %#v", testCheckRequest)
			return nil, http.StatusBadRequest, err
		}

		if waitTime := s.limiter.LimitHost(testCheckRequest.Check.Target.Address); waitTime > 0 {
			log.Info("rate-limited host: %s - wait %d seconds", testCheckRequest.Check.Target.Address, waitTime)
			return nil, http.StatusTooManyRequests, errLimitExceeded
		}

		return context.WithValue(ctx, requestKey, testCheckRequest), 0, nil
	}
}

func (s *service) handleCheck() tp.HandleFunc {
	return func(ctx context.Context) (interface{}, int, error) {
		req, ok := ctx.Value(requestKey).(*opsee.TestCheckRequest)
		if !ok {
			log.Error("error decoding TestCheckRequest")
			return nil, http.StatusBadRequest, errBadRequest
		}

		resp, err := s.TestCheck(ctx, req)
		if err != nil {
			log.WithError(err).Error("error executing TestCheck")
			return nil, http.StatusInternalServerError, errUnknown
		}

		return resp, http.StatusOK, nil
	}
}

func (s *service) handleToken() tp.HandleFunc {
	return func(ctx context.Context) (interface{}, int, error) {
		if waitTime := s.limiter.LimitGenerator(); waitTime > 0 {
			log.Info("rate-limited token generation - wait %d seconds", waitTime)
			return nil, http.StatusTooManyRequests, errLimitExceeded
		}

		token, err := s.getToken()
		if err != nil {
			log.WithError(err).Error("error executing getToken")
			return nil, http.StatusInternalServerError, errUnknown
		}

		return tokenResponse{token}, http.StatusOK, nil
	}
}

// encode pb objects with jsonpb, and others with json
func maybeEncodeProto(msg interface{}) ([]byte, error) {
	var (
		buf bytes.Buffer
		err error
	)

	if pmsg, ok := msg.(proto.Message); ok {
		m := jsonpb.Marshaler{}
		if err = m.Marshal(&buf, pmsg); err != nil {
			return nil, err
		}

		return buf.Bytes(), nil
	}

	if err = json.NewEncoder(&buf).Encode(msg); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

var (
	errMissingCheck              = errors.New("missing check.")
	errMissingCheckTarget        = errors.New("missing check target.")
	errMissingCheckTargetAddress = errors.New("missing check target address.")
	errMissingCheckSpec          = errors.New("missing http_check.")
	errMissingCheckSpecPath      = errors.New("missing http_check path.")
	errMissingCheckSpecPort      = errors.New("missing http_check port.")
	errMissingCheckSpecProtocol  = errors.New("missing http_check protocol.")
)

func validateTestCheckRequest(req *opsee.TestCheckRequest) error {
	if req.Check == nil {
		return errMissingCheck
	}

	target := req.Check.Target
	if target == nil {
		return errMissingCheckTarget
	}

	if target.Address == "" {
		return errMissingCheckTargetAddress
	}

	spec, ok := req.Check.Spec.(*schema.Check_HttpCheck)
	if !ok {
		return errMissingCheckSpec
	}

	if spec.HttpCheck == nil {
		return errMissingCheckSpec
	}

	if spec.HttpCheck.Path == "" {
		return errMissingCheckSpecPath
	}

	if spec.HttpCheck.Port == 0 {
		return errMissingCheckSpecPort
	}

	if spec.HttpCheck.Protocol == "" {
		return errMissingCheckSpecProtocol
	}

	return nil
}
