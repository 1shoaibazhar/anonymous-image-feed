package api

import (
	"encoding/json"
	"net/http"
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
	images, err := h.repo.ListImages(r.Context())
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
