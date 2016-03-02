package service

import (
	"crypto/rand"
	"encoding/hex"
	"github.com/bradfitz/gomemcache/memcache"
	log "github.com/sirupsen/logrus"
)

const (
	entropyBytes = 32
	ttl          = 5 * 60 // 5 minutes
)

func (s *service) getToken() (string, error) {
	toke, err := genToken()
	if err != nil {
		return "", err
	}

	err = s.memcacheClient.Set(&memcache.Item{
		Key:        toke,
		Value:      []byte("true"),
		Expiration: int32(ttl),
	})

	if err != nil {
		return "", err
	}

	return toke, nil
}

func (s *service) authorizeToken(token string) (bool, error) {
	_, err := s.memcacheClient.Get(token)
	if err != nil {
		if err == memcache.ErrCacheMiss {
			log.Info("cache miss: ", token)
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func genToken() (string, error) {
	buf := make([]byte, entropyBytes)
	_, err := rand.Read(buf)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(buf), nil
}
