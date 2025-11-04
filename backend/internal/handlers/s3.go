package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/ancoreraj/s3-clone-go/internal/services"
	"github.com/gorilla/mux"
)

type S3Handler struct {
	storage *services.StorageService
}

func NewS3Handler(storage *services.StorageService) *S3Handler {
	return &S3Handler{storage: storage}
}

type HealthResponse struct {
	Message   string     `json:"message"`
	Endpoints []Endpoint `json:"endpoints"`
}

type Endpoint struct {
	Method      string `json:"method"`
	Path        string `json:"path"`
	Description string `json:"description"`
}

type UploadResponse struct {
	Message  string `json:"message"`
	Bucket   string `json:"bucket"`
	Key      string `json:"key"`
	Size     int64  `json:"size"`
	MimeType string `json:"mimetype"`
}

type BucketListResponse struct {
	Bucket string   `json:"bucket"`
	Files  []string `json:"files"`
}

type BucketsResponse struct {
	Buckets []string `json:"buckets"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

type CreateBucketRequest struct {
	Name string `json:"name"`
}

func (h *S3Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	response := HealthResponse{
		Message: "S3 Clone API is running",
		Endpoints: []Endpoint{
			{Method: "PUT", Path: "/upload/:bucket/:key", Description: "Upload a file to a bucket"},
			{Method: "GET", Path: "/download/:bucket/:key", Description: "Download a file from a bucket"},
			{Method: "GET", Path: "/list/:bucket", Description: "List all files in a bucket"},
			{Method: "DELETE", Path: "/delete/:bucket/:key", Description: "Delete a file from a bucket"},
		},
	}
	jsonResponse(w, http.StatusOK, response)
}

func (h *S3Handler) UploadObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	
	r.ParseMultipartForm(100 << 20) // 100MB max
	
	file, header, err := r.FormFile("file")
	if err != nil {
		jsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "No file uploaded"})
		return
	}
	defer file.Close()

	key := r.URL.Query().Get("key")
	if key == "" {
		key = header.Filename
	} else {
		ext := h.storage.GetFileExtension(header.Header.Get("Content-Type"), key)
		if ext != "" && !strings.HasSuffix(strings.ToLower(key), "."+strings.ToLower(ext)) {
			key = key + "." + ext
		}
	}

	size, err := h.storage.SaveFile(bucket, key, file)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to upload file"})
		return
	}

	response := UploadResponse{
		Message:  "File uploaded successfully",
		Bucket:   bucket,
		Key:      key,
		Size:     size,
		MimeType: header.Header.Get("Content-Type"),
	}
	jsonResponse(w, http.StatusOK, response)
}

func (h *S3Handler) DownloadObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	key := vars["key"]

	filePath, err := h.storage.GetFilePath(bucket, key)
	if err != nil {
		jsonResponse(w, http.StatusNotFound, ErrorResponse{Error: "File not found"})
		return
	}

	http.ServeFile(w, r, filePath)
}

func (h *S3Handler) ListBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	files, err := h.storage.ListBucketContents(bucket)
	if err != nil {
		if err == services.ErrBucketNotFound {
			jsonResponse(w, http.StatusNotFound, ErrorResponse{Error: "Bucket not found"})
		} else {
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to list bucket contents"})
		}
		return
	}

	response := BucketListResponse{
		Bucket: bucket,
		Files:  files,
	}
	jsonResponse(w, http.StatusOK, response)
}

func (h *S3Handler) DeleteObject(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]
	key := vars["key"]

	err := h.storage.DeleteFile(bucket, key)
	if err != nil {
		if err == services.ErrFileNotFound {
			jsonResponse(w, http.StatusNotFound, ErrorResponse{Error: "File not found"})
		} else {
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to delete file"})
		}
		return
	}

	response := MessageResponse{
		Message: fmt.Sprintf("File %s deleted successfully", key),
	}
	jsonResponse(w, http.StatusOK, response)
}

func (h *S3Handler) ListAllBuckets(w http.ResponseWriter, r *http.Request) {
	buckets, err := h.storage.ListAllBuckets()
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to list buckets"})
		return
	}

	response := BucketsResponse{
		Buckets: buckets,
	}
	jsonResponse(w, http.StatusOK, response)
}

func (h *S3Handler) CreateBucket(w http.ResponseWriter, r *http.Request) {
	var req CreateBucketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Invalid request body"})
		return
	}

	if req.Name == "" {
		jsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Bucket name is required"})
		return
	}

	if !isValidBucketName(req.Name) {
		jsonResponse(w, http.StatusBadRequest, ErrorResponse{
			Error: "Invalid bucket name. Use only letters, numbers, dashes, and underscores.",
		})
		return
	}

	created, err := h.storage.CreateBucket(req.Name)
	if err != nil {
		jsonResponse(w, http.StatusInternalServerError, ErrorResponse{Error: "Failed to create bucket"})
		return
	}

	if !created {
		jsonResponse(w, http.StatusConflict, ErrorResponse{
			Error: fmt.Sprintf("Bucket '%s' already exists", req.Name),
		})
		return
	}

	response := MessageResponse{
		Message: fmt.Sprintf("Bucket '%s' created successfully", req.Name),
	}
	jsonResponse(w, http.StatusCreated, response)
}

func (h *S3Handler) DeleteBucket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bucket := vars["bucket"]

	if bucket == "" {
		jsonResponse(w, http.StatusBadRequest, ErrorResponse{Error: "Bucket name is required"})
		return
	}

	deleted, err := h.storage.DeleteBucket(bucket)
	if err != nil {
		if err == services.ErrBucketNotFound {
			jsonResponse(w, http.StatusNotFound, ErrorResponse{
				Error: fmt.Sprintf("Bucket '%s' not found", bucket),
			})
		} else if err == services.ErrBucketNotEmpty {
			jsonResponse(w, http.StatusConflict, ErrorResponse{
				Error: fmt.Sprintf("Cannot delete bucket '%s': bucket is not empty", bucket),
			})
		} else {
			jsonResponse(w, http.StatusInternalServerError, ErrorResponse{
				Error: fmt.Sprintf("Failed to delete bucket '%s'", bucket),
			})
		}
		return
	}

	if deleted {
		response := MessageResponse{
			Message: fmt.Sprintf("Bucket '%s' deleted successfully", bucket),
		}
		jsonResponse(w, http.StatusOK, response)
	}
}

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func isValidBucketName(name string) bool {
	for _, char := range name {
		if !((char >= 'a' && char <= 'z') || 
			(char >= 'A' && char <= 'Z') || 
			(char >= '0' && char <= '9') || 
			char == '-' || char == '_') {
			return false
		}
	}
	return true
}