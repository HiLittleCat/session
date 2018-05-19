package session

import (
	"errors"
	"net/http"

	"github.com/HiLittleCat/core"
)

// Use 初始化并加载session中间件
func Use(rp *RedisProvider) {
	provider = rp
	provider.Cookie.MaxAge = int(rp.Expire.Seconds())
	core.Use(session)
}

// Get get session
func Get(ctx *core.Context) IStore {
	store, ok := ctx.Data["session"].(IStore)
	if ok == false {
		return nil
	}
	return store
}

// Set set session
func Set(ctx *core.Context, key string, values map[string]string) error {
	store, err := provider.Set(key, values)
	if err != nil {
		return err
	}
	store.Cookie = provider.Cookie
	store.Cookie.Value = store.Values[cookieValueKey]
	ctx.Data["session"] = store
	http.SetCookie(ctx.ResponseWriter, &store.Cookie)
	return nil
}

// Delete delete session
func Delete(ctx *core.Context, store IStore) error {
	ctx.Data["session"] = nil
	st, ok := store.(*RedisStore)
	if ok == false {
		return errors.New("Type IStore cannot convert to *RedisStore")
	}
	sid := store.SessionID()
	provider.Destroy(sid)
	st.Cookie.MaxAge = 1
	http.SetCookie(ctx.ResponseWriter, &st.Cookie)
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
		store.Cookie = provider.Cookie
		store.Cookie.Value = cookie.Value
		ctx.Data["session"] = store
		http.SetCookie(ctx.ResponseWriter, &store.Cookie)
	}

	ctx.Next()
}
