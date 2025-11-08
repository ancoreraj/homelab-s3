package main

import (
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/ancoreraj/s3-clone-go/internal/handlers"
	"github.com/ancoreraj/s3-clone-go/internal/middleware"
	"github.com/ancoreraj/s3-clone-go/internal/services"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func getLocalIP() string {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "localhost"
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP.String()
}

func main() {
	godotenv.Load()

	storageService := services.NewStorageService("../../uploads")
	
	s3Handler := handlers.NewS3Handler(storageService)

	router := mux.NewRouter()

	router.HandleFunc("/health", s3Handler.HealthCheck).Methods("GET")
	router.HandleFunc("/upload/{bucket}", s3Handler.UploadObject).Methods("PUT")
	router.HandleFunc("/download/{bucket}/{key:.*}", s3Handler.DownloadObject).Methods("GET")
	router.HandleFunc("/list/{bucket}", s3Handler.ListBucket).Methods("GET")
	router.HandleFunc("/delete/{bucket}/{key:.*}", s3Handler.DeleteObject).Methods("DELETE")
	router.HandleFunc("/buckets", s3Handler.ListAllBuckets).Methods("GET")
	router.HandleFunc("/buckets", s3Handler.CreateBucket).Methods("POST")
	router.HandleFunc("/buckets/{bucket}", s3Handler.DeleteBucket).Methods("DELETE")

	publicDir := filepath.Join("..", "public")
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(publicDir)))

	handler := middleware.EnableCORS(router)
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	localIP := getLocalIP()
	log.Printf("S3 Clone server running at http://0.0.0.0:%s\n", port)
	log.Printf("Access from other machines using: http://%s:%s\n", localIP, port)

	if err := http.ListenAndServe("0.0.0.0:"+port, handler); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}