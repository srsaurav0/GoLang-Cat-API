package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/beego/beego/v2/server/web"
)

// Structs to unmarshal Cat API responses
type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

// FetchCatImages handles the request to fetch random cat images
func (c *CatController) FetchCatImages() {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}
	baseURL, _ := web.AppConfig.String("catapi_base_url")

	apiKey, _ := web.AppConfig.String("catapi_key")

	// Get the breed_id from the request
	breedID := c.GetString("breed_id")

	// Fetch images using the helper function
	images, err := fetchCatImages(baseURL, apiKey, breedID)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]string{"error": err.Error()}
		c.ServeJSON()
		return
	}

	// Return the fetched images as JSON
	c.Data["json"] = images
	c.ServeJSON()
}

func fetchCatImages(baseURL, apiKey, breedID string) ([]map[string]interface{}, error) {
	client := &http.Client{}

	// Construct the request URL
	requestURL := fmt.Sprintf("%s/images/search?limit=15", baseURL)
	if breedID != "" {
		requestURL += fmt.Sprintf("&breed_id=%s", breedID) // Append breed_id if provided
	}

	// Create the request
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to create request: %w", err)
	}
	req.Header.Set("x-api-key", apiKey)

	// Execute the request
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Failed to fetch images: %w", err)
	}
	defer resp.Body.Close()

	// Check for unexpected status codes
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Unexpected response status: %d", resp.StatusCode)
	}

	// Parse the response body
	var images []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return images, nil
}

func submitVote(baseURL, apiKey string, payload []byte) (map[string]interface{}, error) {
	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL+"/votes", bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create vote request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to submit vote: %w", err)
	}
	defer resp.Body.Close()

	rawResponseBody, _ := io.ReadAll(resp.Body)
	log.Println("Raw response from The Cat API:", string(rawResponseBody))

	// Treat both 200 and 201 as successful responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, string(rawResponseBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rawResponseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Cat API response: %w", err)
	}

	log.Println("Parsed response from Cat API:", result)
	return result, nil
}

func (c *CatController) GetVotes() {
	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	// Forward query parameters
	queryParams := c.Ctx.Request.URL.RawQuery

	// Create the GET request to The Cat API
	client := &http.Client{}
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/votes?%s", baseURL, queryParams), nil)
	if err != nil {
		// fmt.Println("Error creating request to Cat API:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to create request")
		return
	}
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		log.Println("Error sending request to Cat API:", err)
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to retrieve votes")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Ctx.WriteString("Failed to retrieve votes")
		return
	}

	var result []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		log.Println("Error parsing Cat API response:", err)
		c.Ctx.Output.SetStatus(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to parse Cat API response")
		return
	}

	log.Println("Votes retrieved from Cat API:", result)
	c.Data["json"] = result
	c.ServeJSON()
}
