package controllers

import (
	"encoding/json"
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
