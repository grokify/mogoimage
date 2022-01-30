package pigoutil

import (
	"bytes"
	_ "embed"
	"io"
	"net/http"
	"os"
)

// DetectContentTypeFile detects the file type by reading MIME type information of the file content.
func DetectContentTypeFile(fname string) (interface{}, error) {
	file, err := os.Open(fname)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return DetectContentType(file)
}

// DetectContentTypeBytes detects the file type by reading MIME type information of the file content.
func DetectContentTypeBytes(b []byte) (string, error) {
	return DetectContentType(bytes.NewReader(b))
}

// DetectContentType detects the file type by reading MIME type information of the file content.
func DetectContentType(rs io.ReadSeeker) (string, error) {
	// Only the first 512 bytes are used to sniff the content type.
	buffer := make([]byte, 512)
	_, err := rs.Read(buffer)
	if err != nil {
		return "", err
	}

	// Reset the read pointer if necessary.
	rs.Seek(0, 0)

	// Always returns a valid content-type and "application/octet-stream" if no others seemed to match.
	contentType := http.DetectContentType(buffer)

	return string(contentType), nil
}
