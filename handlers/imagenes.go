package handlers

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func saveUploadedImage(r *http.Request, formField string) (*string, error) {
	file, header, err := r.FormFile(formField)
	if err != nil {
		return nil, nil
	}
	defer file.Close()

	// Validar tipo de archivo
	ext := strings.ToLower(filepath.Ext(header.Filename))
	allowed := map[string]bool{".jpg": true, ".jpeg": true, ".png": true, ".webp": true, ".gif": true}
	if !allowed[ext] {
		return nil, fmt.Errorf("tipo de archivo no permitido")
	}

	// Crear directorio si no existe
	if err := os.MkdirAll("uploads", 0755); err != nil {
		return nil, err
	}

	// Nombre único basado en timestamp
	filename := fmt.Sprintf("%d%s", time.Now().UnixNano(), ext)
	destPath := filepath.Join("uploads", filename)

	dest, err := os.Create(destPath)
	if err != nil {
		return nil, err
	}
	defer dest.Close()

	// Copiar: máximo 1MB
	maxSize := int64(1 << 20)
	if _, err := io.Copy(dest, io.LimitReader(file, maxSize)); err != nil {
		return nil, err
	}

	// Guardar la URL relativa que el cliente usará
	urlPath := "/uploads/" + filename
	return &urlPath, nil
}
