package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/beego/beego/v2/server/web"
)

// CatController handles API requests
type CatController struct {
	web.Controller
	HTTPClient *http.Client
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
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}
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
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	// Read the raw body
	rawBody, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		if c.Data == nil {
			c.Data = make(map[interface{}]interface{}) // Initialize if nil
		}
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Failed to read request body"}
		c.ServeJSON()
		return
	}
	fmt.Println("Raw request body:", string(rawBody)) // Log the raw request body for debugging

	// Parse the request body into a struct
	var reqBody struct {
		ImageID string `json:"image_id"`
		SubID   string `json:"sub_id"`
	}
	if err := json.Unmarshal(rawBody, &reqBody); err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Invalid JSON format"}
		c.ServeJSON()
		return
	}
	fmt.Println("Parsed request body:", reqBody)

	// Channels to handle tasks
	addToFavoritesChan := make(chan error, 1)             // Buffered channel to avoid blocking
	nextImageChan := make(chan map[string]interface{}, 1) // Buffered channel
	errorChan := make(chan error, 2)                      // Buffered channel

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Task 1: Add image to favorites
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := addToFavorites(baseURL, apiKey, reqBody.ImageID, reqBody.SubID)
		addToFavoritesChan <- err // Send result to channel
	}()

	// Task 2: Fetch the next image
	wg.Add(1)
	go func() {
		defer wg.Done()
		images, err := fetchCatImages(baseURL, apiKey, "")
		if err != nil {
			errorChan <- err
			return
		}
		if len(images) > 0 {
			nextImageChan <- images[0] // Send the first image to the channel
		} else {
			errorChan <- fmt.Errorf("No images available")
		}
	}()

	// Close channels after all goroutines are done
	go func() {
		wg.Wait()
		close(addToFavoritesChan)
		close(nextImageChan)
		close(errorChan)
	}()

	// Aggregate results
	var addToFavoritesErr error
	var nextImage map[string]interface{}
	var errors []error

	for {
		select {
		case err, ok := <-addToFavoritesChan:
			if ok {
				addToFavoritesErr = err
			} else {
				addToFavoritesChan = nil
			}
		case img, ok := <-nextImageChan:
			if ok {
				nextImage = img
			} else {
				nextImageChan = nil
			}
		case err, ok := <-errorChan:
			if ok {
				errors = append(errors, err)
			} else {
				errorChan = nil
			}
		}

		// Break when all channels are closed
		if addToFavoritesChan == nil && nextImageChan == nil && errorChan == nil {
			break
		}
	}

	fmt.Println("Mark!:", nextImage)

	// Handle errors
	if len(errors) > 0 {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"error":   "Failed to complete tasks",
			"details": errors,
		}
		c.ServeJSON()
		return
	}

	// Return success with the next image
	c.Data["json"] = map[string]interface{}{
		"message":          "Image added to favorites and next image fetched",
		"next_image":       nextImage,
		"add_to_favorites": addToFavoritesErr == nil,
	}
	c.ServeJSON()
}

func addToFavorites(baseURL, apiKey, imageID, subID string) error {
	client := &http.Client{}
	payload := map[string]string{
		"image_id": imageID,
		"sub_id":   subID,
	}
	payloadBytes, _ := json.Marshal(payload)
	req, _ := http.NewRequest("POST", fmt.Sprintf("%s/favourites", baseURL), bytes.NewBuffer(payloadBytes))
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("Error adding to favorites: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Failed to add to favorites: %s", resp.Status)
	}
	return nil
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

func (c *CatController) GetFavourites() {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	subID := c.GetString("sub_id") // Get user ID from query string

	requestURL := baseURL + "/favourites"
	if subID != "" {
		requestURL += "?sub_id=" + subID
	}

	client := c.HTTPClient // Use the injected HTTP client
	if client == nil {
		client = &http.Client{} // Fallback to default if not injected
	}

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
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	// Read and parse the request body
	rawBody, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Failed to read request body"}
		c.ServeJSON()
		return
	}

	// Channels for concurrent tasks
	voteChan := make(chan map[string]interface{})
	nextImageChan := make(chan map[string]interface{})
	errorChan := make(chan error, 2)

	// WaitGroup to manage goroutines
	var wg sync.WaitGroup

	// Task 1: Submit Vote
	wg.Add(1)
	go func() {
		defer wg.Done()
		result, err := submitVote(baseURL, apiKey, rawBody)
		if err != nil {
			errorChan <- err
			return
		}
		voteChan <- result
	}()

	// Task 2: Fetch Next Image
	wg.Add(1)
	go func() {
		defer wg.Done()
		images, err := fetchCatImages(baseURL, apiKey, "")
		if err != nil {
			errorChan <- err
			return
		}
		if len(images) > 0 {
			nextImageChan <- images[0] // Send the first image to the channel
		} else {
			errorChan <- fmt.Errorf("No images available")
		}
	}()

	// Close channels after tasks complete
	go func() {
		wg.Wait()
		close(voteChan)
		close(nextImageChan)
		close(errorChan)
	}()

	// Aggregate results
	var voteResult map[string]interface{}
	var nextImage map[string]interface{}
	var errors []error

	for {
		select {
		case result, ok := <-voteChan:
			if ok {
				voteResult = result
			} else {
				voteChan = nil
			}
		case image, ok := <-nextImageChan:
			if ok {
				nextImage = image
			} else {
				nextImageChan = nil
			}
		case err, ok := <-errorChan:
			if ok {
				errors = append(errors, err)
			} else {
				errorChan = nil
			}
		}

		// Break when all channels are closed
		if voteChan == nil && nextImageChan == nil && errorChan == nil {
			break
		}
	}

	// Handle errors
	if len(errors) > 0 {
		fmt.Println("Error details:", errors) // Log all errors in the slice
		for _, err := range errors {
			fmt.Println("Individual error:", err) // Log individual errors
		}
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"error":   "Failed to complete tasks",
			"details": errors,
		}
		c.ServeJSON()
		return
	}

	fmt.Println("Mark 5")

	// Return success with vote result and next image
	c.Data["json"] = map[string]interface{}{
		"message":    "Vote submitted and next image fetched",
		"vote":       voteResult,
		"next_image": nextImage,
	}
	c.ServeJSON()
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
	fmt.Println("Raw response from The Cat API:", string(rawResponseBody))

	// Treat both 200 and 201 as successful responses
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("unexpected response status: %d, body: %s", resp.StatusCode, string(rawResponseBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(rawResponseBody, &result); err != nil {
		return nil, fmt.Errorf("failed to parse Cat API response: %w", err)
	}

	fmt.Println("Parsed response from Cat API:", result)
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
