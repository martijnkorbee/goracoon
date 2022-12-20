# GoRacoon 

Laravel like framework with embeded database and caching.
Special thanks to [Trevor Sawler](https://github.com/tsawler "Trevor Sawler").

***

## Supported databases/caching
Using: [upper/db/v4](https://upper.io/v4/ "upper/db/v4")
* Postgres
* TODO: Adding built in support for MySQL/MariaDB and MSSQL
* SQLite3
* Redis
* BadgerDB

## Building the CLI

#### Clone the project
```bash
  git clone https://github.com/martijnkorbee/goracoon.git
```

#### Go to the project directory
```bash
  cd goracoon
```

#### Build CLI
Builds the CLI in tmp/
```bash
  make build
  ./tmp/goracoon [command] [args]
```

#### Install CLI
Installs the CLI in $GOPATH/bin. Make sure to export this path.
```bash
  make install
```

### CLI commands:
```bash
  goracoon [command] [args]
```

| command           | args          | description                                                                   |
| :-----------------| :-------------| :-----------------------------------------------------------------------------|
| `help`            |               | show help text                                                                |
| `version`         |               | show version                                                                  |
| `maintenance`     | `up/down`     | put the application in or out of maintenance mode                             |   
| `make new`        | `appname`     | creates a new skeleton app                                                    |
| `make migration`  | `name`        | creates 2 new up and down migrations                                          |
| `migrate`         |               | runs all non run up and down migrations                                       |
| `migrate`         | `down`        | reverse the most recent migration                                             |
| `migrate`         | `down force`  | force down migration                                                          |
| `migrate`         | `reset`       | runs all down migrations in reverse order then all up migrations              |
| `make auth`       |               | creates and runs migrations for auth tables and creates models and middleware |
| `make session`    |               | creates a table in the db to use as session store                             |
| `make handler`    | `name`        | creates a stub handler in the handlers dir                                    |
| `make model`      | `name`        | creates a new model in the data dir                                           |
| `make model`      | `--migrate`   | use the migrate flag to also create the db table                              |
| `make key`        |               | generates a 32 character key                                                  |
| `make mail`       | `name`        | make a new mail template                                                      |
