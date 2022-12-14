package goracoon

import (
	"database/sql"
	"fmt"
	stdlog "log"
	"os"
	"strconv"
	"time"

	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"github.com/dgraph-io/badger"
	"github.com/go-chi/chi/v5"
	"github.com/gomodule/redigo/redis"
	"github.com/joho/godotenv"
	"github.com/martijnkorbee/goracoon/cache"
	"github.com/martijnkorbee/goracoon/logger"
	"github.com/martijnkorbee/goracoon/mailer"
	"github.com/martijnkorbee/goracoon/render"
	"github.com/martijnkorbee/goracoon/session"
	"github.com/robfig/cron/v3"
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
)

const version = "1.0.0"

var logOut *os.File

var redisCache *cache.RedisCache
var redisPool *redis.Pool
var badgerCache *cache.BadgerCache
var badgerConnection *badger.DB

var maintenanceMode bool

// goracoon is the overall type for the goracoon package.
// Members that are exported in this type are available to any application that uses it.
type Goracoon struct {
	AppName        string
	Debug          bool
	Version        string
	config         config
	Log            *zerolog.Logger
	HTPPLog        *zerolog.Logger
	RootPath       string
	SessionManager *scs.SessionManager
	Routes         *chi.Mux
	Render         *render.Render
	JetViews       *jet.Set
	DB             Database
	EncryptionKey  string
	Cache          cache.Cache
	Scheduler      *cron.Cron
	Mail           mailer.Mail
}

// config used to extract configuration from .env to be used by application
type config struct {
	host        string
	port        string
	renderer    string
	sessionType string
	cacheType   string
	dbType      string
	dbConfig    databaseConfig
	redis       redisConfig
	cookie      cookieConfig
}

// New reads the .env file, creates the application config, populates the     goracoon app type with settings
// based on .env, and creates necessary folders and files if they don't exist
func (gr *Goracoon) New(rootPath string) error {
	pathConfig := initPaths{
		rootPath: rootPath,
		folderNames: []string{
			"handlers",
			"migrations",
			"views",
			"mail",
			"data",
			"public",
			"tmp",
			"logs",
			"middleware",
		},
	}

	// init racoon app
	err := gr.Init(pathConfig)
	if err != nil {
		return err
	}

	// load .env
	err = gr.checkDotEnv(rootPath)
	if err != nil {
		return err
	}
	err = godotenv.Load(rootPath + "/.env")
	if err != nil {
		return err
	}

	// get application config from .env
	gr.loadConfig()

	// assign application variables
	gr.Debug, _ = strconv.ParseBool(os.Getenv("DEBUG"))
	gr.Version = version
	gr.RootPath = rootPath

	// start loggers
	l := logger.Logger{
		Debug:         gr.Debug,
		ConsoleOutput: os.Getenv("LOG_CONSOLE"),
		FileOutput:    os.Getenv("LOG_FILE"),
		FileConfig: lumberjack.Logger{
			Filename:   gr.RootPath + "/logs/log.log",
			MaxBackups: 2,
			LocalTime:  true,
		},
	}
	gr.Log = l.StartLoggers()

	// overwrite default logger for other packages
	stdlog.SetOutput(gr.Log)

	httpl := logger.Logger{
		Debug:         gr.Debug,
		ConsoleOutput: os.Getenv("LOG_HTTP_CONSOLE"),
		FileOutput:    os.Getenv("LOG_HTTP_FILE"),
		FileConfig: lumberjack.Logger{
			Filename:   gr.RootPath + "/logs/httplog.log",
			MaxBackups: 2,
			LocalTime:  true,
		},
	}
	gr.HTPPLog = httpl.StartLoggers()

	// load encryption key
	gr.EncryptionKey = os.Getenv("KEY")

	// add scheduler
	gr.Scheduler = cron.New()

	// add mailer
	gr.Mail = gr.createMailer()

	// connect to application database
	if gr.config.dbType != "" {
		db, err := gr.OpenDB(gr.config.dbType, gr.BuildDSN(gr.config.dbType))
		if err != nil {
			gr.Log.Fatal().Err(err).Msgf("could not start db connection")
		} else {
			gr.Log.Info().Msg(fmt.Sprintf("connected to application database:%s", gr.config.dbType))
		}

		gr.DB = Database{
			DatabaseType:   gr.config.dbType,
			ConnectionPool: db,
		}
	}

	// connect to redis
	if gr.config.cacheType == "redis" || gr.config.sessionType == "redis" {
		redisCache = gr.createClientRedisCache()

		// check for connection
		ok, err := redisCache.Ping()
		if err != nil {
			gr.Log.Fatal().Err(err).Msgf("could not connect to redis")
		} else {
			gr.Log.Info().Msg(fmt.Sprintf("connected to redis, replied with:%s", ok))
		}

		// add cache client
		gr.Cache = redisCache
		redisPool = redisCache.Pool
	}

	// connect to badgerDB
	if gr.config.cacheType == "badger" || gr.config.sessionType == "badger" {
		badgerCache = gr.createClientBadgerCache()

		// add cache client
		gr.Cache = badgerCache
		badgerConnection = badgerCache.Connection

		// garbage collecting
		_, err = gr.Scheduler.AddFunc("@daily", func() {
			badgerCache.Connection.RunValueLogGC(0.7)
		})
		if err != nil {
			gr.Log.Error().Err(err).Msg("")
		}
	}

	// start mail channels
	go gr.Mail.ListenForMail()

	// add session (init session must be called before routes)
	session := session.Session{
		CookieName:     gr.config.cookie.name,
		CookieLifeTime: gr.config.cookie.lifeTime,
		CookiePersist:  gr.config.cookie.persist,
		CookieSecure:   gr.config.cookie.secure,
		CookieDomain:   gr.config.cookie.domain,
		SessionType:    gr.config.sessionType,
	}
	switch gr.config.sessionType {
	case "postgres", "postgresql", "sqlite":
		session.DBPool = gr.DB.ConnectionPool.Driver().(*sql.DB)
	case "redis":
		session.RedisPool = redisCache.Pool
	case "badger":
		session.BadgerConn = badgerCache.Connection
	default:
		// defaults to cookie
	}
	gr.SessionManager = session.InitSession()
	// add routes

	gr.Routes = gr.routes().(*chi.Mux)

	// Jet engine renderer
	if gr.Debug {
		gr.JetViews = jet.NewSet(
			jet.NewOSFileSystemLoader(rootPath+"/views"),
			jet.InDevelopmentMode(),
		)
	} else {
		gr.JetViews = jet.NewSet(
			jet.NewOSFileSystemLoader(rootPath + "/views"),
		)
	}
	gr.createRenderer()

	return nil
}

