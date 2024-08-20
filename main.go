package main

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"

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

type PrometheusQueryResult struct {
	Status string `json:"status"`
	Data   struct {
		Result []struct {
			Value []interface{} `json:"value"`
		} `json:"result"`
	} `json:"data"`
}

func main() {

	prometheusURL := loadenv("PROM_URL")
	query := "up{job=\"apiserver\"}"

	req, err := http.NewRequest("GET", prometheusURL, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	q := url.Values{}
	q.Add("query", query)
	req.URL.RawQuery = q.Encode()

	// Create an HTTP client with TLS config that skips certificate verification
	client := &http.Client{
		Timeout: 10 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
	}

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var result PrometheusQueryResult
	err = json.Unmarshal(body, &result)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	if result.Status == "success" && len(result.Data.Result) > 0 {
		// Check the value of the "up" metric
		status := result.Data.Result[0].Value[1].(string) // Value[1] is the "up" status
		if status == "1" {
			fmt.Println("OK")
		} else {
			fmt.Println("kube-apiserver is down")
		}
	} else {
		fmt.Println("Failed to retrieve kube-apiserver status or no results found")
	}
}
