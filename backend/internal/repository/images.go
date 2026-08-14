package repository

import (
	"context"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Image struct {
	ID        string
	Title     string
	Tags      []string
	FilePath  string
	CreatedAt time.Time
}

type NewImage struct {
	ID        string
	Title     string
	FilePath  string
	ThumbPath string
	MimeType  string
	SizeBytes int64
	Tags      []string
}

type ImageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) ListImages(ctx context.Context) ([]Image, error) {
	rows, err := r.db.Query(ctx, `
		SELECT i.id::text, i.title, i.file_path, i.created_at,
		       COALESCE(array_agg(t.name ORDER BY t.name) FILTER (WHERE t.name IS NOT NULL), '{}') AS tags
		FROM images i
		LEFT JOIN image_tags it ON it.image_id = i.id
		LEFT JOIN tags t ON t.id = it.tag_id
		WHERE i.status = 'ready'
		GROUP BY i.id
		ORDER BY i.created_at DESC
		LIMIT 50
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.ID, &img.Title, &img.FilePath, &img.CreatedAt, &img.Tags); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *ImageRepository) CreateImage(ctx context.Context, img NewImage) (time.Time, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer tx.Rollback(ctx)

	var createdAt time.Time
	err = tx.QueryRow(ctx, `
		INSERT INTO images (id, title, file_path, thumb_path, mime_type, size_bytes, status)
		VALUES ($1, $2, $3, $4, $5, $6, 'ready')
		RETURNING created_at
	`, img.ID, img.Title, img.FilePath, img.ThumbPath, img.MimeType, img.SizeBytes).Scan(&createdAt)
	if err != nil {
		return time.Time{}, err
	}

	for _, tagName := range img.Tags {
		var tagID int
		err := tx.QueryRow(ctx, `
			INSERT INTO tags (name) VALUES ($1)
			ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
			RETURNING id
		`, tagName).Scan(&tagID)
		if err != nil {
			return time.Time{}, err
		}

		if _, err := tx.Exec(ctx, `
			INSERT INTO image_tags (image_id, tag_id) VALUES ($1, $2)
		`, img.ID, tagID); err != nil {
			return time.Time{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return time.Time{}, err
	}

	return createdAt, nil
}
