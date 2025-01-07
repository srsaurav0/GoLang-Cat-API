package controllers

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/beego/beego/v2/server/web"
)

// FetchBreeds handles the request to fetch all cat breeds
func (c *CatController) FetchBreeds() {
	baseURL, _ := web.AppConfig.String("catapi_base_url")

	apiKey, _ := web.AppConfig.String("catapi_key")

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", baseURL+"/breeds", nil)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to create request")
		return
	}

	req.Header.Set("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to fetch breeds")
		return
	}
	defer resp.Body.Close()

	var breeds []map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.Ctx.WriteString("Failed to parse response")
		return
	}

	c.Data["json"] = breeds
	c.ServeJSON()
}
