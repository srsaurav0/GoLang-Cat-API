package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"cat-voting-app/controllers"
	_ "cat-voting-app/controllers"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

func init() {
	// Initialize Beego configuration with mock values
	err := web.AppConfig.Set("catapi_base_url", "http://default-test-url")
	if err != nil {
		panic("Failed to set default test config: " + err.Error())
	}
	err = web.AppConfig.Set("catapi_key", "default-test-key")
	if err != nil {
		panic("Failed to set default test config: " + err.Error())
	}
}

func setupTest(url string) (*controllers.CatController, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", url, nil)
	recorder := httptest.NewRecorder()

	beegoCtx := &context.Context{
		Input:  context.NewInput(),
		Output: context.NewOutput(),
	}
	beegoCtx.Reset(recorder, req)

	controller := &controllers.CatController{}
	controller.Init(beegoCtx, "CatController", "FetchCatImages", nil)

	return controller, recorder
}

// Mock HTTP server to simulate Cat API responses
func setupMockServer() *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("x-api-key") != "test_api_key" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		mockBreeds := []map[string]interface{}{
			{
				"id":   "abys",
				"name": "Abyssinian",
			},
			{
				"id":   "beng",
				"name": "Bengal",
			},
		}
		json.NewEncoder(w).Encode(mockBreeds)
	}))
}

func setupBreedTest(t *testing.T) (*controllers.CatController, *httptest.ResponseRecorder) {
	req := httptest.NewRequest("GET", "/api/breeds", nil)
	recorder := httptest.NewRecorder()

	beegoCtx := &context.Context{
		Input:  context.NewInput(),
		Output: context.NewOutput(),
	}
	beegoCtx.Reset(recorder, req)

	controller := &controllers.CatController{}
	controller.Init(beegoCtx, "CatController", "FetchBreeds", nil)

	return controller, recorder
}

// func setConfig(baseURL, apiKey string) error {
// 	if err := web.AppConfig.Set("catapi_base_url", baseURL); err != nil {
// 		return err
// 	}
// 	if err := web.AppConfig.Set("catapi_key", apiKey); err != nil {
// 		return err
// 	}
// 	return nil
// }

func clearBreedConfig() {
	web.AppConfig.Set("catapi_base_url", "")
	web.AppConfig.Set("catapi_key", "")
}

// func TestCatController_FetchBreeds_RequestCreationError(t *testing.T) {
// 	// Setup
// 	catController := &controllers.CatController{}

// 	// Create a test context
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/breeds", nil)
// 	ctx := context.NewContext()
// 	ctx.Reset(w, r)
// 	ctx.Output = context.NewOutput()
// 	ctx.Output.Context = ctx
// 	catController.Ctx = ctx

// 	// Initialize the controller
// 	catController.Init(ctx, "CatController", "FetchBreeds", nil)

// 	// Set an invalid URL that should cause NewRequest to fail
// 	web.AppConfig.Set("catapi_base_url", "\x00invalid")
// 	web.AppConfig.Set("catapi_key", "dummy-key")

// 	// Execute
// 	catController.FetchBreeds()

// 	// Debug information
// 	t.Logf("Response Status Code: %d", w.Code)
// 	t.Logf("Response Body: %q", w.Body.String())
// 	t.Logf("Response Headers: %v", w.Header())

// 	// Print the actual response for debugging
// 	fmt.Printf("Status Code: %d\n", w.Code)
// 	fmt.Printf("Response Body: %s\n", w.Body.String())

// 	// Assertions with more detailed failure messages
// 	if w.Code != 500 {
// 		t.Errorf("Expected status code 500, got %d", w.Code)
// 	}

// 	if body := w.Body.String(); body != "Failed to create request" {
// 		t.Errorf("Expected body 'Failed to create request', got %q", body)
// 	}
// }

type mockClient struct{}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
	return nil, errors.New("network error")
}

// func TestCatController_FetchBreeds_ClientError(t *testing.T) {
// 	// Setup
// 	catController := &controllers.CatController{}

