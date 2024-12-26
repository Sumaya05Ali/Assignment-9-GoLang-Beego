package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	beego "github.com/beego/beego/v2/server/web"
)

type BreedsController struct {
	beego.Controller
}

type Breed struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Temperament string `json:"temperament"`
	Origin      string `json:"origin"`
	Description string `json:"description"`
	Image       struct {
		URL string `json:"url"`
	} `json:"image"`
}

type Image struct {
	ID     string `json:"id"`
	URL    string `json:"url"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
}

func fetchBreeds(apiKey string, ch chan<- []Breed, errCh chan<- error) {
	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", "https://api.thecatapi.com/v1/breeds", nil)
	if err != nil {
		errCh <- fmt.Errorf("error creating request: %v", err)
		return
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("error fetching breeds: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- fmt.Errorf("error reading response body: %v", err)
		return
	}

	var breeds []Breed
	if err := json.Unmarshal(body, &breeds); err != nil {
		errCh <- fmt.Errorf("error parsing JSON: %v", err)
		return
	}

	ch <- breeds
}

func fetchBreedImages(breedID string, apiKey string, ch chan<- []Image, errCh chan<- error) {
	client := &http.Client{Timeout: 10 * time.Second}
	url := fmt.Sprintf("https://api.thecatapi.com/v1/images/search?breed_ids=%s&limit=5", breedID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		errCh <- fmt.Errorf("error creating request: %v", err)
		return
	}

	req.Header.Add("x-api-key", apiKey)

	resp, err := client.Do(req)
	if err != nil {
		errCh <- fmt.Errorf("error fetching images: %v", err)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		errCh <- fmt.Errorf("error reading response body: %v", err)
		return
	}

	var images []Image
	if err := json.Unmarshal(body, &images); err != nil {
		errCh <- fmt.Errorf("error parsing JSON: %v", err)
		return
	}

	ch <- images
}

func (b *BreedsController) Get() {
	b.Layout = ""
	b.TplName = "breeds.tpl"
}

func (b *BreedsController) GetBreeds() {
	apiKey, _ := beego.AppConfig.String("catapi_key")

	breedsCh := make(chan []Breed)
	errCh := make(chan error)

	go fetchBreeds(apiKey, breedsCh, errCh)

	select {
	case breeds := <-breedsCh:
		b.Data["json"] = breeds
		b.ServeJSON()
	case err := <-errCh:
		b.Data["json"] = map[string]string{"error": err.Error()}
		b.ServeJSON()
	}
}

func (b *BreedsController) GetBreedImages() {
	breedID := b.Ctx.Input.Param(":breed_id")
	apiKey, _ := beego.AppConfig.String("catapi_key")

	imagesCh := make(chan []Image)
	errCh := make(chan error)

	go fetchBreedImages(breedID, apiKey, imagesCh, errCh)

	select {
	case images := <-imagesCh:
		b.Data["json"] = images
		b.ServeJSON()
	case err := <-errCh:
		b.Ctx.Output.SetStatus(500)
		b.Data["json"] = map[string]string{"error": err.Error()}
		b.ServeJSON()
	case <-time.After(15 * time.Second):
		b.Ctx.Output.SetStatus(504)
		b.Data["json"] = map[string]string{"error": "request timed out"}
		b.ServeJSON()
	}
}
