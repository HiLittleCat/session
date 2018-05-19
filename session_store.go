package session

import (
	"net/http"
	"time"

	"github.com/HiLittleCat/conn"
	redis "gopkg.in/redis.v5"
)

// provider redis session provider
var provider = &RedisProvider{}

var cookieValueKey = "_id"

// redisStore session store
type redisStore struct {
	SID    string
	Values map[string]string
	Cookie http.Cookie
}

// Set value
func (rs *redisStore) Set(key, value string) error {
	rs.Values[key] = value
	if key == cookieValueKey {
		rs.Cookie.Value = value
	}
	err := provider.refresh(rs)
	return err
}

// Get value
func (rs *redisStore) Get(key string) string {
	if v, ok := rs.Values[key]; ok {
		return v
	}
	return ""
}

// Delete value in redis session
func (rs *redisStore) Delete(key string) error {
	delete(rs.Values, key)
	if key == cookieValueKey {
		rs.Cookie.Value = ""
	}
	rs.Cookie.MaxAge = 1
	err := provider.refresh(rs)
	return err
}

// SessionID get redis session id
func (rs *redisStore) SessionID() string {
	return rs.SID
}

// RedisProvider redis session RedisProvider
type RedisProvider struct {
	Expire time.Duration
	Pool   *conn.RedisPool
	Cookie http.Cookie
}

// Set value in redis session
func (rp *RedisProvider) Set(key string, values map[string]string) (*redisStore, error) {
	rs := &redisStore{SID: key, Values: values, Cookie: provider.Cookie}
	rs.Cookie.Value = values[cookieValueKey]
	err := provider.refresh(rs)
	return rs, err
}

// refresh refresh store to redis
func (rp *RedisProvider) refresh(rs *redisStore) error {
	var err error
	rp.Pool.Exec(func(c *redis.Client) {
		err = c.HMSet(rs.SID, rs.Values).Err()
		if err != nil {
			return
		}
		err = c.Expire(rs.SID, rp.Expire).Err()
	})
	return nil
}

// Get read redis session by sid
func (rp *RedisProvider) Get(sid string) (*redisStore, error) {
	var rs = &redisStore{}
	var val map[string]string
	var err error
	rp.Pool.Exec(func(c *redis.Client) {
		val, err = c.HGetAll(sid).Result()
		rs.Values = val
	})
	return rs, err
}

// Destroy delete redis session by id
func (rp *RedisProvider) Destroy(sid string) error {
	var err error
	rp.Pool.Exec(func(c *redis.Client) {
		err = c.Del(sid).Err()
	})
	return err
}

// UpExpire refresh session expire
func (rp *RedisProvider) UpExpire(sid string) error {
	var err error
	rp.Pool.Exec(func(c *redis.Client) {
		err = c.Expire(sid, rp.Expire).Err()
	})
	return err
}
