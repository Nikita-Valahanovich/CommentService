package storage

import (
	"CommentService/pkg/models"
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type DB struct {
	pool *pgxpool.Pool
}

func NewDB() (*DB, error) {
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:admin@192.168.1.165/newsdb"
	}
	pool, err := pgxpool.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}
	return &DB{pool: pool}, nil
}

func (db *DB) Close() {
	db.pool.Close()
}

func (db *DB) InsertComment(c *models.Comment) error {
	query := `INSERT INTO comments (news_id, parent_id, content, author, created_at)
	          VALUES ($1, $2, $3, $4, $5)`
	_, err := db.pool.Exec(context.Background(), query,
		c.NewsID, c.ParentID, c.Content, c.Author, c.CreatedAt)
	return err
}

func (db *DB) GetCommentsByNewsID(newsID int) ([]models.Comment, error) {
	query := `SELECT id, news_id, parent_id, content, author, created_at
	          FROM comments WHERE news_id = $1`
	rows, err := db.pool.Query(context.Background(), query, newsID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []models.Comment
	for rows.Next() {
		var c models.Comment
		err := rows.Scan(&c.ID, &c.NewsID, &c.ParentID, &c.Content, &c.Author, &c.CreatedAt)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}
