package controllers

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/config"
	"github.com/go-chi/chi/v5"
)

const (
	maxUploadBytes = 5 << 20

	maxUploadRequestBytes = maxUploadBytes + 1<<20
)

var allowedImages = map[string]string{
	"image/jpeg": ".jpg",
	"image/png":  ".png",
}

var storedFileName = regexp.MustCompile(`^[a-f0-9]{32}\.(jpg|png)$`)

type UploadController struct {
	dir    string
	logger *slog.Logger
}

func NewUploadController(cfg *config.Config, logger *slog.Logger) (*UploadController, error) {
	if err := os.MkdirAll(cfg.UploadDir, 0o755); err != nil {
		return nil, fmt.Errorf("create upload dir: %w", err)
	}
	return &UploadController{dir: cfg.UploadDir, logger: logger}, nil
}

type uploadResponse struct {
	URL string `json:"url"`
}

func (c *UploadController) Upload(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxUploadRequestBytes)

	file, header, err := r.FormFile("file")
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeError(w, http.StatusRequestEntityTooLarge, errors.New("file must be at most 5 MiB"))
			return
		}
		writeError(w, http.StatusBadRequest, errors.New(`multipart field "file" is required`))
		return
	}
	defer file.Close()

	if header.Size > maxUploadBytes {
		writeError(w, http.StatusRequestEntityTooLarge, errors.New("file must be at most 5 MiB"))
		return
	}

	sniff := make([]byte, 512)
	n, err := io.ReadFull(file, sniff)
	if err != nil && !errors.Is(err, io.ErrUnexpectedEOF) && !errors.Is(err, io.EOF) {
		writeError(w, http.StatusBadRequest, errors.New("could not read file"))
		return
	}
	ext, ok := allowedImages[http.DetectContentType(sniff[:n])]
	if !ok {
		writeError(w, http.StatusUnprocessableEntity, errors.New("file must be a .jpg or .png image"))
		return
	}

	token := make([]byte, 16)
	if _, err := rand.Read(token); err != nil {
		writeError(w, http.StatusInternalServerError, errors.New("could not store file"))
		return
	}
	name := hex.EncodeToString(token) + ext

	dst, err := os.OpenFile(filepath.Join(c.dir, name), os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0o644)
	if err != nil {
		c.logger.ErrorContext(r.Context(), "failed to create upload", "error", err)
		writeError(w, http.StatusInternalServerError, errors.New("could not store file"))
		return
	}
	defer dst.Close()

	if _, err := dst.Write(sniff[:n]); err == nil {
		_, err = io.Copy(dst, file)
	}
	if err != nil {
		os.Remove(dst.Name())
		c.logger.ErrorContext(r.Context(), "failed to write upload", "error", err)
		writeError(w, http.StatusInternalServerError, errors.New("could not store file"))
		return
	}

	writeJSON(w, http.StatusCreated, uploadResponse{URL: publicBaseURL(r) + "/uploads/" + name})
}

func (c *UploadController) Serve(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if !storedFileName.MatchString(name) {
		http.NotFound(w, r)
		return
	}
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("Cache-Control", "public, max-age=31536000, immutable")
	http.ServeFile(w, r, filepath.Join(c.dir, name))
}

func publicBaseURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https" {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
