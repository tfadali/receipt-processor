package main

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"
)

func TestReceiptProcessing(t *testing.T) {
	baseUrl := "http://localhost:8080"

	requestJson := readFileAsString(t, "testFiles/sampleReceipt1.json")
	idMap := getPostJsonResponse(t, baseUrl+"/receipts/process", requestJson)
	pointsMap := doGetRequestWithJsonResponse(t, baseUrl+"/receipts/"+idMap["id"]+"/points")

	assertEquals(t, 28, pointsMap["points"])
}

func TestReceiptProcessingAlternate(t *testing.T) {
	baseUrl := "http://localhost:8080"

	requestJson := readFileAsString(t, "testFiles/sampleReceipt2.json")
	idMap := getPostJsonResponse(t, baseUrl+"/receipts/process", requestJson)
	pointsMap := doGetRequestWithJsonResponse(t, baseUrl+"/receipts/"+idMap["id"]+"/points")

	assertEquals(t, 109, pointsMap["points"])
}

// read file and fail test if not successful.
func readFileAsString(t *testing.T, path string) string {
	bytes, err := os.ReadFile(path) // just pass the file name
	if err != nil {
		t.Errorf("Error while loading file: %q", path)
	}
	return string(bytes)
}

// do post request, get parsed response, fail test if not successful.
func getPostJsonResponse(t *testing.T, urlPath string, data string) map[string]string {
	response, err := http.Post(urlPath, "application/json", bytes.NewBuffer([]byte(data)))
	if err != nil {
		t.Errorf("Error sending POST request to %q: %q", urlPath, err)
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Errorf("Failed to read response from %q %q", urlPath, err)
	}
	var parsed map[string]string
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Errorf("Failed to parse json response from %q. %q", urlPath, err)
	}
	return parsed
}

// do get request, get parsed response, fail test if not successful.
func doGetRequestWithJsonResponse(t *testing.T, urlPath string) map[string]int {
	response, err := http.Get(urlPath)
	if err != nil {
		t.Errorf("Error sending GET request to %q %q", urlPath, err)
	}
	body, err := io.ReadAll(response.Body)
	defer response.Body.Close()
	if err != nil {
		t.Errorf("Failed to read request from %q: %q", urlPath, err)
	}

	var parsed map[string]int
	err = json.Unmarshal(body, &parsed)
	if err != nil {
		t.Errorf("Failed to parse response from %q : %q", urlPath, err)
	}
	return parsed
}

func assertEquals(t *testing.T, expected int, received int) {
	if expected != received {
		t.Errorf("Expected: %d, but was: %d", expected, received)
	}
}
