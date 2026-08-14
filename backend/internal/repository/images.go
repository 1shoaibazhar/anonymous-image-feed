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
	CreatedAt time.Time
}

type ImageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{db: db}
}

func (r *ImageRepository) ListImages(ctx context.Context) ([]Image, error) {
	rows, err := r.db.Query(ctx, `
		SELECT i.id::text, i.title, i.created_at,
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
		if err := rows.Scan(&img.ID, &img.Title, &img.CreatedAt, &img.Tags); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}
