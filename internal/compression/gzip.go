package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"strings"
)

func CompressGzip(data []byte) ([]byte, error) {
	var b bytes.Buffer
	w, err := gzip.NewWriterLevel(&b, gzip.BestSpeed)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize gzip writer: %v", err)
	}
	_, err = w.Write(data)
	if err != nil {
		w.Close()
		return nil, fmt.Errorf("failed to write data to gzip writer: %v", err)
	}
	err = w.Close()
	if err != nil {
		return nil, fmt.Errorf("failed to finalize gzip compression: %v", err)
	}
	return b.Bytes(), nil
}

func WithUncompressGzip(h http.Handler) http.Handler {
	uncompressFn := func(w http.ResponseWriter, r *http.Request) {
		if strings.Contains(r.Header.Get("Content-Encoding"), "gzip") {
			gz, err := gzip.NewReader(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// закрытие gzip-читателя опционально, т.к. все данные уже прочитаны и
			// текущая реализация не требует закрытия, тем не менее лучше это делать -
			// некоторые реализации могут рассчитывать на закрытие читателя
			// gz.Close() не вызывает закрытия r.Body - это будет сделано позже, http-сервером
			defer gz.Close()

			// при чтении вернётся распакованный слайс байт
			decompressedBody, err := io.ReadAll(gz)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			r.Body = io.NopCloser(bytes.NewReader(decompressedBody))
			r.ContentLength = int64(len(decompressedBody))
			r.Header.Del("Content-Encoding")
		}
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(uncompressFn)
}
