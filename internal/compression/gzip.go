package compression

import (
	"bytes"
	"compress/gzip"
	"fmt"
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