// Init creates necessary folders for the     goracoon application
func (gr *Goracoon) Init(p initPaths) error {
	root := p.rootPath

	// iterate through the required foldernames and create them if they don't exist
	for _, path := range p.folderNames {
		err := gr.CreateDirIfNotExists(root + "/" + path)
		if err != nil {
			return err
		}
	}

	return nil
}

func (gr *Goracoon) loadConfig() {
	gr.config = config{
		host:        os.Getenv("HOST"),
		port:        os.Getenv("PORT"),
		renderer:    os.Getenv("RENDERER"),
		sessionType: os.Getenv("SESSION_TYPE"),
		cacheType:   os.Getenv("CACHE"),
		dbType:      os.Getenv("DATABASE_TYPE"),
		dbConfig: databaseConfig{
			host:     os.Getenv("DATABASE_HOST"),
			port:     os.Getenv("DATABASE_PORT"),
			user:     os.Getenv("DATABASE_USER"),
			name:     os.Getenv("DATABASE_NAME"),
			password: os.Getenv("DATABASE_PASS"),
			sslMode:  os.Getenv("DATABASE_SSL_MODE"),
		},
		cookie: cookieConfig{
			name:     os.Getenv("COOKIE_NAME"),
			lifeTime: os.Getenv("COOKIE_LIFETIME"),
			persist:  os.Getenv("COOKIE_PERSIST"),
			secure:   os.Getenv("COOKIE_SECURE"),
			domain:   os.Getenv("COOKIE_DOMAIN"),
		},
		redis: redisConfig{
			host:     os.Getenv("REDIS_HOST"),
			password: os.Getenv("REDIS_PASSWORD"),
			prefix:   os.Getenv("REDIS_PREFIX"),
		},
	}
}

// checkDotEnv creates the .env file if doesn't exist
func (gr *Goracoon) checkDotEnv(path string) error {
	err := gr.CreateFileIfNotExists(path + "/.env")
	if err != nil {
		return err
	}

	return nil
}

// createRenderer
func (gr *Goracoon) createRenderer() {
	myRenderer := render.Render{
		Renderer: gr.config.renderer,
		RootPath: gr.RootPath,
		Port:     gr.config.port,
		JetViews: gr.JetViews,
		Session:  gr.SessionManager,
	}

	gr.Render = &myRenderer
}

func (gr *Goracoon) createMailer() mailer.Mail {
	port, _ := strconv.Atoi(os.Getenv("SMTP_PORT"))

	return mailer.Mail{
		Domain:      os.Getenv("MAILER_DOMAIN"),
		Templates:   gr.RootPath + "/mail",
		Host:        os.Getenv("SMTP_HOST"),
		Port:        port,
		Username:    os.Getenv("SMTP_USERNAME"),
		Password:    os.Getenv("SMTP_PASSWORD"),
		Encryption:  os.Getenv("SMTP_ENCRYPTION"),
		FromAddress: os.Getenv("SMTP_FROM_ADDRESS"),
		FromName:    os.Getenv("SMTP_FROM_NAME"),
		Jobs:        make(chan mailer.Message, 20),
		Results:     make(chan mailer.Result, 20),
		API:         os.Getenv("MAILER_API"),
		API_KEY:     os.Getenv("MAILER_KEY"),
		API_URL:     os.Getenv("MAILER_URL"),
	}
}

func (gr *Goracoon) createClientRedisCache() *cache.RedisCache {
	cacheClient := cache.RedisCache{
		Pool:   gr.createRedisPool(),
		Prefix: gr.config.redis.prefix,
	}
	return &cacheClient
}

// createRedisPool creates a redis connection pool
func (gr *Goracoon) createRedisPool() *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   10000,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp",
				gr.config.redis.host,
				redis.DialPassword(gr.config.redis.password))
		},

		TestOnBorrow: func(conn redis.Conn, t time.Time) error {
			_, err := conn.Do("PING")
			return err
		},
	}
}

func (gr *Goracoon) createClientBadgerCache() *cache.BadgerCache {
	cacheClient := cache.BadgerCache{
		Connection: gr.createBadgerConnection(),
	}
	return &cacheClient
}

func (gr *Goracoon) createBadgerConnection() *badger.DB {
	err := gr.CreateDirIfNotExists("./tmp/badger")
	if err != nil {
		gr.Log.Error().Err(err).Msg("")
	}

	db, err := badger.Open(badger.DefaultOptions("./tmp/badger"))
	if err != nil {
		return nil
	}
	return db
}
