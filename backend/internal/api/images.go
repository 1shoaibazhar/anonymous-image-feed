package api

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"

	"imagefeed/internal/repository"
)

const cacheTTL = 30 * time.Second

type ImageHandler struct {
	repo *repository.ImageRepository
	rdb  *redis.Client
}

func NewImageHandler(repo *repository.ImageRepository, rdb *redis.Client) *ImageHandler {
	return &ImageHandler{repo: repo, rdb: rdb}
}

type imageResponse struct {
	ID        string   `json:"id"`
	Title     string   `json:"title"`
	Tags      []string `json:"tags"`
	URL       string   `json:"url"`
	CreatedAt string   `json:"created_at"`
}

type imagesListResponse struct {
	Images     []imageResponse `json:"images"`
	NextCursor *string         `json:"next_cursor"`
}

func encodeCursor(c repository.Cursor) string {
	raw := c.CreatedAt.Format(time.RFC3339Nano) + "|" + c.ID
	return base64.URLEncoding.EncodeToString([]byte(raw))
}

func decodeCursor(s string) (*repository.Cursor, error) {
	raw, err := base64.URLEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	parts := strings.SplitN(string(raw), "|", 2)
	if len(parts) != 2 {
		return nil, fmt.Errorf("malformed cursor")
	}
	createdAt, err := time.Parse(time.RFC3339Nano, parts[0])
	if err != nil {
		return nil, err
	}
	return &repository.Cursor{CreatedAt: createdAt, ID: parts[1]}, nil
}

func (h *ImageHandler) List(w http.ResponseWriter, r *http.Request) {
	tags := []string{}
	if raw := r.URL.Query().Get("tags"); raw != "" {
		for _, t := range strings.Split(raw, ",") {
			if trimmed := strings.TrimSpace(t); trimmed != "" {
				tags = append(tags, trimmed)
			}
		}
	}
	sort.Strings(tags)

	cursorParam := r.URL.Query().Get("cursor")
	var cursor *repository.Cursor
	if cursorParam != "" {
		decoded, err := decodeCursor(cursorParam)
		if err != nil {
			http.Error(w, "invalid cursor", http.StatusBadRequest)
			return
		}
		cursor = decoded
	}

	cacheKey := "images:" + strings.Join(tags, ",") + ":" + cursorParam
	if cached, err := h.rdb.Get(r.Context(), cacheKey).Result(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	images, hasMore, err := h.repo.ListImages(r.Context(), tags, cursor)
	if err != nil {
		http.Error(w, "failed to load images", http.StatusInternalServerError)
		return
	}

	resp := imagesListResponse{Images: make([]imageResponse, 0, len(images))}
	for _, img := range images {
		resp.Images = append(resp.Images, imageResponse{
			ID:        img.ID,
			Title:     img.Title,
			Tags:      img.Tags,
			URL:       img.FilePath,
			CreatedAt: img.CreatedAt.Format(time.RFC3339),
		})
	}
	if hasMore {
		last := images[len(images)-1]
		next := encodeCursor(repository.Cursor{CreatedAt: last.CreatedAt, ID: last.ID})
		resp.NextCursor = &next
	}

	body, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, "failed to load images", http.StatusInternalServerError)
		return
	}
	h.rdb.Set(r.Context(), cacheKey, body, cacheTTL)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}

func (h *ImageHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	cacheKey := "tags:all"
	if cached, err := h.rdb.Get(r.Context(), cacheKey).Result(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	tags, err := h.repo.ListTags(r.Context())
	if err != nil {
		http.Error(w, "failed to load tags", http.StatusInternalServerError)
		return
	}

	body, err := json.Marshal(tags)
	if err != nil {
		http.Error(w, "failed to load tags", http.StatusInternalServerError)
		return
	}
	h.rdb.Set(r.Context(), cacheKey, body, cacheTTL)

	w.Header().Set("Content-Type", "application/json")
	w.Write(body)
}
