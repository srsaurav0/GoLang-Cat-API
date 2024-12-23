package controllers

import (
	// "encoding/json"
	// "fmt"
	// "net/http"
	// "sync"
	// "time"

	// "github.com/beego/beego/v2/server/web"
	beego "github.com/beego/beego/v2/server/web"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.TplName = "index.tpl"
}

// func (c *MainController) HomePage() {
// 	baseURL, _ := web.AppConfig.String("catapi_base_url")
// 	apiKey, _ := web.AppConfig.String("catapi_key")

// 	// Channels for concurrent fetching
// 	imagesChan := make(chan []map[string]interface{})
// 	breedsChan := make(chan []map[string]interface{})
// 	favoritesChan := make(chan []map[string]interface{})
// 	errorChan := make(chan error, 3)

// 	// WaitGroup to wait for all goroutines
// 	var wg sync.WaitGroup

// 	// Fetch cat images
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		fmt.Println("Fetching images: Started at", time.Now())
// 		client := &http.Client{}
// 		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/images/search?limit=10", baseURL), nil)
// 		req.Header.Set("x-api-key", apiKey)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("Error fetching images: %w", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		var images []map[string]interface{}
// 		if err := json.NewDecoder(resp.Body).Decode(&images); err != nil {
// 			errorChan <- fmt.Errorf("Error parsing images response: %w", err)
// 			return
// 		}
// 		fmt.Println("Fetching images: Completed at", time.Now())
// 		imagesChan <- images
// 	}()

// 	// Fetch breeds
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		fmt.Println("Fetching breeds: Started at", time.Now())
// 		client := &http.Client{}
// 		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/breeds", baseURL), nil)
// 		req.Header.Set("x-api-key", apiKey)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("Error fetching breeds: %w", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		var breeds []map[string]interface{}
// 		if err := json.NewDecoder(resp.Body).Decode(&breeds); err != nil {
// 			errorChan <- fmt.Errorf("Error parsing breeds response: %w", err)
// 			return
// 		}
// 		fmt.Println("Fetching breeds: Completed at", time.Now())
// 		breedsChan <- breeds
// 	}()

// 	// Fetch favorites
// 	wg.Add(1)
// 	go func() {
// 		defer wg.Done()
// 		fmt.Println("Fetching favorites: Started at", time.Now())
// 		client := &http.Client{}
// 		req, _ := http.NewRequest("GET", fmt.Sprintf("%s/favourites?sub_id=user-123", baseURL), nil)
// 		req.Header.Set("x-api-key", apiKey)

// 		resp, err := client.Do(req)
// 		if err != nil {
// 			errorChan <- fmt.Errorf("Error fetching favorites: %w", err)
// 			return
// 		}
// 		defer resp.Body.Close()

// 		var favorites []map[string]interface{}
// 		if err := json.NewDecoder(resp.Body).Decode(&favorites); err != nil {
// 			errorChan <- fmt.Errorf("Error parsing favorites response: %w", err)
// 			return
// 		}
// 		fmt.Println("Fetching favorites: Completed at", time.Now())
// 		favoritesChan <- favorites
// 	}()

// 	// Wait for all goroutines to complete
// 	go func() {
// 		wg.Wait()
// 		close(imagesChan)
// 		close(breedsChan)
// 		close(favoritesChan)
// 		close(errorChan)
// 	}()

// 	// Aggregate results
// 	var images []map[string]interface{}
// 	var breeds []map[string]interface{}
// 	var favorites []map[string]interface{}
// 	errors := []error{}

// 	for {
// 		select {
// 		case img, ok := <-imagesChan:
// 			if ok {
// 				images = img
// 				fmt.Println("Images fetched:", time.Now())
// 			} else {
// 				imagesChan = nil // Set to nil when channel is closed
// 			}
// 		case brd, ok := <-breedsChan:
// 			if ok {
// 				breeds = brd
// 				fmt.Println("Breeds fetched:", time.Now())
// 			} else {
// 				breedsChan = nil // Set to nil when channel is closed
// 			}
// 		case fav, ok := <-favoritesChan:
// 			if ok {
// 				favorites = fav
// 				fmt.Println("Favorites fetched:", time.Now())
// 			} else {
// 				favoritesChan = nil // Set to nil when channel is closed
// 			}
// 		case err, ok := <-errorChan:
// 			if ok {
// 				errors = append(errors, err)
// 				fmt.Println("Error occurred:", time.Now(), err)
// 			} else {
// 				errorChan = nil // Set to nil when channel is closed
// 			}
// 		}

// 		// Break when all channels are closed
// 		if imagesChan == nil && breedsChan == nil && favoritesChan == nil && errorChan == nil {
// 			break
// 		}
// 	}

// 	fmt.Println("Check1", time.Now())
// 	// Check for errors
// 	if len(errors) > 0 {
// 		for _, err := range errors {
// 			fmt.Println("Error:", err)
// 		}
// 		c.Ctx.Output.SetStatus(500)
// 		c.Data["json"] = map[string]string{"error": "Failed to fetch some data"}
// 		c.ServeJSON()
// 		return
// 	}

// 	// Pass data to the template
// 	c.Data["images"] = images
// 	c.Data["breeds"] = breeds
// 	c.Data["favorites"] = favorites

// 	// Render index.tpl
// 	c.TplName = "index.tpl"
// }
