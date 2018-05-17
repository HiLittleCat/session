package session

import (
	"net/http"
	"time"

	"github.com/HiLittleCat/conn"
	redis "gopkg.in/redis.v5"
)

var (
	storeOption StoreOption
	storeS      *Store
)

// Values session values
type Values map[string]string

// StoreOption store option
type StoreOption struct {
	Expire time.Duration
	Pool   *conn.RedisPool
	Cookie http.Cookie
}

// Store redis session store
type Store struct {
	ID     string
	Values Values
	Cookie http.Cookie
}

// Generate init a session by the key
func (rs *Store) Generate(key string, values Values) *Store {
	store := Store{
		ID:     rs.sessionID(key),
		Values: values,
		Cookie: storeOption.Cookie,
	}
	return &store
}

// Flush the session to redis
func (rs *Store) Flush() error {
	var err error
	storeOption.Pool.Exec(func(c *redis.Client) {
		err = c.HMSet(rs.ID, rs.Values).Err()
		if err != nil {
			return
		}
		err = c.Expire(rs.ID, storeOption.Expire).Err()
	})

	return err
}

// Get get a session from redis
func (rs *Store) Get() error {
	var val map[string]string
	var err error
	storeOption.Pool.Exec(func(c *redis.Client) {
		val, err = c.HGetAll(rs.ID).Result()
		rs.Values = val
	})
	return err
}

// Delete generate a session by sid
func (rs *Store) Delete() error {
	var err error
	storeOption.Pool.Exec(func(c *redis.Client) {
		err = c.Del(rs.ID).Err()
	})
	return err
}

// GetFieldValue get key value
func (rs *Store) GetFieldValue(key string) string {
	return rs.Values[key]
}

// SetFieldValue set key value in redis session
func (rs *Store) SetFieldValue(key string, value string) error {
	var err error
	storeOption.Pool.Exec(func(c *redis.Client) {
		err = c.HSet(rs.ID, key, value).Err()
		if err != nil {
			return
		}
		err = c.Expire(rs.ID, storeOption.Expire).Err()
	})
	if err == nil {
		rs.Values[key] = value
	}
	return err
}

// UpExpire update session expire
func (rs *Store) UpExpire() error {
	var err error
	storeOption.Pool.Exec(func(c *redis.Client) {
		err = c.Expire(rs.ID, storeOption.Expire).Err()
	})
	return err
}

// SessionID generate redis session id by value
func (rs *Store) sessionID(value string) string {
	return value
}