// 	// Create a test context
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/breeds", nil)
// 	beegoCtx := &context.Context{
// 		Input:  context.NewInput(),
// 		Output: context.NewOutput(),
// 	}
// 	beegoCtx.Reset(w, r)
// 	catController.Init(beegoCtx, "CatController", "FetchBreeds", nil)

// 	// Set valid URL and API key
// 	web.AppConfig.Set("catapi_base_url", "https://api.example.com")
// 	web.AppConfig.Set("catapi_key", "dummy-key")

// 	// Execute
// 	catController.FetchBreeds()

// 	// Assertions
// 	assert.Equal(t, http.StatusInternalServerError, w.Code, "Expected 500 status code")
// 	assert.Contains(t, w.Body.String(), "Failed to fetch breeds", "Expected error message doesn't match")
// }

func TestFetchCatImages_FetchError(t *testing.T) {
	controller, recorder := setupTest("/api/cats?breed_id=abc123")

	// Set invalid base URL to trigger fetch error
	_ = web.AppConfig.Set("catapi_base_url", "http://invalid-url-that-will-fail")
	_ = web.AppConfig.Set("catapi_key", "test_api_key")

	controller.FetchCatImages()

	if recorder.Code != http.StatusInternalServerError {
		t.Errorf("Expected status code 500, but got %d", recorder.Code)
	}

	// Check if response is JSON and contains error message
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to be application/json, but got %s", contentType)
	}

	responseBody := recorder.Body.String()
	if !strings.Contains(responseBody, "error") {
		t.Errorf("Expected error message in response body, but got: %s", responseBody)
	}
}

// mockClientInvalidJSON is a custom HTTP client that returns invalid JSON
type mockClientInvalidJSON struct{}

func (m *mockClientInvalidJSON) Do(req *http.Request) (*http.Response, error) {
	// Create a new response with invalid JSON
	r := &http.Response{
		Status:     "200 OK",
		StatusCode: 200,
		Body:       io.NopCloser(strings.NewReader(`{"invalid json`)), // Malformed JSON
	}
	return r, nil
}

// func TestCatController_FetchBreeds_InvalidJSON(t *testing.T) {
// 	// Setup
// 	catController := &controllers.CatController{
// 		Client: &mockClientInvalidJSON{},
// 	}

// 	// Create a test context
// 	w := httptest.NewRecorder()
// 	r := httptest.NewRequest("GET", "/breeds", nil)
// 	ctx := context.NewContext()

// 	// Create a new response writer wrapper
// 	ctx.Input = context.NewInput()
// 	ctx.Request = r
// 	ctx.Output = context.NewOutput()
// 	ctx.ResponseWriter = &context.Response{
// 		ResponseWriter: w,
// 	}

// 	// Initialize the controller
// 	catController.Init(ctx, "CatController", "FetchBreeds", nil)

// 	// Set valid URL and API key
// 	web.AppConfig.Set("catapi_base_url", "https://api.example.com")
// 	web.AppConfig.Set("catapi_key", "dummy-key")

// 	// Execute
// 	catController.FetchBreeds()

// 	// Assertions
// 	assert.Equal(t, 500, w.Code, "Expected 500 status code")
// 	assert.Equal(t, "Failed to parse response", w.Body.String(), "Expected error message doesn't match")
// }

func TestFetchBreeds_Success(t *testing.T) {
	// Start mock server
	mockServer := setupMockServer()
	defer mockServer.Close()

	controller, recorder := setupBreedTest(t)

	// Set up valid configuration
	clearBreedConfig()
	web.AppConfig.Set("catapi_base_url", mockServer.URL)
	web.AppConfig.Set("catapi_key", "test_api_key")

	// Call the controller method
	controller.FetchBreeds()

	// Check status code
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", recorder.Code)
	}

	// Verify response is JSON
	contentType := recorder.Header().Get("Content-Type")
	if !strings.Contains(contentType, "application/json") {
		t.Errorf("Expected Content-Type to be application/json, but got %s", contentType)
	}

	// Verify response contains breeds
	var breeds []map[string]interface{}
	if err := json.NewDecoder(recorder.Body).Decode(&breeds); err != nil {
		t.Errorf("Failed to decode response body: %v", err)
	}

	if len(breeds) != 2 {
		t.Errorf("Expected 2 breeds, but got %d", len(breeds))
	}
}

