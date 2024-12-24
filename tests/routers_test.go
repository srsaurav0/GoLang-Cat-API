package test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/beego/beego/v2/server/web"
	"github.com/stretchr/testify/assert"

	_ "cat-voting-app/routers"
)

func TestRouters(t *testing.T) {
	// Initialize the Beego app and router
	web.TestBeegoInit("./") // Adjust the path to your app's root

	tests := []struct {
		method   string
		route    string
		expected string
	}{
		{"GET", "/", "MainController"},
		{"GET", "/api/cats", "CatController.FetchCatImages"},
		{"GET", "/api/breeds", "CatController.FetchBreeds"},
		{"POST", "/api/add-to-favourites", "CatController.AddToFavourites"},
		{"GET", "/api/get-favourites", "CatController.GetFavourites"},
		{"DELETE", "/api/remove-favourite", "CatController.RemoveFavourite"},
		{"POST", "/api/vote", "CatController.Vote"},
		{"GET", "/api/votes", "CatController.GetVotes"},
	}

	for _, test := range tests {
		t.Run(test.route, func(t *testing.T) {
			// Create a new HTTP request
			req, _ := http.NewRequest(test.method, test.route, nil)
			resp := httptest.NewRecorder()

			// Serve the request
			web.BeeApp.Handlers.ServeHTTP(resp, req)

			// Verify the response status code
			assert.NotEqual(t, http.StatusNotFound, resp.Code, "Route should not return 404")
		})
	}
}

func TestRoutes(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/cats", nil)
	recorder := httptest.NewRecorder()
	web.BeeApp.Handlers.ServeHTTP(recorder, req)

	if recorder.Code != http.StatusOK {
		t.Errorf("Expected status OK but got %v", recorder.Code)
	}
}
