package session

import (
	"github.com/zyl0501/go-push/api"
	"github.com/zyl0501/go-push/tools/config"
	"time"
	"github.com/zyl0501/go-push/tools/crypto"
	"strconv"
)

var sessionExpireTime = config.SessionExpireTime

type ReusableSession struct {
	SessionId  string
	ExpireTime int64
	Context    api.SessionContext
}

func NewSession(context api.SessionContext) *ReusableSession {
	now := time.Now()
	session := ReusableSession{}
	session.Context = context;
	session.SessionId = crypto.MD5([]byte(context.DeviceId + strconv.FormatInt(now.Unix(), 10)));
	session.ExpireTime = now.Add(sessionExpireTime).Unix();
	return &session
}

type ReusableSessionManager struct {
	cache map[string]ReusableSession
}

func NewReusableSessionManager() (*ReusableSessionManager) {
	return &ReusableSessionManager{make(map[string]ReusableSession)}
}

func (manager *ReusableSessionManager) CacheSession(session *ReusableSession) {
	manager.cache[getSessionKey(session.SessionId)] = *session
}

func (manager *ReusableSessionManager) QuerySession(sessionId string) (*ReusableSession) {
	session := manager.cache[sessionId]
	return &session
}

func getSessionKey(sessionId string) string {
	return "mp:rs:" + sessionId;
}
