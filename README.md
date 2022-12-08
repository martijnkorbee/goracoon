# GoRacoon 

Laravel like framework. Special thanks to Trevor Sawler.

Available commands:
    help                    - show the help commands
    version                 - print application version
    make new <appname>      - creates a new skeleton app in current directory
    make migration <name>   - creates 2 new up and down migrations in the migrations folder
    migrate                 - runs all up migrations that have not been run previously
    migrate down            - reverse the most recent migration
                            - use migrate down force to force the migration
    migrate reset           - runs all down migrations in reverse order, and then all up migrations
    make auth               - creates and runs migrations for authentication tables and creates models and middleware
    make model <name>       - creates a new model in the data directory
                            - use the --migrate flag to also create the db table
    make handler <name>     - creates a stub handler in the handlers directory
    make session            - creates a table in the db as a session store
