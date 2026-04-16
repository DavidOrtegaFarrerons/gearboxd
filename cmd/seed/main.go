package main

import (
	"context"
	"database/sql"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"time"

	_ "github.com/lib/pq"
)

//go:embed seeds/*.sql
var seedFiles embed.FS

type config struct {
	dsn string
}

func main() {
	var cfg config

	flag.StringVar(&cfg.dsn, "dsn", os.Getenv("GEARBOXD_DB_DSN"), "DSN to connect to the database")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	files, err := seedFiles.ReadDir("seeds")
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	for _, f := range files {
		if f.IsDir() {
			continue
		}

		if filepath.Ext(f.Name()) != ".sql" {
			continue
		}

		query, err := seedFiles.ReadFile("seeds/" + f.Name())
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}

		logger.Info(fmt.Sprintf("Executing seeder from: %s", f.Name()))

		_, err = db.Exec(string(query))
		if err != nil {
			logger.Error(err.Error())
			os.Exit(1)
		}
	}

	logger.Info("The seeder was executed successfully")
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.dsn)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
