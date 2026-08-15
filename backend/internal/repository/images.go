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
	MimeType  string
	SizeBytes int64
	Tags      []string
}

type Cursor struct {
	CreatedAt time.Time
	ID        string
}

type ImageRepository struct {
	db *pgxpool.Pool
}

func NewImageRepository(db *pgxpool.Pool) *ImageRepository {
	return &ImageRepository{db: db}
}

const imagesPageSize = 50

const listImagesQuery = `
	SELECT i.id::text, i.title, i.file_path, i.created_at
	FROM images i
	WHERE i.status = 'ready'
	  AND ($1 = false OR i.created_at < $2 OR (i.created_at = $2 AND i.id::text < $3))
	ORDER BY i.created_at DESC, i.id::text DESC
	LIMIT $4
`

const listImagesByTagQuery = `
	SELECT i.id::text, i.title, i.file_path, i.created_at
	FROM images i
	WHERE i.status = 'ready'
	  AND EXISTS (
	    SELECT 1 FROM image_tags it
	    JOIN tags t ON t.id = it.tag_id
	    WHERE it.image_id = i.id AND t.name = ANY($1::text[])
	  )
	  AND ($2 = false OR i.created_at < $3 OR (i.created_at = $3 AND i.id::text < $4))
	ORDER BY i.created_at DESC, i.id::text DESC
	LIMIT $5
`

const tagsForImagesQuery = `
	SELECT it.image_id::text, t.name
	FROM tags t
	JOIN image_tags it ON it.tag_id = t.id
	WHERE it.image_id = ANY($1::uuid[])
	ORDER BY it.image_id, t.name
`

const listAllTagsQuery = `
	SELECT DISTINCT t.name
	FROM tags t
	JOIN image_tags it ON it.tag_id = t.id
	JOIN images i ON i.id = it.image_id
	WHERE i.status = 'ready'
	ORDER BY t.name
`

const insertImageQuery = `
	INSERT INTO images (id, title, file_path, mime_type, size_bytes, status)
	VALUES ($1, $2, $3, $4, $5, 'ready')
	RETURNING created_at
`

const upsertTagQuery = `
	INSERT INTO tags (name) VALUES ($1)
	ON CONFLICT (name) DO UPDATE SET name = EXCLUDED.name
	RETURNING id
`

const insertImageTagQuery = `
	INSERT INTO image_tags (image_id, tag_id) VALUES ($1, $2)
`

func (r *ImageRepository) ListImages(ctx context.Context, tags []string, after *Cursor) ([]Image, bool, error) {
	hasCursor := after != nil
	var afterCreatedAt time.Time
	var afterID string
	if after != nil {
		afterCreatedAt = after.CreatedAt
		afterID = after.ID
	}

	var images []Image
	var err error
	if len(tags) == 0 {
		images, err = r.queryImages(ctx, listImagesQuery, hasCursor, afterCreatedAt, afterID, imagesPageSize+1)
	} else {
		images, err = r.queryImages(ctx, listImagesByTagQuery, tags, hasCursor, afterCreatedAt, afterID, imagesPageSize+1)
	}
	if err != nil {
		return nil, false, err
	}

	images, hasMore := paginate(images, imagesPageSize)

	if len(images) > 0 {
		ids := make([]string, len(images))
		for i, img := range images {
			ids[i] = img.ID
		}
		tagsByImage, err := r.queryTagsForImages(ctx, ids)
		if err != nil {
			return nil, false, err
		}
		attachTags(images, tagsByImage)
	}

	return images, hasMore, nil
}

func paginate(images []Image, pageSize int) ([]Image, bool) {
	hasMore := len(images) > pageSize
	if hasMore {
		images = images[:pageSize]
	}
	return images, hasMore
}

func attachTags(images []Image, tagsByImage map[string][]string) {
	for i := range images {
		tags := tagsByImage[images[i].ID]
		if tags == nil {
			tags = []string{}
		}
		images[i].Tags = tags
	}
}

func (r *ImageRepository) queryImages(ctx context.Context, query string, args ...any) ([]Image, error) {
	rows, err := r.db.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var images []Image
	for rows.Next() {
		var img Image
		if err := rows.Scan(&img.ID, &img.Title, &img.FilePath, &img.CreatedAt); err != nil {
			return nil, err
		}
		images = append(images, img)
	}
	return images, rows.Err()
}

func (r *ImageRepository) queryTagsForImages(ctx context.Context, imageIDs []string) (map[string][]string, error) {
	rows, err := r.db.Query(ctx, tagsForImagesQuery, imageIDs)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]string, len(imageIDs))
	for rows.Next() {
		var imageID, name string
		if err := rows.Scan(&imageID, &name); err != nil {
			return nil, err
		}
		result[imageID] = append(result[imageID], name)
	}
	return result, rows.Err()
}

func (r *ImageRepository) ListTags(ctx context.Context) ([]string, error) {
	rows, err := r.db.Query(ctx, listAllTagsQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	tags := []string{}
	for rows.Next() {
		var tag string
		if err := rows.Scan(&tag); err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}
	return tags, rows.Err()
}

func (r *ImageRepository) CreateImage(ctx context.Context, img NewImage) (time.Time, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return time.Time{}, err
	}
	defer tx.Rollback(ctx)

	var createdAt time.Time
	err = tx.QueryRow(ctx, insertImageQuery,
		img.ID, img.Title, img.FilePath, img.MimeType, img.SizeBytes).Scan(&createdAt)
	if err != nil {
		return time.Time{}, err
	}

	for _, tagName := range img.Tags {
		var tagID int
		if err := tx.QueryRow(ctx, upsertTagQuery, tagName).Scan(&tagID); err != nil {
			return time.Time{}, err
		}

		if _, err := tx.Exec(ctx, insertImageTagQuery, img.ID, tagID); err != nil {
			return time.Time{}, err
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return time.Time{}, err
	}

	return createdAt, nil
}
