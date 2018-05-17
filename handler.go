package session

import (
	"net/http"

	"github.com/HiLittleCat/core"
)

// Use 初始化并加载session中间件
func Use(options StoreOption) {
	storeOption = options
	storeOption.Cookie.MaxAge = int(options.Expire.Seconds())
	core.Use(session)
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

	store := storeS.Generate(cookie.Value, map[string]string{})

	err := store.Get()
	if err != nil {
		ctx.Fail(err)
		return
	}

	if len(store.Values) > 0 {
		err := store.UpExpire()

		if err != nil {
			ctx.Fail(err)
			return
		}
		store.Cookie.Value = cookie.Value
		ctx.Data["session"] = store
		http.SetCookie(ctx.ResponseWriter, &store.Cookie)
	}

	ctx.Next()
}
