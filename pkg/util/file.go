package util

import (
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func ValidatePhoto(photo *multipart.FileHeader) error {
	// Open the photo file
	photoFile, err := photo.Open()
	if err != nil {
		return fmt.Errorf("failed to open photo file")
	}
	defer photoFile.Close()

	// Read the first 512 bytes to check MIME type
	buffer := make([]byte, 512)
	if _, err := photoFile.Read(buffer); err != nil {
		return fmt.Errorf("failed to read photo file")
	}

	// Detect content type
	contentType := http.DetectContentType(buffer)
	if contentType != "image/jpeg" && contentType != "image/png" {
		return fmt.Errorf("file must be a JPG, JPEG, or PNG image")
	}

	// Validate file extension
	ext := strings.ToLower(filepath.Ext(photo.Filename))
	if ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return fmt.Errorf("invalid file extension")
	}

	return nil
}
func GetFileExtension(photo *multipart.FileHeader) string {
	return strings.ToLower(filepath.Ext(photo.Filename))
}

func ConcatWithServerURL(serverURL, path string) string {
	return fmt.Sprintf("%s/%s", serverURL, path)
}

func DeleteFile(filePath string) error {

	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete file")
	}
	return nil
}
