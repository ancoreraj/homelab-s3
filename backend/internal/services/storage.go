package services

import (
	"errors"
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrBucketNotFound = errors.New("bucket not found")
	ErrFileNotFound   = errors.New("file not found")
	ErrBucketNotEmpty = errors.New("bucket is not empty")
)

type StorageService struct {
	basePath string
}

func NewStorageService(basePath string) *StorageService {
	os.MkdirAll(basePath, 0755)
	return &StorageService{basePath: basePath}
}

func (s *StorageService) SaveFile(bucket, key string, file io.Reader) (int64, error) {
	bucketPath := filepath.Join(s.basePath, bucket)
	if err := os.MkdirAll(bucketPath, 0755); err != nil {
		return 0, err
	}

	filePath := filepath.Join(bucketPath, key)
	dst, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer dst.Close()

	return io.Copy(dst, file)
}

func (s *StorageService) GetFilePath(bucket, key string) (string, error) {
	filePath := filepath.Join(s.basePath, bucket, key)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", ErrFileNotFound
	}
	
	return filePath, nil
}

func (s *StorageService) DeleteFile(bucket, key string) error {
	filePath := filepath.Join(s.basePath, bucket, key)
	
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return ErrFileNotFound
	}
	
	return os.Remove(filePath)
}

func (s *StorageService) ListBucketContents(bucket string) ([]string, error) {
	bucketPath := filepath.Join(s.basePath, bucket)
	
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return nil, ErrBucketNotFound
	}

	files, err := os.ReadDir(bucketPath)
	if err != nil {
		return nil, err
	}

	fileNames := []string{}
	for _, file := range files {
		if !file.IsDir() {
			fileNames = append(fileNames, file.Name())
		}
	}

	return fileNames, nil
}

func (s *StorageService) ListAllBuckets() ([]string, error) {
	files, err := os.ReadDir(s.basePath)
	if err != nil {
		if os.IsNotExist(err) {
			return []string{}, nil
		}
		return nil, err
	}

	buckets := []string{}
	for _, file := range files {
		if file.IsDir() {
			buckets = append(buckets, file.Name())
		}
	}

	return buckets, nil
}

func (s *StorageService) CreateBucket(name string) (bool, error) {
	bucketPath := filepath.Join(s.basePath, name)
	
	if _, err := os.Stat(bucketPath); !os.IsNotExist(err) {
		return false, nil
	}

	if err := os.MkdirAll(bucketPath, 0755); err != nil {
		return false, err
	}

	return true, nil
}

func (s *StorageService) DeleteBucket(name string) (bool, error) {
	bucketPath := filepath.Join(s.basePath, name)
	
	if _, err := os.Stat(bucketPath); os.IsNotExist(err) {
		return false, ErrBucketNotFound
	}

	files, err := os.ReadDir(bucketPath)
	if err != nil {
		return false, err
	}

	if len(files) > 0 {
		return false, ErrBucketNotEmpty
	}

	if err := os.Remove(bucketPath); err != nil {
		return false, err
	}

	return true, nil
}

func (s *StorageService) GetFileExtension(mimeType, filename string) string {
	// First check if filename already has an extension
	if ext := filepath.Ext(filename); ext != "" {
		return strings.TrimPrefix(ext, ".")
	}

	mimeToExt := map[string]string{
		"image/jpeg":      "jpg",
		"image/jpg":       "jpg",
		"image/png":       "png",
		"image/gif":       "gif",
		"image/bmp":       "bmp",
		"image/webp":      "webp",
		"image/svg+xml":   "svg",
		"text/plain":      "txt",
		"text/html":       "html",
		"text/css":        "css",
		"text/csv":        "csv",
		"text/javascript": "js",
		"application/json": "json",
		"application/xml":  "xml",
		"application/pdf":  "pdf",
		"application/zip":  "zip",
		"application/x-zip-compressed": "zip",
		"application/x-7z-compressed":  "7z",
		"application/x-tar":            "tar",
		"application/x-rar-compressed": "rar",
		"application/msword":           "doc",
		"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   "docx",
		"application/vnd.ms-excel":                                                  "xls",
		"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet":         "xlsx",
		"application/vnd.ms-powerpoint":                                             "ppt",
		"application/vnd.openxmlformats-officedocument.presentationml.presentation": "pptx",
		"audio/mpeg":     "mp3",
		"audio/wav":      "wav",
		"audio/ogg":      "ogg",
		"video/mp4":      "mp4",
		"video/mpeg":     "mpeg",
		"video/quicktime": "mov",
		"video/webm":     "webm",
	}

	if ext, ok := mimeToExt[mimeType]; ok {
		return ext
	}

	exts, _ := mime.ExtensionsByType(mimeType)
	if len(exts) > 0 {
		return strings.TrimPrefix(exts[0], ".")
	}

	return ""
}