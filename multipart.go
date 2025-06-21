package httpclient

import (
	"bytes"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
)

// MultipartData represents multipart form data
type MultipartData struct {
	Fields map[string]string
	Files  map[string]*MultipartFile
}

// MultipartFile represents a file to be uploaded
type MultipartFile struct {
	Path     string
	Filename string
	Content  []byte
}

// NewMultipartData creates a new multipart data container
func NewMultipartData() *MultipartData {
	return &MultipartData{
		Fields: make(map[string]string),
		Files:  make(map[string]*MultipartFile),
	}
}

// AddField adds a form field
func (md *MultipartData) AddField(name, value string) {
	md.Fields[name] = value
}

// AddFileFromPath adds a file from file path
func (md *MultipartData) AddFileFromPath(fieldName, filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}

	filename := filepath.Base(filePath)
	md.Files[fieldName] = &MultipartFile{
		Path:     filePath,
		Filename: filename,
		Content:  content,
	}

	return nil
}

// AddFileFromBytes adds a file from bytes
func (md *MultipartData) AddFileFromBytes(fieldName, filename string, content []byte) {
	md.Files[fieldName] = &MultipartFile{
		Filename: filename,
		Content:  content,
	}
}

// ToReader converts multipart data to a reader
func (md *MultipartData) ToReader() (io.Reader, string, error) {
	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// Add fields
	for name, value := range md.Fields {
		if err := writer.WriteField(name, value); err != nil {
			return nil, "", err
		}
	}

	// Add files
	for fieldName, file := range md.Files {
		part, err := writer.CreateFormFile(fieldName, file.Filename)
		if err != nil {
			return nil, "", err
		}

		if _, err := part.Write(file.Content); err != nil {
			return nil, "", err
		}
	}

	if err := writer.Close(); err != nil {
		return nil, "", err
	}

	return &buf, writer.FormDataContentType(), nil
}

// MultipartRequestOptions extends RequestOptions with multipart support
type MultipartRequestOptions struct {
	*RequestOptions
	Multipart *MultipartData
}

// NewMultipartRequestOptions creates new multipart request options
func NewMultipartRequestOptions() *MultipartRequestOptions {
	return &MultipartRequestOptions{
		RequestOptions: &RequestOptions{},
		Multipart:      NewMultipartData(),
	}
}

// AddField adds a form field to multipart data
func (mro *MultipartRequestOptions) AddField(name, value string) {
	mro.Multipart.AddField(name, value)
}

// AddFile adds a file to multipart data
func (mro *MultipartRequestOptions) AddFile(fieldName, filePath string) error {
	return mro.Multipart.AddFileFromPath(fieldName, filePath)
}

// AddFileFromBytes adds a file from bytes to multipart data
func (mro *MultipartRequestOptions) AddFileFromBytes(fieldName, filename string, content []byte) {
	mro.Multipart.AddFileFromBytes(fieldName, filename, content)
} 