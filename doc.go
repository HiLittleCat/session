/*
Package session provides session middlerware and options
session.Use(&session.RedisProvider{
	Expire: time.Hour,
	Pool:   pool,
	Cookie: http.Cookie{
		Name:     "MySession",
		HttpOnly: true,
		Domain:   "",
		Secure:   false,
	},
})

session.Set(ctx)
session.Get(ctx, _id, map[string]string{"_id": _id, name: "XXX"} )
session.Delete(ctx, store)
*/
package session
