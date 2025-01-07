package test

// func TestGetVotes_CreateRequestError(t *testing.T) {
// 	// Create a new controller instance
// 	controller := &controllers.CatController{}

// 	// Create a new request
// 	req := httptest.NewRequest("GET", "/api/votes", nil)
// 	req.Header.Set("Content-Type", "application/json")
// 	w := httptest.NewRecorder()

// 	// Create Beego context
// 	beegoContext := context.NewContext()
// 	beegoContext.Reset(w, req)

// 	// Initialize controller with Beego context
// 	controller.Init(beegoContext, "CatController", "GetVotes", nil)

// 	// Mock the web.AppConfig to return invalid base URL to trigger the error
// 	originalBaseURL, _ := web.AppConfig.String("catapi_base_url")
// 	originalApiKey, _ := web.AppConfig.String("catapi_key")
// 	defer func() {
// 		web.AppConfig.Set("catapi_base_url", originalBaseURL)
// 		web.AppConfig.Set("catapi_key", originalApiKey)
// 	}()

// 	web.AppConfig.Set("catapi_base_url", "invalid-url")

// 	// Execute the method
// 	controller.GetVotes()

// 	// Assert response
// 	assert.Equal(t, http.StatusInternalServerError, w.Code, "Status code should match expected")
// 	respBody := w.Body.String()
// 	t.Logf("Response Body: %s", respBody) // Add this for debugging
// 	assert.Contains(t, respBody, "Failed to create request", "Response should contain expected message")
// }
