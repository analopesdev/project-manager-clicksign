package controllers

import (
	"bytes"
	"encoding/json"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/analopesdev/project-manager-clicksign/api/internal/infra/config"
	"github.com/go-chi/chi/v5"
)

var tinyPNG = []byte{
	0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a, 0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
	0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01, 0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
	0xde, 0x00, 0x00, 0x00, 0x0c, 0x49, 0x44, 0x41, 0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
	0x00, 0x03, 0x01, 0x01, 0x00, 0x18, 0xdd, 0x8d, 0xb0, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
	0x44, 0xae, 0x42, 0x60, 0x82,
}

func newUploadRouter(t *testing.T) (*chi.Mux, string) {
	t.Helper()
	dir := filepath.Join(t.TempDir(), "uploads")
	c, err := NewUploadController(&config.Config{UploadDir: dir}, slog.New(slog.NewTextHandler(io.Discard, nil)))
	if err != nil {
		t.Fatal(err)
	}
	r := chi.NewRouter()
	r.Post("/uploads", c.Upload)
	r.Get("/uploads/{name}", c.Serve)
	return r, dir
}

func postFile(t *testing.T, r http.Handler, field string, content []byte) *httptest.ResponseRecorder {
	t.Helper()
	var body bytes.Buffer
	w := multipart.NewWriter(&body)
	part, _ := w.CreateFormFile(field, "cover.png")
	part.Write(content)
	w.Close()

	req := httptest.NewRequest(http.MethodPost, "http://api.test/uploads", &body)
	req.Header.Set("Content-Type", w.FormDataContentType())
	rec := httptest.NewRecorder()
	r.ServeHTTP(rec, req)
	return rec
}

func TestUploadStoresImageAndServesIt(t *testing.T) {
	r, dir := newUploadRouter(t)

	rec := postFile(t, r, "file", tinyPNG)
	if rec.Code != http.StatusCreated {
		t.Fatalf("status = %d, body = %s", rec.Code, rec.Body)
	}
	var res uploadResponse
	json.Unmarshal(rec.Body.Bytes(), &res)
	if !strings.HasPrefix(res.URL, "http://api.test/uploads/") || !strings.HasSuffix(res.URL, ".png") {
		t.Fatalf("unexpected url %q", res.URL)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 1 {
		t.Fatalf("expected 1 stored file, got %d", len(entries))
	}

	get := httptest.NewRecorder()
	r.ServeHTTP(get, httptest.NewRequest(http.MethodGet, strings.TrimPrefix(res.URL, "http://api.test"), nil))
	if get.Code != http.StatusOK || !bytes.Equal(get.Body.Bytes(), tinyPNG) {
		t.Fatalf("serve: status %d, body differs", get.Code)
	}
	if got := get.Header().Get("X-Content-Type-Options"); got != "nosniff" {
		t.Fatalf("missing nosniff, got %q", got)
	}
}

func TestUploadRejectsNonImageEvenWithPngName(t *testing.T) {
	r, dir := newUploadRouter(t)
	rec := postFile(t, r, "file", []byte("<script>alert(1)</script>"))
	if rec.Code != http.StatusUnprocessableEntity {
		t.Fatalf("status = %d", rec.Code)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatal("rejected file must not be stored")
	}
}

func TestUploadRejectsMissingField(t *testing.T) {
	r, _ := newUploadRouter(t)
	if rec := postFile(t, r, "other", tinyPNG); rec.Code != http.StatusBadRequest {
		t.Fatalf("status = %d", rec.Code)
	}
}

func TestUploadRejectsFilesOver5MiB(t *testing.T) {
	r, dir := newUploadRouter(t)
	big := append(append([]byte{}, tinyPNG...), make([]byte, 6<<20)...)
	if rec := postFile(t, r, "file", big); rec.Code != http.StatusRequestEntityTooLarge {
		t.Fatalf("status = %d", rec.Code)
	}
	if entries, _ := os.ReadDir(dir); len(entries) != 0 {
		t.Fatal("oversized file must not be stored")
	}
}

func TestServeOnlyAcceptsGeneratedNames(t *testing.T) {
	r, dir := newUploadRouter(t)
	os.WriteFile(filepath.Join(dir, "secret.txt"), []byte("x"), 0o644)

	for _, name := range []string{"secret.txt", "..%2Fgo.mod", "abc.png", "../uploads"} {
		rec := httptest.NewRecorder()
		r.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/uploads/"+name, nil))
		if rec.Code != http.StatusNotFound {
			t.Errorf("GET /uploads/%s = %d, want 404", name, rec.Code)
		}
	}
}
