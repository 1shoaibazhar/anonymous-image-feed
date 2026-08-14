package api

import (
	"encoding/json"
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

	cacheKey := "images:" + strings.Join(tags, ",")
	if cached, err := h.rdb.Get(r.Context(), cacheKey).Result(); err == nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte(cached))
		return
	}

	images, err := h.repo.ListImages(r.Context(), tags)
	if err != nil {
		http.Error(w, "failed to load images", http.StatusInternalServerError)
		return
	}

	resp := make([]imageResponse, 0, len(images))
	for _, img := range images {
		resp = append(resp, imageResponse{
			ID:        img.ID,
			Title:     img.Title,
			Tags:      img.Tags,
			URL:       img.FilePath,
			CreatedAt: img.CreatedAt.Format(time.RFC3339),
		})
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
