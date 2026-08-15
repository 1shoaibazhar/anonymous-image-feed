package api

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"imagefeed/internal/imaging"
	"imagefeed/internal/repository"
	"imagefeed/internal/ws"
)

const (
	maxUploadSize = 5 << 20 // 5 MB
	uploadsDir    = "./uploads"
)

var allowedMimeTypes = map[string]bool{
	"image/jpeg": true,
	"image/png":  true,
	"image/webp": true,
}

type UploadHandler struct {
	repo *repository.ImageRepository
	hub  *ws.Hub
	rdb  *redis.Client
}

func NewUploadHandler(repo *repository.ImageRepository, hub *ws.Hub, rdb *redis.Client) *UploadHandler {
	return &UploadHandler{repo: repo, hub: hub, rdb: rdb}
}

func (h *UploadHandler) Create(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadSize)

	if err := r.ParseMultipartForm(maxUploadSize); err != nil {
		http.Error(w, "file too large or invalid form", http.StatusBadRequest)
		return
	}

	title := strings.TrimSpace(r.FormValue("title"))
	if title == "" {
		http.Error(w, "title is required", http.StatusBadRequest)
		return
	}
	tags := parseTags(r.FormValue("tags"))

	file, _, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "file is required", http.StatusBadRequest)
		return
	}
	defer file.Close()

	buf := make([]byte, 512)
	n, err := file.Read(buf)
	if err != nil && err != io.EOF {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}
	mimeType := http.DetectContentType(buf[:n])

	if !allowedMimeTypes[mimeType] {
		http.Error(w, "unsupported file type, only jpeg, png and webp are allowed", http.StatusBadRequest)
		return
	}

	if _, err := file.Seek(0, io.SeekStart); err != nil {
		http.Error(w, "failed to read file", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(uploadsDir, 0o755); err != nil {
		http.Error(w, "failed to prepare storage", http.StatusInternalServerError)
		return
	}

	normalized, err := imaging.Normalize(file)
	if err != nil {
		http.Error(w, "failed to process image", http.StatusInternalServerError)
		return
	}

	id := uuid.New().String()
	filename := id + ".jpg"
	err = os.WriteFile(filepath.Join(uploadsDir, filename), normalized, 0o644)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	relativePath := "/uploads/" + filename
	createdAt, err := h.repo.CreateImage(r.Context(), repository.NewImage{
		ID:        id,
		Title:     title,
		FilePath:  relativePath,
		MimeType:  "image/jpeg",
		SizeBytes: int64(len(normalized)),
		Tags:      tags,
	})
	if err != nil {
		http.Error(w, "failed to save image record", http.StatusInternalServerError)
		return
	}

	h.rdb.Del(r.Context(), "images::", "tags:all")

	created := imageResponse{
		ID:        id,
		Title:     title,
		Tags:      tags,
		URL:       relativePath,
		CreatedAt: createdAt.Format(time.RFC3339),
	}

	if msg, err := json.Marshal(wsMessage{Type: "image.created", Payload: created}); err == nil {
		h.hub.Broadcast(msg)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(created)
}

func parseTags(raw string) []string {
	if raw == "" {
		return []string{}
	}
	parts := strings.Split(raw, ",")
	tags := make([]string, 0, len(parts))
	for _, p := range parts {
		if t := strings.TrimSpace(p); t != "" {
			tags = append(tags, t)
		}
	}
	return tags
}
