package goracoon

import (
	"github.com/upper/db/v4"
)

type initPaths struct {
	rootPath    string
	folderNames []string
}

type Database struct {
	DatabaseType   string
	ConnectionPool db.Session
}

type databaseConfig struct {
	host     string
	port     string
	user     string
	password string
	name     string
	sslMode  string
}

type redisConfig struct {
	host     string
	password string
	prefix   string
}

type cookieConfig struct {
	name     string
	lifeTime string
	persist  string
	secure   string
	domain   string
}
