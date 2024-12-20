package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
)

// CatController handles API requests
type CatController struct {
	web.Controller
}

// Structs to unmarshal Cat API responses
type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type FavouriteRequest struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
}

// FetchCatImages handles the request to fetch random cat images
func (c *CatController) FetchCatImages() {
	baseURL, err := web.AppConfig.String("catapi_base_url")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to load Cat API base URL from config")
		return
	}

	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to load Cat API key from config")
		return
	}

	// Get the breed_id from the request
	breedID := c.GetString("breed_id")
	requestURL := baseURL + "/images/search?limit=20"
	if breedID != "" {
		requestURL += "&breed_id=" + breedID // Append breed_id if provided
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to create request")
		return
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to fetch cat images")
		return
	}
	defer resp.Body.Close()

	var images []CatImage
	if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to parse response")
		return
	}

	c.Data["json"] = images
	c.ServeJSON()
}

// FetchBreeds handles the request to fetch all cat breeds
func (c *CatController) FetchBreeds() {
	baseURL, err := web.AppConfig.String("catapi_base_url")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to load API base URL from config")
		return
	}

	apiKey, err := web.AppConfig.String("catapi_key")
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to load API key from config")
		return
	}

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", baseURL+"/breeds", nil)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to create request")
		return
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to fetch breeds")
		return
	}
	defer resp.Body.Close()

	var breeds []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to parse response")
		return
	}

	c.Data["json"] = breeds
	c.ServeJSON()
}

func (c *CatController) AddToFavourites() {
	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	// Read the raw body
	rawBody, err := io.ReadAll(c.Ctx.Request.Body)
	fmt.Println("Raw request body:", string(rawBody)) // Debug the raw body received

	// Parse the JSON body
	var payload map[string]string
	fmt.Println("Payload is:", payload) // Log parsing errors
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		fmt.Println("Error unmarshaling JSON:", err) // Log parsing errors
		c.Ctx.Output.SetStatus(400)
		c.Ctx.WriteString("Invalid request body")
		return
	}

	// Debug the parsed payload
	fmt.Println("Parsed payload:", payload)

	// Create the request to The Cat API
	requestBody, _ := json.Marshal(payload)
	client := &http.Client{Timeout: 10 * time.Second}
	req, _ := http.NewRequest("POST", baseURL+"/favourites", bytes.NewBuffer(requestBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request to The Cat API:", err)
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to add favourite")
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Println("Error response from The Cat API:", resp.Status)
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Ctx.WriteString("Failed to add favourite")
		return
	}

	// Parse and forward the response from The Cat API
	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error parsing The Cat API response:", err)
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to parse response from Cat API")
		return
	}

	// Log and return the response
	fmt.Println("Response from The Cat API:", result)
	c.Data["json"] = result
	c.ServeJSON()
}

func (c *CatController) GetFavourites() {
	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	subID := c.GetString("sub_id") // Get user ID from query string

	requestURL := baseURL + "/favourites"
	if subID != "" {
		requestURL += "?sub_id=" + subID
	}

	client := &http.Client{}
	req, _ := http.NewRequest("GET", requestURL, nil)
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to fetch favourites")
		return
	}
	defer resp.Body.Close()

	var favourites []map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&favourites)
	c.Data["json"] = favourites
	c.ServeJSON()
}

func (c *CatController) RemoveFavourite() {
	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	favouriteID := c.GetString("favourite_id") // Get favourite ID from query string

	requestURL := baseURL + "/favourites/" + favouriteID

	client := &http.Client{}
	req, _ := http.NewRequest("DELETE", requestURL, nil)
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to remove favourite")
		return
	} else {
		c.Ctx.Output.SetStatus(resp.StatusCode)
		c.Ctx.WriteString("Favourite removed successfully")
	}
	defer resp.Body.Close()
}

func (c *CatController) Vote() {
	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	// Read the request body
	rawBody, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(400) // Set status before writing response
		c.Ctx.WriteString("Failed to read request body")
		return
	}

	// Parse the request body
	var payload map[string]interface{}
	if err := json.Unmarshal(rawBody, &payload); err != nil {
		c.Ctx.Output.SetStatus(400) // Set status before writing response
		c.Ctx.WriteString("Invalid request body")
		return
	}

	// Forward vote to The Cat API
	client := &http.Client{}
	req, err := http.NewRequest("POST", baseURL+"/votes", bytes.NewBuffer(rawBody))
	if err != nil {
		c.Ctx.Output.SetStatus(500) // Set status before writing response
		c.Ctx.WriteString("Failed to create vote request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.Output.SetStatus(500) // Set status before writing response
		c.Ctx.WriteString("Failed to submit vote")
		return
	}
	defer resp.Body.Close()

	// Log the raw response body from The Cat API
	rawResponseBody, _ := io.ReadAll(resp.Body)
	fmt.Println("Raw response from The Cat API:", string(rawResponseBody))

	if resp.StatusCode != http.StatusOK {
		errorMessage, _ := io.ReadAll(resp.Body) // Read the error message
		c.Ctx.Output.SetStatus(resp.StatusCode)  // Set the HTTP status
		c.Data["json"] = map[string]string{
			"error": fmt.Sprintf("Failed to submit vote: %s", string(errorMessage)),
		}
		c.ServeJSON() // Use ServeJSON to write the response
		return
	}

	// Parse and forward the response
	var result map[string]interface{}
	if err := json.Unmarshal(rawResponseBody, &result); err != nil {
		fmt.Println("Error parsing Cat API response:", err)
		c.Ctx.Output.SetStatus(500) // Set HTTP status
		c.Data["json"] = map[string]string{"error": "Failed to parse Cat API response"}
		c.ServeJSON() // Use ServeJSON to write the response
		return
	}

	fmt.Println("Parsed response from Cat API:", result)
	c.Data["json"] = result
	c.ServeJSON()
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
		fmt.Println("Error creating request to Cat API:", err)
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to create request")
		return
	}
	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request to Cat API:", err)
		c.Ctx.Output.SetStatus(500)
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
		fmt.Println("Error parsing Cat API response:", err)
		c.Ctx.Output.SetStatus(500)
		c.Ctx.WriteString("Failed to parse Cat API response")
		return
	}

	fmt.Println("Votes retrieved from Cat API:", result)
	c.Data["json"] = result
	c.ServeJSON()
}
