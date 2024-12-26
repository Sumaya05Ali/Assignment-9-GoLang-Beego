package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	//"strings"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type CatController struct {
	beego.Controller
}

// Struct to handle the response from different APIs
type CatImage struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

type CatBreed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Temperament string `json:"temperament"`
	Origin      string `json:"origin"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
}

type CatVote struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
	Value   int    `json:"value"`
}

type CatFavourite struct {
	ID      int    `json:"id"`
	ImageID string `json:"image_id"`
	Image   struct {
		URL string `json:"url"`
	} `json:"image"`
}

func (c *CatController) Get() {
	apiKey, _ := beego.AppConfig.String("catapi_key")
	apiURL, _ := beego.AppConfig.String("catapi_url")

	// Create channels for concurrent API calls
	responseChannel := make(chan []CatImage)
	breedsChannel := make(chan []CatBreed)
	votesChannel := make(chan []CatVote)
	favouritesChannel := make(chan []CatFavourite)

	// Fetch images in a separate goroutine
	go func() {
		client := &http.Client{Timeout: time.Second * 10}
		req, err := http.NewRequest("GET", apiURL, nil)
		if err != nil {
			fmt.Println("Error creating request for cat images:", err)
			close(responseChannel)
			return
		}
		req.Header.Add("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching cat images:", err)
			close(responseChannel)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body for cat images:", err)
			close(responseChannel)
			return
		}

		var images []CatImage
		if err := json.Unmarshal(body, &images); err != nil {
			fmt.Println("Error parsing cat images JSON:", err)
			close(responseChannel)
			return
		}
		responseChannel <- images
	}()

	// Fetch breeds in a separate goroutine
	go func() {
		client := &http.Client{Timeout: time.Second * 10}
		req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
		if err != nil {
			fmt.Println("Error creating request for cat breeds:", err)
			close(breedsChannel)
			return
		}
		req.Header.Add("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching cat breeds:", err)
			close(breedsChannel)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body for cat breeds:", err)
			close(breedsChannel)
			return
		}

		var breeds []CatBreed
		if err := json.Unmarshal(body, &breeds); err != nil {
			fmt.Println("Error parsing cat breeds JSON:", err)
			close(breedsChannel)
			return
		}
		breedsChannel <- breeds
	}()

	// Fetch votes in a separate goroutine
	go func() {
		client := &http.Client{Timeout: time.Second * 10}
		req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/votes", nil)
		if err != nil {
			fmt.Println("Error creating request for cat votes:", err)
			close(votesChannel)
			return
		}
		req.Header.Add("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching cat votes:", err)
			close(votesChannel)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body for cat votes:", err)
			close(votesChannel)
			return
		}

		var votes []CatVote
		if err := json.Unmarshal(body, &votes); err != nil {
			fmt.Println("Error parsing cat votes JSON:", err)
			close(votesChannel)
			return
		}
		votesChannel <- votes
	}()

	// Fetch favourites in a separate goroutine
	go func() {
		client := &http.Client{Timeout: time.Second * 10}
		req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/favourites", nil)
		if err != nil {
			fmt.Println("Error creating request for cat favourites:", err)
			close(favouritesChannel)
			return
		}
		req.Header.Add("x-api-key", apiKey)

		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("Error fetching cat favourites:", err)
			close(favouritesChannel)
			return
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body for cat favourites:", err)
			close(favouritesChannel)
			return
		}

		var favourites []CatFavourite
		if err := json.Unmarshal(body, &favourites); err != nil {
			fmt.Println("Error parsing cat favourites JSON:", err)
			close(favouritesChannel)
			return
		}
		favouritesChannel <- favourites
	}()

	// Collect results from all channels
	var catImages []CatImage
	var catBreeds []CatBreed
	var catVotes []CatVote
	var catFavourites []CatFavourite

	select {
	case catImages = <-responseChannel:
	default:
		fmt.Println("No data received from cat images API")
	}

	select {
	case catBreeds = <-breedsChannel:
	default:
		fmt.Println("No data received from cat breeds API")
	}

	select {
	case catVotes = <-votesChannel:
	default:
		fmt.Println("No data received from cat votes API")
	}

	select {
	case catFavourites = <-favouritesChannel:
	default:
		fmt.Println("No data received from cat favourites API")
	}

	// Send the data to the template
	c.Data["CatImages"] = catImages
	c.Data["CatBreeds"] = catBreeds
	c.Data["CatVotes"] = catVotes
	c.Data["CatFavourites"] = catFavourites
	c.TplName = "cat.tpl"
}
func (c *CatController) SaveFavorite() {
	// Get API key from Beego configuration
	apiKey, err := beego.AppConfig.String("catapi_key")
	if err != nil {
		c.Data["json"] = map[string]string{"error": "API key not configured"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}

	// Parse request body
	var requestBody struct {
		ImageID string `json:"image_id"`
	}

	// Read the request body
	body, err := io.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read request body"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}
	// Log the received body for debugging
	fmt.Printf("Received request body: %s\n", string(body))

	if err := json.Unmarshal(body, &requestBody); err != nil {
		c.Data["json"] = map[string]string{"error": "Invalid JSON payload"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	if requestBody.ImageID == "" {
		c.Data["json"] = map[string]string{"error": "Image ID is required"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		c.ServeJSON()
		return
	}

	// Make request to The Cat API
	client := &http.Client{Timeout: time.Second * 10}
	apiURL := "https://api.thecatapi.com/v1/favourites"

	// Create request body
	requestBodyData, err := json.Marshal(map[string]string{
		"image_id": requestBody.ImageID,
	})
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to create request body"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}

	// Create request
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBodyData))
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to create request"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}

	// Add headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", apiKey)

	// Make the request
	resp, err := client.Do(req)
	if err != nil {
		c.Data["json"] = map[string]string{"error": fmt.Sprintf("Failed to make request: %v", err)}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	// Read response body
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		c.Data["json"] = map[string]string{"error": "Failed to read response body"}
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		c.ServeJSON()
		return
	}

	// Check response status
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		c.Data["json"] = map[string]string{
			"error": fmt.Sprintf("API error (status %d): %s", resp.StatusCode, string(responseBody)),
		}
		c.Ctx.ResponseWriter.WriteHeader(resp.StatusCode)
		c.ServeJSON()
		return
	}

	// Return success response
	c.Data["json"] = map[string]string{
		"message":  "Favorite saved successfully",
		"image_id": requestBody.ImageID,
	}
	c.ServeJSON()
}

func (c *CatController) ShowFavorites() {
	apiKey, _ := beego.AppConfig.String("catapi_key")

	// Fetch favorites
	client := &http.Client{Timeout: time.Second * 10}
	req, _ := http.NewRequest("GET", "https://api.thecatapi.com/v1/favourites", nil)
	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error fetching favourites:", err)
		c.Data["json"] = map[string]string{"error": "Failed to fetch favorites"}
		c.ServeJSON()
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var favourites []CatFavourite
	if err := json.Unmarshal(body, &favourites); err != nil {
		fmt.Println("Error parsing favourites JSON:", err)
		c.Data["json"] = map[string]string{"error": "Failed to parse favorites"}
		c.ServeJSON()
		return
	}

	// Pass data to the template
	c.Data["CatFavourites"] = favourites
	c.TplName = "favorite.tpl"
}
