package db

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/MichaelWaters001/youtube-recommender/pkg/config"
	"github.com/jackc/pgx"
	"github.com/jackc/pgx/v5/pgxpool"
)

var DB *pgxpool.Pool

// InitDB initializes the PostgreSQL connection pool
func InitDB() error {
	cfg := config.AppConfig.Database

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Dbname, cfg.Sslmode)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	pool, err := pgxpool.New(ctx, connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to database: %w", err)
	}

	// Ping to verify connection
	err = pool.Ping(ctx)
	if err != nil {
		return fmt.Errorf("database is unreachable: %w", err)
	}

	DB = pool
	log.Println("Connected to database")
	return nil
}

// EnsureUser checks if a user exists in the database and inserts if not
func EnsureUser(ctx context.Context, email string) (int, error) {
	var userID int

	err := DB.QueryRow(ctx, "SELECT id FROM users WHERE google_id = $1", email).Scan(&userID)
	if err != nil {
		// Handle case where the user does not exist
		if errors.Is(err, pgx.ErrNoRows) {
			// Insert new user
			err = DB.QueryRow(ctx, "INSERT INTO users (google_id) VALUES ($1) RETURNING id", email).Scan(&userID)
			if err != nil {
				return 0, fmt.Errorf("failed to insert user: %w", err)
			}
			return userID, nil
		}
		return 0, fmt.Errorf("database query error: %w", err)
	}

	return userID, nil
}