func TestFetchCatImages(t *testing.T) {
	// Create a mock HTTP request
	req := httptest.NewRequest("GET", "/api/cats?breed_id=abc123", nil)
	recorder := httptest.NewRecorder()

	// Initialize the Beego context
	beegoCtx := &context.Context{
		Input:  context.NewInput(),
		Output: context.NewOutput(),
	}
	beegoCtx.Reset(recorder, req)

	// Set up the CatController
	controller := &controllers.CatController{}
	controller.Init(beegoCtx, "CatController", "FetchCatImages", nil)

	// Mock configuration settings
	_ = web.AppConfig.Set("catapi_base_url", "https://api.thecatapi.com/v1")
	_ = web.AppConfig.Set("catapi_key", "test_api_key")

	// Call the FetchCatImages method
	controller.FetchCatImages()

	// Validate the response
	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status code 200, but got %d", recorder.Code)
	}

	expectedContentType := "application/json"
	if contentType := recorder.Header().Get("Content-Type"); !strings.Contains(contentType, expectedContentType) {
		t.Errorf("Expected Content-Type to include %s, but got %s", expectedContentType, contentType)
	}

	// You can also validate the response body, if necessary
	responseBody := recorder.Body.String()
	if responseBody == "" {
		t.Errorf("Expected response body, but got empty")
	}
}

func TestAddToFavourites(t *testing.T) {
	payload := map[string]string{
		"image_id": "test-image-id", // Must not be empty
		"sub_id":   "user-123",
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		t.Fatalf("Failed to marshal payload: %v", err)
	}

	req := httptest.NewRequest("POST", "/api/add-to-favourites", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	controller := &controllers.CatController{}
	mockCtx := context.NewContext()
	mockCtx.Reset(resp, req)
	controller.Ctx = mockCtx

	controller.AddToFavourites()

	if resp.Code != http.StatusOK {
		t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	}

	var responseBody map[string]interface{}
	err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	t.Logf("Response: %v", responseBody)
}

type brokenReader struct{}

func (b *brokenReader) Read([]byte) (int, error) {
	return 0, fmt.Errorf("simulated read error")
}

func (b *brokenReader) Close() error {
	return nil
}

func TestAddToFavourites_ErrorReadingBody(t *testing.T) {
	// Simulate a request with a broken reader
	req := httptest.NewRequest("POST", "/api/add-to-favourites", &brokenReader{})
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()

	controller := &controllers.CatController{}
	mockCtx := context.NewContext()
	mockCtx.Reset(resp, req)
	controller.Ctx = mockCtx

	controller.AddToFavourites()

	// Check that the response status code is 400
	if resp.Code != http.StatusBadRequest {
		t.Errorf("Expected status code %d, got %d", http.StatusBadRequest, resp.Code)
	}

	// Check the response body for the appropriate error message
	var responseBody map[string]interface{}
	err := json.Unmarshal(resp.Body.Bytes(), &responseBody)
	if err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	expectedError := "Failed to read request body"
	if responseBody["error"] != expectedError {
		t.Errorf("Expected error message '%s', got '%s'", expectedError, responseBody["error"])
	}

	t.Logf("Error response: %v", responseBody)
}

func TestGetFavourites(t *testing.T) {
	// Mock data
	mockFavourites := []map[string]interface{}{
		{"id": 1, "image_id": "3q", "sub_id": "user-7899"},
	}
	mockResponse, _ := json.Marshal(mockFavourites)

	// Mock transport
	mockTransport := &mockTransport{
		responseBody: mockResponse,
		statusCode:   http.StatusOK,
	}
	mockClient := &http.Client{Transport: mockTransport}

	// Mock request and response
	req := httptest.NewRequest("GET", "/api/get-favourites?sub_id=user-7899", nil)
	resp := httptest.NewRecorder()

	// Mock context
	mockCtx := &context.Context{}
	mockCtx.Request = req
	mockCtx.ResponseWriter = &context.Response{ResponseWriter: resp}
	mockCtx.Input = &context.BeegoInput{Context: mockCtx}
	mockCtx.Output = &context.BeegoOutput{Context: mockCtx}

	// Mock controller with injected client
	controller := &controllers.CatController{
		HTTPClient: mockClient, // Inject mocked HTTP client
	}
	controller.Ctx = mockCtx

	// Call the method
	controller.GetFavourites()

	// Validate the response
	if resp.Code != http.StatusOK {
		t.Fatalf("Expected status 200, got %d", resp.Code)
	}

	var responseBody []map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to parse response: %v", err)
	}

	if len(responseBody) != len(mockFavourites) {
		t.Errorf("Expected %d favourites, got %d", len(mockFavourites), len(responseBody))
	}
}

