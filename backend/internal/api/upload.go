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

	"imagefeed/internal/repository"
)

const (
	maxUploadSize = 10 << 20 // 10 MB
	uploadsDir    = "./uploads"
)

var allowedMimeTypes = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
	"image/webp": ".webp",
	"image/gif":  ".gif",
}

type UploadHandler struct {
	repo *repository.ImageRepository
}

func NewUploadHandler(repo *repository.ImageRepository) *UploadHandler {
	return &UploadHandler{repo: repo}
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

	ext, ok := allowedMimeTypes[mimeType]
	if !ok {
		http.Error(w, "unsupported file type, only jpeg, png, webp, and gif are allowed", http.StatusBadRequest)
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

	id := uuid.New().String()
	filename := id + ext
	dst, err := os.Create(filepath.Join(uploadsDir, filename))
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	size, err := io.Copy(dst, file)
	if err != nil {
		http.Error(w, "failed to save file", http.StatusInternalServerError)
		return
	}

	relativePath := "/uploads/" + filename
	createdAt, err := h.repo.CreateImage(r.Context(), repository.NewImage{
		ID:        id,
		Title:     title,
		FilePath:  relativePath,
		ThumbPath: relativePath, // no real thumbnail yet, same file for both
		MimeType:  mimeType,
		SizeBytes: size,
		Tags:      tags,
	})
	if err != nil {
		http.Error(w, "failed to save image record", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(imageResponse{
		ID:        id,
		Title:     title,
		Tags:      tags,
		CreatedAt: createdAt.Format(time.RFC3339),
	})
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
