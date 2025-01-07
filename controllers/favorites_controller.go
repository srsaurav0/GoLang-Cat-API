package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sync"

	"github.com/beego/beego/v2/server/web"
)

// CatController handles API requests
type CatController struct {
	web.Controller
	HTTPClient *http.Client
	Client     HTTPClient
}

type HTTPClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type FavouriteRequest struct {
	ImageID string `json:"image_id"`
	SubID   string `json:"sub_id"`
}

func (c *CatController) AddToFavourites() {
	if c.Data == nil {
		c.Data = make(map[interface{}]interface{})
	}

	baseURL, _ := web.AppConfig.String("catapi_base_url")
	apiKey, _ := web.AppConfig.String("catapi_key")

	rawBody, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		fmt.Println("Error reading request body:", err)
		c.Ctx.Output.SetStatus(400)
		c.Data["json"] = map[string]string{"error": "Failed to read request body"}
		c.ServeJSON()
		return
	}
	fmt.Println("Raw request body:", string(rawBody)) // Log the raw request body for debugging

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

	addToFavoritesChan := make(chan error, 1)
	nextImageChan := make(chan map[string]interface{}, 1)
	errorChan := make(chan error, 2)

	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		err := addToFavorites(baseURL, apiKey, reqBody.ImageID, reqBody.SubID)
		addToFavoritesChan <- err
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		images, err := fetchCatImages(baseURL, apiKey, "")
		if err != nil {
			errorChan <- err
			return
		}
		if len(images) > 0 {
			nextImageChan <- images[0]
		} else {
			errorChan <- fmt.Errorf("No images available")
		}
	}()

	go func() {
		wg.Wait()
		close(addToFavoritesChan)
		close(nextImageChan)
		close(errorChan)
	}()

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

		if addToFavoritesChan == nil && nextImageChan == nil && errorChan == nil {
			break
		}
	}

	fmt.Println("Mark!:", nextImage)

	if len(errors) > 0 {
		c.Ctx.Output.SetStatus(500)
		c.Data["json"] = map[string]interface{}{
			"error":   "Failed to complete tasks",
			"details": errors,
		}
		c.ServeJSON()
		return
	}

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
