package main

func doMigrate(arg2, arg3 string) error {
	dsn := getDSN()

	// run the migration command
	switch arg2 {
	case "up":
		err := racoon.MigrateUp(dsn)
		if err != nil {
			return err
		}

	case "down":
		if arg3 != "" {
			switch arg3 {
			case "all":
				err := racoon.MigrateDownAll(dsn)
				if err != nil {
					return err
				}
			case "force":
				err := racoon.MigrateForce(dsn)
				if err != nil {
					return err
				}
			}
		} else {
			err := racoon.Steps(-1, dsn)
			if err != nil {
				return err
			}
		}

	case "reset":
		err := racoon.MigrateDownAll(dsn)
		if err != nil {
			return err
		}

		err = racoon.MigrateUp(dsn)
		if err != nil {
			return err
		}

	default:
		showHelp()
	}

	return nil
}
