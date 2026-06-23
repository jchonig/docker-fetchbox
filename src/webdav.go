package main

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"
)

type webDAVUploader struct {
	dest Destination
}

func newWebDAVUploader(dest Destination) (FileUploader, error) {
	if dest.Type != "webdav" {
		return nil, fmt.Errorf("unsupported destination type %q", dest.Type)
	}
	return &webDAVUploader{dest: dest}, nil
}

func (w *webDAVUploader) Upload(filename string, data []byte) error {
	base := strings.TrimRight(w.dest.URL, "/") + "/" + strings.TrimLeft(w.dest.Path, "/")
	url := strings.TrimRight(base, "/") + "/" + filename

	req, err := http.NewRequest(http.MethodPut, url, bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("build request: %w", err)
	}

	req.Header.Set("Content-Type", "application/octet-stream")
	if w.dest.Username != "" {
		req.SetBasicAuth(w.dest.Username, w.dest.Password())
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return fmt.Errorf("put %s: %w", url, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("put %s: HTTP %d", url, resp.StatusCode)
	}

	return nil
}
