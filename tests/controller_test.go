package test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
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

// func TestAddToFavourites(t *testing.T) {
//     payload := map[string]string{
//         "image_id": "test-image-id", // Must not be empty
//         "sub_id":   "user-123",
//     }
//     payloadBytes, err := json.Marshal(payload)
//     if err != nil {
//         t.Fatalf("Failed to marshal payload: %v", err)
//     }

//     req := httptest.NewRequest("POST", "/api/add-to-favourites", bytes.NewReader(payloadBytes))
//     req.Header.Set("Content-Type", "application/json")

//     resp := httptest.NewRecorder()

//     controller := &controllers.CatController{}
//     mockCtx := context.NewContext()
//     mockCtx.Reset(req, resp)
//     controller.Ctx = mockCtx

//     controller.AddToFavourites()

//     if resp.Code != http.StatusOK {
//         t.Errorf("Expected status code %d, got %d", http.StatusOK, resp.Code)
//     }

//     var responseBody map[string]interface{}
//     err = json.Unmarshal(resp.Body.Bytes(), &responseBody)
//     if err != nil {
//         t.Fatalf("Failed to parse response body: %v", err)
//     }

//     t.Logf("Response: %v", responseBody)
// }

// func TestAddToFavouritesError(t *testing.T) {
// 	// Mock request body
// 	requestBody := map[string]string{
// 		"image_id": "invalid",
// 		"sub_id":   "user-123",
// 	}
// 	body, _ := json.Marshal(requestBody)

// 	// Create a mock request
// 	req := httptest.NewRequest(http.MethodPost, "/api/add-to-favourites", bytes.NewReader(body))
// 	req.Header.Set("Content-Type", "application/json")

// 	// Create a mock response recorder
// 	resp := httptest.NewRecorder()

// 	// Initialize the controller
// 	controller := &controllers.CatController{}
// 	controller.Ctx = &context.Context{
// 		Input: &context.BeegoInput{
// 			Context: &context.Context{
// 				Request: req,
// 			},
// 		},
// 		Output: &context.BeegoOutput{
// 			Context: &context.Context{
// 				ResponseWriter: resp,
// 			},
// 		},
// 	}

// 	// Mock helper functions
// 	controllers.AddToFavorites = func(baseURL, apiKey, imageID, subID string) error {
// 		return assert.AnError
// 	}

// 	controllers.FetchCatImages = func(baseURL, apiKey, breedID string) ([]map[string]interface{}, error) {
// 		return nil, assert.AnError
// 	}

// 	// Call the AddToFavourites method
// 	controller.AddToFavourites()

// 	// Assert the response
// 	assert.Equal(t, http.StatusInternalServerError, resp.Code)

// 	// Parse response body
// 	var responseBody map[string]interface{}
// 	err := json.NewDecoder(resp.Body).Decode(&responseBody)
// 	assert.NoError(t, err)

// 	// Validate the error response
// 	assert.Equal(t, "Failed to complete tasks", responseBody["error"])
// 	assert.NotNil(t, responseBody["details"])
// }

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
