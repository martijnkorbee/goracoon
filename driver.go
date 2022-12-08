package goracoon

import (
	"fmt"

	"github.com/upper/db/v4"
	pg "github.com/upper/db/v4/adapter/postgresql"
	sqlt "github.com/upper/db/v4/adapter/sqlite"
)

func (gr *Goracoon) OpenDB(dbType string, settings db.ConnectionURL) (db db.Session, err error) {
	switch dbType {
	case "postgres", "postgresql":
		db, err = pg.Open(settings)
		if err != nil {
			return nil, err
		}

		err = db.Ping()
		if err != nil {
			return nil, err
		}

	case "sqlite":
		err = gr.CreateDirIfNotExists("./tmp/db-data")
		if err != nil {
			return nil, err
		}

		err = gr.CreateFileIfNotExists("./tmp/db-data/" + gr.config.dbConfig.name + ".db")
		if err != nil {
			return nil, err
		}

		db, err = sqlt.Open(settings)
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func (gr *Goracoon) BuildDSN(dbType string) (settings db.ConnectionURL) {
	switch dbType {
	case "postgres", "postgresql":
		settings = pg.ConnectionURL{
			User:     gr.config.dbConfig.user,
			Password: gr.config.dbConfig.password,
			Host:     gr.config.dbConfig.host,
			Database: gr.config.dbConfig.name,
			Options: map[string]string{
				"sslmode": gr.config.dbConfig.sslMode,
			},
		}
	case "sqlite":
		settings = sqlt.ConnectionURL{
			Database: fmt.Sprintf("./tmp/db-data/%s.db", gr.config.dbConfig.name),
		}
	default:
		// do nothing
	}

	return settings
}
