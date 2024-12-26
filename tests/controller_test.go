package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"cat-voting-app/controllers"
	_ "cat-voting-app/controllers"

	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
)

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

// type MockConfig struct {
// 	data map[string]string
// }

// func (m *MockConfig) String(key string) (string, error) {
// 	val, ok := m.data[key]
// 	if !ok {
// 		return "", fmt.Errorf("missing config key: %s", key)
// 	}
// 	return val, nil
// }

// func (m *MockConfig) Int(key string) (int, error) {
// 	return 0, nil // Example for integer configuration
// }

// func TestFetchCatImages_MissingBaseURL(t *testing.T) {
// 	// Backup the original AppConfig
// 	originalAppConfig := web.AppConfig
// 	defer func() { web.AppConfig = originalAppConfig }() // Restore after test

// 	// Mock configuration
// 	mockConfig := &MockConfig{
// 		data: map[string]string{
// 			// "catapi_base_url": "", // Simulate missing key
// 			"catapi_key": "test_api_key",
// 		},
// 	}

// 	// Use reflection to replace web.AppConfig
// 	v := reflect.ValueOf(&web.AppConfig).Elem()
// 	v.Set(reflect.ValueOf(mockConfig))

// 	// Create a mock HTTP request
// 	req := httptest.NewRequest("GET", "/api/cats?breed_id=abc123", nil)
// 	recorder := httptest.NewRecorder()

// 	// Initialize the Beego context
// 	beegoCtx := &context.Context{
// 		Input:  context.NewInput(),
// 		Output: context.NewOutput(),
// 	}
// 	beegoCtx.Reset(recorder, req)

// 	// Set up the CatController
// 	controller := &controllers.CatController{}
// 	controller.Init(beegoCtx, "CatController", "FetchCatImages", nil)

// 	// Call the FetchCatImages method
// 	controller.FetchCatImages()

// 	// Validate the response
// 	if recorder.Code != http.StatusInternalServerError {
// 		t.Errorf("Expected status code 500, but got %d", recorder.Code)
// 	}

// 	expectedMessage := "Failed to load Cat API base URL from config"
// 	if !strings.Contains(recorder.Body.String(), expectedMessage) {
// 		t.Errorf("Expected error message '%s', but got '%s'", expectedMessage, recorder.Body.String())
// 	}
// }

type MockConfig struct {
	data map[string]string
}

func (m *MockConfig) Set(key, val string) error {
	m.data[key] = val
	return nil
}

func (m *MockConfig) String(key string) string {
	return m.data[key]
}

func (m *MockConfig) Strings(key string) []string {
	return []string{m.data[key]}
}

func (m *MockConfig) Int(key string) (int, error) {
	if val, ok := m.data[key]; ok {
		return strconv.Atoi(val)
	}
	return 0, nil
}

func (m *MockConfig) Int64(key string) (int64, error) {
	if val, ok := m.data[key]; ok {
		return strconv.ParseInt(val, 10, 64)
	}
	return 0, nil
}

func (m *MockConfig) Bool(key string) (bool, error) {
	if val, ok := m.data[key]; ok {
		return strconv.ParseBool(val)
	}
	return false, nil
}

func (m *MockConfig) Float(key string) (float64, error) {
	if val, ok := m.data[key]; ok {
		return strconv.ParseFloat(val, 64)
	}
	return 0.0, nil
}

func (m *MockConfig) DefaultString(key, defaultVal string) string {
	if val, ok := m.data[key]; ok {
		return val
	}
	return defaultVal
}

func (m *MockConfig) DefaultStrings(key string, defaultVal []string) []string {
	if val, ok := m.data[key]; ok && val != "" {
		return []string{val}
	}
	return defaultVal
}

func (m *MockConfig) DefaultInt(key string, defaultVal int) int {
	if val, ok := m.data[key]; ok {
		intVal, err := strconv.Atoi(val)
		if err == nil {
			return intVal
		}
	}
	return defaultVal
}

func (m *MockConfig) DefaultInt64(key string, defaultVal int64) int64 {
	if val, ok := m.data[key]; ok {
		int64Val, err := strconv.ParseInt(val, 10, 64)
		if err == nil {
			return int64Val
		}
	}
	return defaultVal
}

func (m *MockConfig) DefaultBool(key string, defaultVal bool) bool {
	if val, ok := m.data[key]; ok {
		boolVal, err := strconv.ParseBool(val)
		if err == nil {
			return boolVal
		}
	}
	return defaultVal
}

func (m *MockConfig) DefaultFloat(key string, defaultVal float64) float64 {
	if val, ok := m.data[key]; ok {
		floatVal, err := strconv.ParseFloat(val, 64)
		if err == nil {
			return floatVal
		}
	}
	return defaultVal
}

func (m *MockConfig) DIY(key string) (interface{}, error) {
	if val, ok := m.data[key]; ok {
		return val, nil
	}
	return nil, nil
}

func (m *MockConfig) GetSection(section string) (map[string]string, error) {
	return m.data, nil
}

func (m *MockConfig) SaveConfigFile(filename string) error {
	return nil
}

// func TestFetchCatImages_ConfigError(t *testing.T) {
// 	// Backup the original AppConfig and restore it after the test
// 	originalAppConfig := web.AppConfig
// 	defer func() { web.AppConfig = originalAppConfig }()

// 	// Mock configuration to simulate an error
// 	mockConfig := &MockConfig{
// 		data: map[string]string{
// 			"catapi_base_url": "", // Simulate missing configuration key
// 			"catapi_key":      "test_api_key",
// 		},
// 	}
// 	web.AppConfig = config.Configer(mockConfig) // Cast MockConfig to Configer interface

// 	// Create a mock HTTP request
// 	req := httptest.NewRequest("GET", "/api/cats?breed_id=abc123", nil)
// 	recorder := httptest.NewRecorder()

// 	// Initialize the Beego context
// 	beegoCtx := &context.Context{
// 		Input:  context.NewInput(),
// 		Output: context.NewOutput(),
// 	}
// 	beegoCtx.Reset(recorder, req)

// 	// Set up the CatController
// 	controller := &controllers.CatController{}
// 	controller.Init(beegoCtx, "CatController", "FetchCatImages", nil)

// 	// Call the FetchCatImages method
// 	controller.FetchCatImages()

// 	// Validate the response
// 	if recorder.Code != http.StatusInternalServerError {
// 		t.Errorf("Expected status code 500, but got %d", recorder.Code)
// 	}

// 	expectedError := "Failed to load Cat API base URL from config"
// 	if !strings.Contains(recorder.Body.String(), expectedError) {
// 		t.Errorf("Expected error message '%s', but got '%s'", expectedError, recorder.Body.String())
// 	}

// 	t.Logf("Error response: %s", recorder.Body.String())
// }

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

	// // Validate response
	// if resp.Code != http.StatusOK {
	// 	t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
	// }

	// // Parse and validate response body
	var responseBody map[string]interface{}
	if err := json.Unmarshal(resp.Body.Bytes(), &responseBody); err != nil {
		t.Fatalf("Failed to parse response body: %v", err)
	}

	// // Validate response fields
	// if responseBody["message"] != "Vote submitted and next image fetched" {
	// 	t.Errorf("Unexpected message in response: %v", responseBody["message"])
	// }

	// if _, ok := responseBody["vote"]; !ok {
	// 	t.Errorf("Missing 'vote' field in response: %v", responseBody)
	// }

	// if _, ok := responseBody["next_image"]; !ok {
	// 	t.Errorf("Missing 'next_image' field in response: %v", responseBody)
	// }

	t.Logf("Response: %v", responseBody)
}
