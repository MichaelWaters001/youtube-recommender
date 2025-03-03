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

// GetCreator fetches a creator by ID
func GetCreator(ctx context.Context, id string) (map[string]interface{}, error) {
	var creator struct {
		ID          int
		YouTubeID   string
		Name        string
		Description string
	}

	err := DB.QueryRow(ctx, "SELECT id, youtube_id, name, description FROM creators WHERE id = $1", id).
		Scan(&creator.ID, &creator.YouTubeID, &creator.Name, &creator.Description)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch creator: %w", err)
	}

	return map[string]interface{}{
		"id":          creator.ID,
		"youtube_id":  creator.YouTubeID,
		"name":        creator.Name,
		"description": creator.Description,
	}, nil
}

// SearchCreatorsByTag finds creators with a specific tag
func SearchCreatorsByTag(ctx context.Context, tag string) ([]map[string]interface{}, error) {
	rows, err := DB.Query(ctx, `
		SELECT c.id, c.youtube_id, c.name, c.description
		FROM creators c
		JOIN creator_tags ct ON c.id = ct.creator_id
		JOIN tags t ON ct.tag_id = t.id
		WHERE t.name ILIKE $1
	`, tag)
	if err != nil {
		return nil, fmt.Errorf("failed to search creators: %w", err)
	}
	defer rows.Close()

	var creators []map[string]interface{}
	for rows.Next() {
		var id int
		var youtubeID, name, description string
		if err := rows.Scan(&id, &youtubeID, &name, &description); err != nil {
			return nil, fmt.Errorf("failed to scan creator row: %w", err)
		}
		creators = append(creators, map[string]interface{}{
			"id":          id,
			"youtube_id":  youtubeID,
			"name":        name,
			"description": description,
		})
	}

	return creators, nil
}

// AddTag stores a tag for a creator in the database
func AddTag(ctx context.Context, creatorID int, tagName string, userID int) (int, error) {
	var tagID int
	err := DB.QueryRow(ctx, `
		INSERT INTO creator_tags (creator_id, tag_id, user_id)
		VALUES ($1, (SELECT id FROM tags WHERE name = $2), $3)
		RETURNING id
	`, creatorID, tagName, userID).Scan(&tagID)

	if err != nil {
		return 0, fmt.Errorf("failed to add tag: %w", err)
	}

	return tagID, nil
}

// VoteTag adds or updates a user's vote for a tag
func VoteTag(ctx context.Context, userID int, creatorTagID int, voteType int) error {
	_, err := DB.Exec(ctx, `
		INSERT INTO votes (user_id, creator_tag_id, vote_type)
		VALUES ($1, $2, $3)
		ON CONFLICT (user_id, creator_tag_id)
		DO UPDATE SET vote_type = EXCLUDED.vote_type
	`, userID, creatorTagID, voteType)

	if err != nil {
		return fmt.Errorf("failed to vote on tag: %w", err)
	}

	return nil
}

// AddCreator inserts a new creator into the database
func AddCreator(ctx context.Context, youtubeID, name, description string) (int, error) {
	var creatorID int
	err := DB.QueryRow(ctx, `
		INSERT INTO creators (youtube_id, name, description)
		VALUES ($1, $2, $3) RETURNING id
	`, youtubeID, name, description).Scan(&creatorID)

	if err != nil {
		return 0, fmt.Errorf("failed to add creator: %w", err)
	}

	return creatorID, nil
}

// GetTags retrieves all tags associated with a given creator
func GetTags(ctx context.Context, creatorID int) ([]map[string]interface{}, error) {
	rows, err := DB.Query(ctx, `
		SELECT t.id, t.name, ct.user_id
		FROM creator_tags ct
		JOIN tags t ON ct.tag_id = t.id
		WHERE ct.creator_id = $1
	`, creatorID)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch tags: %w", err)
	}
	defer rows.Close()

	var tags []map[string]interface{}
	for rows.Next() {
		var tagID int
		var tagName string
		var userID int
		if err := rows.Scan(&tagID, &tagName, &userID); err != nil {
			return nil, fmt.Errorf("failed to scan tag row: %w", err)
		}
		tags = append(tags, map[string]interface{}{
			"id":      tagID,
			"name":    tagName,
			"user_id": userID,
		})
	}

	return tags, nil
}

// RemoveVote deletes a user's vote for a specific tag
func RemoveVote(ctx context.Context, userID int, creatorTagID int) error {
	_, err := DB.Exec(ctx, "DELETE FROM votes WHERE user_id = $1 AND creator_tag_id = $2", userID, creatorTagID)
	if err != nil {
		return fmt.Errorf("failed to remove vote: %w", err)
	}

	return nil
}
