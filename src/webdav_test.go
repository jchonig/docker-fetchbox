package main

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestWebDAVUpload(t *testing.T) {
	var (
		gotMethod string
		gotPath   string
		gotBody   []byte
		gotAuth   string
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotBody, _ = io.ReadAll(r.Body)
		gotAuth = r.Header.Get("Authorization")
		w.WriteHeader(http.StatusCreated)
	}))
	defer srv.Close()

	dest := Destination{
		Type:        "webdav",
		URL:         srv.URL,
		Path:        "/files/",
		Username:    "user",
		PasswordEnv: "",
	}

	u, err := newWebDAVUploader(dest)
	if err != nil {
		t.Fatalf("newWebDAVUploader: %v", err)
	}

	payload := []byte("hello webdav")
	if err := u.Upload("doc.txt", payload); err != nil {
		t.Fatalf("Upload: %v", err)
	}

	if gotMethod != http.MethodPut {
		t.Errorf("method: got %q, want PUT", gotMethod)
	}
	if gotPath != "/files/doc.txt" {
		t.Errorf("path: got %q, want /files/doc.txt", gotPath)
	}
	if string(gotBody) != string(payload) {
		t.Errorf("body: got %q, want %q", gotBody, payload)
	}
	if gotAuth == "" {
		t.Error("expected Authorization header")
	}
}

func TestWebDAVUploadError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusForbidden)
	}))
	defer srv.Close()

	dest := Destination{Type: "webdav", URL: srv.URL, Path: "/"}
	u, err := newWebDAVUploader(dest)
	if err != nil {
		t.Fatalf("newWebDAVUploader: %v", err)
	}

	err = u.Upload("file.txt", []byte("data"))
	if err == nil {
		t.Fatal("expected error for HTTP 403, got nil")
	}
}

func TestWebDAVUnsupportedType(t *testing.T) {
	dest := Destination{Type: "s3"}
	_, err := newWebDAVUploader(dest)
	if err == nil {
		t.Fatal("expected error for unsupported type, got nil")
	}
}
