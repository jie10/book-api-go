package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jie10/book-api-go/internal/config"
	"github.com/jie10/book-api-go/internal/logger"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBPool is the interface for database operations
type DBPool interface {
	Acquire(ctx context.Context) (*pgxpool.Conn, error)
	Exec(ctx context.Context, sql string, arguments ...interface{}) (int64, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row
	Close()
}

type pgxPool struct {
	*pgxpool.Pool
	once sync.Once
}

var (
	instance *pgxPool
	once     sync.Once
)

// NewDBPool creates a new database connection pool (singleton pattern)
func NewDBPool(ctx context.Context, cfg *config.Config) (DBPool, error) {
	var initErr error
	once.Do(func() {
		connString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			cfg.Database.User,
			cfg.Database.Password,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
			cfg.Database.SSLMode,
		)

		poolConfig, err := pgxpool.ParseConfig(connString)
		if err != nil {
			initErr = err
			return
		}

		// Performance optimizations
		poolConfig.MaxConns = 50                      // Maximum number of connections
		poolConfig.MinConns = 10                      // Minimum number of connections
		poolConfig.MaxConnLifetime = time.Hour        // Maximum lifetime of a connection
		poolConfig.MaxConnIdleTime = time.Minute * 30 // Maximum idle time of a connection
		poolConfig.HealthCheckPeriod = time.Minute    // Health check frequency
		poolConfig.ConnConfig.RuntimeParams = map[string]string{
			"standard_conforming_strings": "on",
			"timezone":                    "UTC",
		}

		pool, err := pgxpool.NewWithConfig(ctx, poolConfig)
		if err != nil {
			initErr = err
			return
		}

		// Verify connection immediately
		if err := pool.Ping(ctx); err != nil {
			initErr = err
			return
		}

		instance = &pgxPool{Pool: pool}

		stats := pool.Stat()
		logger.Info("TotalConns: %d, IdleConns: %d\n", stats.TotalConns(), stats.IdleConns())
	})

	if initErr != nil {
		return nil, initErr
	}

	logger.Info("Successfully connected to database")

	return instance, nil
}

// Exec executes a query without returning any rows
func (p *pgxPool) Exec(ctx context.Context, sql string, arguments ...interface{}) (int64, error) {
	ct, err := p.Pool.Exec(ctx, sql, arguments...)
	return ct.RowsAffected(), err
}

// Close closes all connections in the pool
func (p *pgxPool) Close() {
	p.once.Do(func() {
		p.Pool.Close()
	})
}
