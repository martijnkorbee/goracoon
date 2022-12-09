# GoRacoon 

Laravel like framework. Special thanks to Trevor Sawler.

***

## Building the CLI

Clone the project

```bash
  git clone https://github.com/martijnkorbee/goracoon.git
```

Go to the project directory

```bash
  cd goracoon
```

Build CLI

```bash
  make build
```

### CLI commands:

```bash
  goracoon [command] [args]
```

| command           | args          | description                                                                   |
| :-----------------| :-------------| :-----------------------------------------------------------------------------|
| `help`            |               | show help text                                                                |
| `version`         |               | show version                                                                  |
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
