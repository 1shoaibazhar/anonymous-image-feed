package api

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"imagefeed/internal/repository"
)

type ImageHandler struct {
	repo *repository.ImageRepository
}

func NewImageHandler(repo *repository.ImageRepository) *ImageHandler {
	return &ImageHandler{repo: repo}
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

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *ImageHandler) ListTags(w http.ResponseWriter, r *http.Request) {
	tags, err := h.repo.ListTags(r.Context())
	if err != nil {
		http.Error(w, "failed to load tags", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tags)
}
