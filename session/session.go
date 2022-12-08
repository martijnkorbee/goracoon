package session

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/alexedwards/scs/badgerstore"
	"github.com/alexedwards/scs/postgresstore"
	"github.com/alexedwards/scs/redisstore"
	"github.com/alexedwards/scs/sqlite3store"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger/v3"
	"github.com/gomodule/redigo/redis"
)

type Session struct {
	CookieName     string
	CookieLifeTime string
	CookiePersist  string
	CookieSecure   string
	CookieDomain   string
	SessionType    string
	DBPool         *sql.DB
	RedisPool      *redis.Pool
	BadgerConn     *badger.DB
}

// initSession initializes the session manager
func (s *Session) InitSession() *scs.SessionManager {
	var persist, secure bool

	// cookie lifetime
	minutes, err := strconv.Atoi(s.CookieLifeTime)
	if err != nil {
		minutes = 60
	}

	// cookie persist
	if strings.ToLower(s.CookiePersist) == "true" {
		persist = true
	}

	// cookie secure
	if strings.ToLower(s.CookieSecure) == "true" {
		secure = true
	}

	sessionManager := scs.New()
	sessionManager.Cookie.Name = s.CookieName
	sessionManager.Lifetime = time.Duration(minutes) * time.Minute
	sessionManager.Cookie.Persist = persist
	sessionManager.Cookie.Secure = secure
	sessionManager.Cookie.Domain = s.CookieDomain

	switch strings.ToLower(s.SessionType) {
	case "postgres", "postgresql":
		sessionManager.Store = postgresstore.New(s.DBPool)
	case "sqlite":
		sessionManager.Store = sqlite3store.New(s.DBPool)
	case "redis":
		sessionManager.Store = redisstore.New(s.RedisPool)
	case "badger":
		sessionManager.Store = badgerstore.New(s.BadgerConn)
	default:
		// cookie
	}

	return sessionManager
}
