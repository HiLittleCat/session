package session

import (
	"net/http"
	"time"

	"github.com/HiLittleCat/conn"
	"github.com/HiLittleCat/core"
)

// Use 初始化并加载session中间件
func Use(expire time.Duration, pool *conn.RedisPool, cookie http.Cookie) {
	sessExpire = expire
	redisPool = pool
	httpCookie = cookie
	httpCookie.MaxAge = int(sessExpire.Seconds())
	core.Use(session)
}

// Get get session
func Get(ctx *core.Context) IStore {
	store := ctx.Data["session"]
	if store == nil {
		return nil
	}
	st, ok := store.(IStore)
	if ok == false {
		return nil
	}
	return st
}

// Set set session
func Set(ctx *core.Context, key string, values map[string]string) error {
	store, err := provider.Set(key, values)
	if err != nil {
		return err
	}
	cookie := httpCookie
	cookie.Value = store.Values[cookieValueKey]
	ctx.Data["session"] = store
	http.SetCookie(ctx.ResponseWriter, &cookie)
	return nil
}

// FreshExpire set session
func FreshExpire(ctx *core.Context, key string) error {
	err := provider.UpExpire(key)
	if err != nil {
		return err
	}
	return nil
}

// Delete delete session
func Delete(ctx *core.Context, sid string) error {
	ctx.Data["session"] = nil
	provider.Destroy(sid)
	cookie := httpCookie
	cookie.MaxAge = 1
	http.SetCookie(ctx.ResponseWriter, &cookie)
	return nil
}

// session session处理
func session(ctx *core.Context) {
	var cookie *http.Cookie
	cookies := ctx.Request.Cookies()
	if len(cookies) > 0 {
		cookie = cookies[0]
	} else {
		ctx.Next()
		return
	}
	sid := cookie.Value
	store, err := provider.Get(sid)
	if err != nil {
		ctx.Fail(err)
		return
	}

	if len(store.Values) > 0 {
		err := provider.refresh(store)
		if err != nil {
			ctx.Fail(err)
			return
		}
		cookie := httpCookie
		cookie.Value = store.Values[cookieValueKey]
		ctx.Data["session"] = store
		http.SetCookie(ctx.ResponseWriter, &cookie)
	}

	ctx.Next()
}
