package utils

import (
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/fentezi/olx-scraper/models"
)

// fetchHTML retrieves the HTML content of a given URL.
// It uses the provided URL to create a new HTTP request.
// The function sets the User-Agent header to mimic a web browser.
// It sends the request and retrieves the response.
// If the response status code is not 200 OK, it returns an error.
// Otherwise, it returns the response body as an io.ReadCloser.
func fetchHTML(url string) (io.ReadCloser, error) {
	// Create a new HTTP client
	client := http.Client{}

	// Create a new HTTP request with the provided URL
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// If failed to create the request, return an error
		return nil, fmt.Errorf("failed to create request: %v", err)
	}

	// Set the User-Agent header to mimic a web browser
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36")
	// Send the request and retrieve the response
	resp, err := client.Do(req)
	if err != nil {
		// If failed to send the request, return an error
		return nil, fmt.Errorf("failed to send request: %v", err)
	}

	// If the response status code is not 200 OK, return an error
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch %s: %s", url, resp.Status)
	}

	// Return the response body as an io.ReadCloser
	return resp.Body, nil
}

// fetchAndParseHTML fetches and parses the HTML content of a given URL.
// It uses the provided URL to create a new HTTP request and retrieves the
// HTML content of the URL. It returns a goquery.Document and an error if
// any occurred.
func FetchAndParseHTML(url string) (*goquery.Document, error) {
	// Fetch the HTML content of the URL
	html, err := fetchHTML(url)
	if err != nil {
		// If failed to fetch the HTML, return the error
		return nil, err
	}
	defer html.Close() // Defer the closing of the HTML content

	// Parse the HTML content and return the goquery.Document
	return goquery.NewDocumentFromReader(html)
}

// saveImage downloads an image from the provided URL and saves it to the "image.png" file.
//
// Parameters:
// - url: a string representing the URL of the image to be downloaded.
//
// Returns:
// - error: an error if the download or saving process fails.
func SaveImage(url string, log slog.Logger) error {
	// Check if the URL is empty.
	if url == "" {
		return errors.New("image url is empty")
	}

	// Send a GET request to the URL and retrieve the response.
	resp, err := http.Get(url)
	if err != nil {
		return errors.New("failed to get image")
	}
	defer resp.Body.Close()

	// Read the body of the response.
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Write the body to the "image.png" file.
	err = os.WriteFile("images/image.png", body, 0644)
	if err != nil {
		return err
	}

	// Log the success message.
	log.Info(
		"saveImage",
		slog.String("imageInfo", "image saved to images/image.png"),
	)
	// Return nil indicating success.
	return nil
}

// shouldPrintPublished checks if the published ad should be printed based on its time and the current time.
// It takes a Published struct and the current time as input and returns a boolean.
// It parses the time of the published ad and checks if it is after the current time.
// If the published ad is after the current time, it should be printed and the function returns true.
func ShouldPrintPublished(published *models.Published, currentTime time.Time) bool {
	// Parse the time of the published ad
	timeParse, _ := time.Parse("15:04", published.TimePublished)
	timeParse = timeParse.Add(3 * time.Hour)
	published.TimePublished = timeParse.Format("15:04")
	// Check if the published ad is after the current time
	return timeParse.After(currentTime)
}