package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	//"github.com/astaxie/beego"
	"github.com/beego/beego/v2/server/web"
	"github.com/beego/beego/v2/server/web/context"
	"github.com/stretchr/testify/assert"
)

// Mock configuration function for setting up during tests
func mockConfig() {
	// Set configuration value directly in Beego's global config
	web.BConfig.AppName = "TestApp"
	web.AppConfig.Set("catapi_key", "test-api-key")
}

func TestMain(m *testing.M) {
	// Call the mock configuration setup
	mockConfig()

	// Run the tests
	code := m.Run()

	// Exit with the result code from the tests
	os.Exit(code)
}

func TestCatController_ShowFavorites(t *testing.T) {
	tests := []struct {
		name             string
		apiResponse      string
		apiStatusCode    int
		expectedError    bool
		expectedTemplate string
		expectedCount    int
	}{
		{
			name:             "Successful favorites fetch",
			apiResponse:      `[{"id":1,"image_id":"123","image":{"url":"http://example.com/cat.jpg"}}]`,
			apiStatusCode:    http.StatusOK,
			expectedError:    false,
			expectedTemplate: "favorite.tpl",
			expectedCount:    1,
		},
		{
			name:             "API error response",
			apiResponse:      `{"error": "Internal Server Error"}`,
			apiStatusCode:    http.StatusInternalServerError,
			expectedError:    true,
			expectedTemplate: "",
			expectedCount:    0,
		},
		{
			name:             "Invalid JSON response",
			apiResponse:      `invalid json`,
			apiStatusCode:    http.StatusOK,
			expectedError:    true,
			expectedTemplate: "",
			expectedCount:    0,
		},
	}

	originalClient := http.DefaultClient

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create test server
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request headers
				assert.Equal(t, "test-api-key", r.Header.Get("x-api-key"))
				assert.Equal(t, "GET", r.Method)

				w.WriteHeader(tt.apiStatusCode)
				w.Write([]byte(tt.apiResponse))
			}))
			defer ts.Close()

			// Override the API endpoint
			http.DefaultClient = ts.Client()
			mockFavoritesData := `{"favorites": [{"id": 1, "name": "Test"}]}`
			// Setup controller
			controller := &CatController{}

			// Create a fake context
			ctx := &context.Context{
				Input:  &context.BeegoInput{},
				Output: &context.BeegoOutput{},
			}

			// Initialize the controller
			controller.Init(ctx, "CatController", "ShowFavorites", nil)
			ctx.Output.JSON(mockFavoritesData, false, false)
			// Set test configuration
			web.AppConfig.Set("catapi_key", "test-api-key")

			// Prepare mock data for successful case
			if !tt.expectedError && tt.apiStatusCode == http.StatusOK {
				var favorites []CatFavourite
				err := json.Unmarshal([]byte(tt.apiResponse), &favorites)
				assert.NoError(t, err)
				controller.Data["CatFavourites"] = favorites
				controller.TplName = tt.expectedTemplate
			}

			// Execute the function
			controller.ShowFavorites()

			// Assertions
			if tt.expectedError {
				jsonData, ok := controller.Data["json"].(map[string]string)
				assert.True(t, ok, "Expected json data to be map[string]string")
				assert.Contains(t, jsonData, "error")
			} else {
				// Verify template name
				assert.Equal(t, tt.expectedTemplate, controller.TplName)

				// Verify favorites data
				favorites, ok := controller.Data["CatFavourites"].([]CatFavourite)
				assert.True(t, ok, "Expected CatFavourites to be []CatFavourite")
				assert.Len(t, favorites, tt.expectedCount)

				if tt.expectedCount > 0 {
					// Verify the structure of the first favorite
					assert.Equal(t, 1, favorites[0].ID)
					assert.Equal(t, "123", favorites[0].ImageID)
					assert.Equal(t, "http://example.com/cat.jpg", favorites[0].Image.URL)
				}
			}
		})
	}

	// Restore the original client
	http.DefaultClient = originalClient
}
