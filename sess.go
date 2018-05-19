package session

// IProvider 用以表征session管理器底层存储结构
type IProvider interface {
	Set(rs IStore) error            //设置存储的session
	Get(sid string) (IStore, error) //函数返回sid所代表的Session变量
	Destroy(sid string) error       //函数用来销毁sid对应的Session
	UpExpire()                      //刷新session有效期
}

// IStore session操作
type IStore interface {
	Set(key, value string) error //设置session
	Get(key string) string       //读取session
	Delete(key string) error     //删除session
	SessionID() string           //生成sessionID
}
