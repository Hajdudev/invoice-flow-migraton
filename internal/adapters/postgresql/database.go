package database

import (
	"context"
	"fmt"
	"io/fs"
	"log"
	"net/url"
	"strconv"
	"time"

	"github.com/Hajdudev/invoice-flow/internal/env"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Service interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection pool.
	Close() error

	// Migrate runs database migrations from a directory (not implemented here).
	Migrate(dir string) error
	MigrateFS(migrationsFS fs.FS, dir string) error
	MigrateDownFS(migrationsFS fs.FS, dir string) error
	MigrateDownAllFS(migrationsFS fs.FS, dir string) error

	// Pool returns the underlying *pgxpool.Pool.
	Pool() *pgxpool.Pool

	// GetDB returns a new *pgx.Conn. Caller is responsible for closing the connection.
	// Use this when you explicitly need a pgx.Conn for operations that require it.
	GetDB(ctx context.Context) (*pgx.Conn, error)
}

type service struct {
	pool    *pgxpool.Pool
	connStr string
}

var (
	database = env.GetString("DB_DATABASE", "invoice-flow")
	password = env.GetString("DB_PASSWORD", "pa55word")
	username = env.GetString("DB_USERNAME", "user")
	port     = env.GetString("DB_PORT", "5432")
	host     = env.GetString("DB_HOST", "localhost")
	schema   = env.GetString("DB_SCHEMA", "public")

	dbInstance *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Build DSN using net/url to ensure proper escaping
	u := &url.URL{
		Scheme: "postgres",
		User:   url.UserPassword(username, password),
		Host:   fmt.Sprintf("%s:%s", host, port),
		Path:   database,
	}
	q := u.Query()
	q.Set("sslmode", "disable")
	q.Set("search_path", schema)
	u.RawQuery = q.Encode()

	connStr := u.String()

	cfg, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Fatalf("unable to parse pgxpool config: %v", err)
	}

	pool, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		log.Fatalf("unable to create pgx pool: %v", err)
	}

	dbInstance = &service{
		pool:    pool,
		connStr: connStr,
	}
	return dbInstance
}

// Pool returns the underlying pgxpool.Pool.
func (s *service) Pool() *pgxpool.Pool {
	return s.pool
}

// GetDB returns a new pgx.Conn created from the stored connection string.
// Caller is responsible for closing the connection: conn.Close(ctx)
func (s *service) GetDB(ctx context.Context) (*pgx.Conn, error) {
	return pgx.Connect(ctx, s.connStr)
}

func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)
	// ping by acquiring from pool and doing a simple query
	if err := s.pool.Ping(ctx); err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Printf("db down: %v", err)
		return stats
	}

	stats["status"] = "up"
	stats["message"] = "It's healthy"

	// pool.Stat provides some useful counters
	p := s.pool.Stat()
	stats["acquired_conns"] = strconv.FormatUint(uint64(p.AcquiredConns()), 10)

	return stats
}

// Close closes the pool.
func (s *service) Close() error {
	if s.pool != nil {
		s.pool.Close()
		log.Printf("Disconnected from database: %s", database)
	}
	return nil
}
