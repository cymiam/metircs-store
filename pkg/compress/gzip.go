package compress

import (
	"bytes"
	"compress/gzip"
	"io"
)

func GzipDecompress(data io.Reader) (*gzip.Reader, error) {
	gz, err := gzip.NewReader(data)
	if err != nil {
		return nil, err
	}

	return gz, nil
}

func GzipCompress(data []byte) ([]byte, error) {
	var buffer bytes.Buffer

	w := gzip.NewWriter(&buffer)

	_, err := w.Write(data)
	if err != nil {
		return nil, err
	}

	err = w.Close()

	if err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}
