/*
Package session provides session middlerware and options
session.Use(&session.redisProvider{
	Expire: time.Hour,
	Pool:   pool,
	Cookie: http.Cookie{
		Name:     "MySession",
		HttpOnly: true,
		Domain:   "",
		Secure:   false,
	},
})

session.Get(ctx)
session.Set(ctx, _id, map[string]string{"_id": _id, "name": "XXX"} )
session.Delete(ctx, store)
*/
package session