// Mock transport for HTTP client
type mockTransport struct {
	responseBody []byte
	statusCode   int
}

func (m *mockTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	rec := httptest.NewRecorder()
	rec.WriteHeader(m.statusCode)
	rec.Write(m.responseBody)
	return rec.Result(), nil
}

// Mock transport for the HTTP client
type mockTransport2 struct {
	responses map[string]mockResponse
}

type mockResponse struct {
	body       []byte
	statusCode int
}

func (m *mockTransport2) RoundTrip(req *http.Request) (*http.Response, error) {
	fullURL := req.URL.Path
	if req.URL.RawQuery != "" {
		fullURL += "?" + req.URL.RawQuery
	}

	response, ok := m.responses[fullURL]
	if !ok {
		return nil, fmt.Errorf("unexpected request to URL: %s", fullURL)
	}

	rec := httptest.NewRecorder()
	rec.WriteHeader(response.statusCode)
	rec.Write(response.body)
	return rec.Result(), nil
}

func TestVote(t *testing.T) {
	// Prepare mock responses
	mockVoteResponse := map[string]interface{}{
		"message":  "SUCCESS",
		"id":       123456,
		"image_id": "test-image-id",
		"sub_id":   "user-7899",
		"value":    1,
	}
	mockVoteBody, _ := json.Marshal(mockVoteResponse)

	mockImageResponse := []map[string]interface{}{
		{"id": "image123", "url": "https://cdn2.thecatapi.com/images/123.jpg"},
	}
	mockImageBody, _ := json.Marshal(mockImageResponse)

	// Set up mock HTTP transport
	mockTransport2 := &mockTransport2{
		responses: map[string]mockResponse{
			"/votes": {
				body:       mockVoteBody,
				statusCode: http.StatusCreated,
			},
			"/images/search?limit=15": {
				body:       mockImageBody,
				statusCode: http.StatusOK,
			},
		},
	}

	// Replace the HTTP client globally
	httpClient := &http.Client{Transport: mockTransport2}
	oldClient := http.DefaultClient
	http.DefaultClient = httpClient
	defer func() { http.DefaultClient = oldClient }() // Restore the original HTTP client after the test

	// Prepare test payload
	payload := map[string]interface{}{
		"image_id": "test-image-id",
		"sub_id":   "user-7899",
		"value":    1,
	}
	payloadBytes, _ := json.Marshal(payload)

	// Create HTTP request and response recorder
	req := httptest.NewRequest("POST", "/api/vote", bytes.NewReader(payloadBytes))
	req.Header.Set("Content-Type", "application/json")
	resp := httptest.NewRecorder()

	// Set up controller and mock context
	controller := &controllers.CatController{}
	mockCtx := context.NewContext()
	mockCtx.Reset(resp, req)
	controller.Ctx = mockCtx

	// Execute the controller method
	controller.Vote()

	// // Parse and validate response body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	t.Logf("Response: %v", responseBody)
}
