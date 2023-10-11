package main

import (
	"flag"
	"fmt"
	"go.uber.org/zap"
	"log"
	"os"
	"pismo-ledger-service/config"

	_ "github.com/lib/pq"
	"github.com/pressly/goose"
)

const (
	usage = `
Commands:
    up                   Migrate the DB to the most recent version available
    up-by-one            Migrate the DB up by 1
    up-to VERSION        Migrate the DB to a specific VERSION
    down                 Roll back the version by 1
    down-to VERSION      Roll back to a specific VERSION
    redo                 Re-run the latest migration
    reset                Roll back all migrations
    status               Dump the migration status for the current DB
    version              Print the current version of the database
    create NAME [sql|go] Creates new migration file with the current timestamp
    fix                  Apply sequential ordering to migrations
	`
)

var (
	// we are printing usage, so we have continue on error and exit with our own
	// status.
	flags = flag.NewFlagSet("migrate", flag.ContinueOnError)
	dir   = flags.String(
		"dir", "./migrations",
		"directory with migration files")
	configFile = flags.String("cfg", "config.yml", "config file path")
)

func main() {
	err := flags.Parse(os.Args[1:])
	if err != nil {
		log.Fatalf("migrate: unable to parse flags: %v\n", err)
	}

	args := flags.Args()
	if len(args) < 1 {
		flags.Usage()
		fmt.Println(usage)

		os.Exit(1)
	}

	configPath := config.GetConfigPath(os.Getenv("config"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatal("Loading config:", zap.Error(err))
	}

	// use sequential instead of timestamp for migrations
	goose.SetSequential(true)

	dbURI := fmt.Sprintf(
		`postgres://%v:%v@%v:%v/%v?sslmode=%v`,
		cfg.Postgres.PostgresqlUser,
		cfg.Postgres.PostgresqlPassword,
		cfg.Postgres.PostgresqlHost,
		cfg.Postgres.PostgresqlPort,
		cfg.Postgres.PostgresqlDbname,
		cfg.Postgres.PostgresqlSSLMode,
	)

	log.Printf("using connection uri: %s\n", dbURI)
	db, err := goose.OpenDBWithDriver("postgres", dbURI)
	if err != nil {
		log.Printf("migrate: failed to open db connection: %v\n", err)
		return
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Printf("migrate: failed to close db connection: %v\n", err)
		}
	}()

	// command to execute, e.g. up, down, etc.
	command := args[0]

	// arguments that command may take, e.g. create NAME sql
	arguments := []string{}
	if len(args) > 1 {
		arguments = append(arguments, args[1:]...)
	}

	if err := goose.Run(command, db, *dir, arguments...); err != nil {
		log.Printf("migrate %v: %v", command, err)
	}
}
