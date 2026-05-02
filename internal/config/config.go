package config

import (
    "fmt"
    "os"
    "strconv"
)

type Config struct {
    Port            string
    DatabaseURL     string
    UploadDir       string
    MaxUploadSizeMB int64
}

func Load() (*Config, error) {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    dbURL := os.Getenv("DATABASE_URL")
    if dbURL == "" {
        return nil, fmt.Errorf("DATABASE_URL no está configurado")
    }

    uploadDir := os.Getenv("UPLOAD_DIR")
    if uploadDir == "" {
        uploadDir = "./static/uploads"
    }

    maxMB, err := strconv.ParseInt(os.Getenv("MAX_UPLOAD_SIZE_MB"), 10, 64)
    if err != nil || maxMB <= 0 {
        maxMB = 1
    }

    return &Config{
        Port:            port,
        DatabaseURL:     dbURL,
        UploadDir:       uploadDir,
        MaxUploadSizeMB: maxMB,
    }, nil
}