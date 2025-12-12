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
	defer func() { _ = file.Close() }()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot read uploaded file: %w", err)
	}

	hash := sha256.Sum256(fileBytes)
	checksum := fmt.Sprintf("%x", hash[:])

	filename := fmt.Sprintf("%d_receipt.pdf", time.Now().UnixNano())
	uploadDir := "./uploads"

	// Crear carpeta si no existe
	if err := os.MkdirAll(uploadDir, os.ModePerm); err != nil {
		return "", "", "", fmt.Errorf("cannot create uploads dir: %w", err)
	}

	fullPath := filepath.Join(uploadDir, filename)
	out, err := os.Create(fullPath)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot create file: %w", err)
	}
	defer func() { _ = out.Close() }()

	// Resetear el puntero para copiar nuevamente
	if _, err := file.Seek(0, 0); err != nil {
		return "", "", "", fmt.Errorf("cannot reset file pointer: %w", err)
	}

	if _, err := io.Copy(out, file); err != nil {
		return "", "", "", fmt.Errorf("cannot save file: %w", err)
	}

	// Convertir a ruta absoluta
	absPath, err := filepath.Abs(fullPath)
	if err != nil {
		return "", "", "", fmt.Errorf("cannot get absolute path: %w", err)
	}

	return filename, checksum, absPath, nil
}

func (fsl *FileStorageLocal) DeletePDF(filePath string) error {
	if filePath == "" {
		return nil
	}

	absPath, err := filepath.Abs(filePath)
	if err != nil {
		return fmt.Errorf("cannot get absolute path: %w", err)
	}

	if _, err := os.Stat(absPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}

	return os.Remove(absPath)
}
