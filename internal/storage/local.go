// Package storage implements file storage functionalities.
package storage

import (
	"crypto/sha256"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/SaidMg10/gestor-one/internal/domain"
)

type FileStorageLocal struct{}

func NewFileStorageLocal() domain.FileStorage {
	return &FileStorageLocal{}
}

func (fsl *FileStorageLocal) SavePDF(fileHeader *multipart.FileHeader) (string, string, string, error) {
	if fileHeader == nil {
		return "", "", "", fmt.Errorf("file is required")
	}

	if filepath.Ext(fileHeader.Filename) != ".pdf" {
		return "", "", "", fmt.Errorf("only PDF files are allowed")
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", "", "", fmt.Errorf("cannot open file: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot read uploaded file: %w", err)
	}

	hash := sha256.Sum256(fileBytes)
	checksum := fmt.Sprintf("%x", hash[:])

	filename := fmt.Sprintf("%d_receipt.pdf", time.Now().UnixNano())

	uploadDir := "./uploads"
	fullDiskPath := filepath.Join(uploadDir, filename)

	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", "", fmt.Errorf("cannot create uploads dir: %w", err)
	}

	out, err := os.Create(fullDiskPath)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot create file: %w", err)
	}
	defer func() {
		_ = out.Close()
	}()
	if _, err := file.Seek(0, 0); err != nil {
		return "", "", "", fmt.Errorf("cannot reset file pointer: %w", err)
	}

	if _, err := io.Copy(out, file); err != nil {
		return "", "", "", fmt.Errorf("cannot save file: %w", err)
	}

	// PATH RELATIVA para BD
	relPath := "/uploads/" + filename

	return filename, checksum, relPath, nil
}

func (fsl *FileStorageLocal) DeletePDF(filePath string) error {
	if filePath == "" {
		return nil
	}

	// recibes "/uploads/xxx.pdf"
	filename := filepath.Base(filePath)

	fullDiskPath := filepath.Join("./uploads", filename)

	if _, err := os.Stat(fullDiskPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return os.Remove(fullDiskPath)
}
