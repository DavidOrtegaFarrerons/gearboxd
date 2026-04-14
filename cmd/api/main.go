package main

import (
	"context"
	"database/sql"
	"errors"
	"flag"
	"gearboxd/internal/data"
	"gearboxd/internal/mailer"
	"log"
	"log/slog"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
)

type config struct {
	port int
	env  string
	db   struct {
		dsn          string
		maxOpenConns int
		maxIdleConns int
		maxIdleTime  time.Duration
	}
	smtp struct {
		host     string
		port     int
		username string
		password string
		sender   string
	}
	limiter struct {
		rps     int64
		burst   int
		enabled bool
	}
	cors struct {
		trustedOrigins []string
	}
}
type application struct {
	config config
	logger *slog.Logger
	models data.Models
	mailer mailer.Mailer
}

func main() {
	var cfg config
	serverPort := 4000
	if v := os.Getenv("SERVER_PORT"); v != "" {
		sp, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid SMTP_PORT: %v", err)
		}

		serverPort = sp
	}

	flag.IntVar(&cfg.port, "port", serverPort, "Define the port used in the app")
	flag.StringVar(&cfg.env, "env", "dev", "Define the env used in the app")

	flag.StringVar(&cfg.db.dsn, "db-dsn", os.Getenv("GEARBOXD_DB_DSN"), "Define the PostgreSQL DSN")
	flag.IntVar(&cfg.db.maxOpenConns, "db-max-open-conns", 25, "PostgreSQL max open connections")
	flag.IntVar(&cfg.db.maxIdleConns, "db-max-idle-conns", 25, "PostgreSQL max idle connections")
	flag.DurationVar(&cfg.db.maxIdleTime, "db-max-idle-time", 15*time.Minute, "PostgreSQL max connection idle time")

	flag.StringVar(&cfg.smtp.host, "smtp-host", os.Getenv("SMTP_HOST"), "SMTP host")

	smtpPort := 25
	if v := os.Getenv("SMTP_PORT"); v != "" {
		port, err := strconv.Atoi(v)
		if err != nil {
			log.Fatalf("invalid SMTP_PORT: %v", err)
		}

		smtpPort = port
	}

	flag.IntVar(&cfg.smtp.port, "smtp-port", smtpPort, "SMTP port")
	flag.StringVar(&cfg.smtp.username, "smtp-username", os.Getenv("SMTP_USERNAME"), "SMTP username")
	flag.StringVar(&cfg.smtp.password, "smtp-password", os.Getenv("SMTP_PASSWORD"), "SMTP password")
	flag.StringVar(&cfg.smtp.sender, "smtp-sender", os.Getenv("SMTP_SENDER"), "SMTP sender")

	flag.Int64Var(&cfg.limiter.rps, "limiter-rps", 2, "Rate limiter maximum requests per second")
	flag.IntVar(&cfg.limiter.burst, "limiter-burst", 4, "Rate limiter maximum burst")
	flag.BoolVar(&cfg.limiter.enabled, "limiter-enabled", true, "Enable rate limiter")

	flag.Func("cors-trusted-origins", "Trusted CORS origins (space separated)", func(val string) error {
		cfg.cors.trustedOrigins = strings.Fields(val)
		return nil
	})

	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	db, err := openDB(cfg)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	defer db.Close()

	logger.Info("database connection pool has been established successfully")

	migrationDriver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	migrator, err := migrate.NewWithDatabaseInstance("file:///app/migrations", "postgres", migrationDriver)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	err = migrator.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		logger.Error(err.Error())
		os.Exit(1)
	}

	logger.Info("database migrations have been applied")

	models := data.NewModels(db)

	app := application{
		config: cfg,
		logger: logger,
		models: models,
		mailer: mailer.New(cfg.smtp.host, cfg.smtp.port, cfg.smtp.username, cfg.smtp.password, cfg.smtp.sender),
	}

	err = app.serve()
	if err != nil {
		app.logger.Error(err.Error())
		os.Exit(1)
	}
}

func openDB(cfg config) (*sql.DB, error) {
	db, err := sql.Open("postgres", cfg.db.dsn)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(cfg.db.maxOpenConns)
	db.SetMaxIdleConns(cfg.db.maxIdleConns)
	db.SetConnMaxIdleTime(cfg.db.maxIdleTime)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	err = db.PingContext(ctx)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}
