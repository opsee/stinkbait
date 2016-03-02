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
	opsee "github.com/opsee/basic/service"
	"github.com/opsee/basic/tp"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/context"
	"net/http"
	"time"
)

const (
	requestKey = iota
	tokenKey
)

var (
	errBadRequest   = errors.New("Bad request.")
	errUnknown      = errors.New("An unknown error happened.")
	errUnauthorized = errors.New("Unauthorized request.")
)

type service struct {
	router         *tp.Router
	memcacheClient *memcache.Client
}

type tokenResponse struct {
	Token string `json:"token"`
}

func New(memcachedList []string) *service {
	s := &service{memcacheClient: memcache.New(memcachedList...)}

	router := tp.NewHTTPRouter(context.Background())

	router.Timeout(1 * time.Minute)
	router.CORS(
		[]string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		[]string{`https?://localhost:8080`, `https?://localhost:8008`, `https://(\w+\.)?(opsy\.co|opsee\.co|opsee\.com)`},
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
