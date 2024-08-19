package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

// loads env from a .env file
func loadenv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}
	return os.Getenv(key)
}

func main() {

	promurl := loadenv("PROM_URL")
	clusterurl := loadenv("CLUSTER_URL")
	endpointurl := promurl + clusterurl

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	_, err := http.Get(endpointurl)
	if err != nil {
		fmt.Println(err)
	}
}
