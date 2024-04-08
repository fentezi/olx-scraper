package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

const URL = "https://www.olx.ua/uk/nedvizhimost/kvartiry/dolgosrochnaya-arenda-kvartir/dnepr/?currency=UAH&search%5Border%5D=created_at:desc&view=list"

type Published struct {
	Title         string
	Image         string
	City          string
	Price         string
	HrefPublished string
	TimePublished string
}

// main is the entry point of the program.
// It continuously fetches and prints the published ads
// that have a time greater than the current time.
func main() {
	var currentTime time.Time

	// Enter an infinite loop to continuously check for new ads
	for {
		// Fetch and parse the HTML content of the URL
		doc, err := fetchAndParseHTML(URL)
		if err != nil {
			// If failed to fetch and parse the HTML, log the error and exit
			log.Fatal(err)
		}
		// Get the published ad from the HTML document
		published := getPublished(doc)
		// Check if the published ad should be printed
		if shouldPrintPublished(&published, currentTime) {
			// Print the published
			printPublished(published)
			err = saveImage(published.Image)
			if err != nil {
				log.Println(err)
			}
		}

		timeParse, _ := time.Parse("15:04", published.TimePublished)
		currentTime = timeParse
		// Wait for seconds before checking for new ads
		time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	}
}

<<<<<<< HEAD
// saveImage downloads an image from the provided URL and saves it to the "image.png" file.
//
// Parameters:
// - url: a string representing the URL of the image to be downloaded.
//
// Returns:
// - error: an error if the download or saving process fails.
func saveImage(url string) error {
	// Check if the URL is empty.
	if url == "" {
		return errors.New("image url is empty")
	}

	// Send a GET request to the URL and retrieve the response.
=======
func saveImage(url string) error {
	if url == "" {
		return errors.New("image url is empty")
	}
>>>>>>> e9aa00d7a89ad798ac39d59daaa8cd0cdabdc067
	resp, err := http.Get(url)
	if err != nil {
		return errors.New("failed to get image")
	}
	defer resp.Body.Close()

<<<<<<< HEAD
	// Read the body of the response.
=======
>>>>>>> e9aa00d7a89ad798ac39d59daaa8cd0cdabdc067
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
<<<<<<< HEAD

	// Write the body to the "image.png" file.
=======
>>>>>>> e9aa00d7a89ad798ac39d59daaa8cd0cdabdc067
	err = os.WriteFile("image.png", body, 0644)
	if err != nil {
		return err
	}
<<<<<<< HEAD

	// Log the success message.
	log.Println("Image saved to image.png")

	// Return nil indicating success.
=======
	log.Println("Image saved to image.png")
>>>>>>> e9aa00d7a89ad798ac39d59daaa8cd0cdabdc067
	return nil
}

// getPublished fetches and returns the published ad from the HTML document.
//
// It takes a goquery.Document pointer as input and returns a Published struct.
// It calls the returnPublished function to fetch and parse the HTML content.
func getPublished(doc *goquery.Document) Published {
	// Call the returnPublished function to fetch and parse the HTML content.
	// It returns a Published struct.
	return returnPublished(doc)
}

// shouldPrintPublished checks if the published ad should be printed based on its time and the current time.
// It takes a Published struct and the current time as input and returns a boolean.
// It parses the time of the published ad and checks if it is after the current time.
// If the published ad is after the current time, it should be printed and the function returns true.
func shouldPrintPublished(published *Published, currentTime time.Time) bool {
	// Parse the time of the published ad
	timeParse, _ := time.Parse("15:04", published.TimePublished)
	timeParse = timeParse.Add(3 * time.Hour)
	published.TimePublished = timeParse.Format("15:04")
	// Check if the published ad is after the current time
	return timeParse.After(currentTime)
}

// printPublished prints the published ad details to the console.
// It takes a Published struct as input and prints the title, URL, and current time.
func printPublished(published Published) {
	// Print the title of the published ad
	fmt.Println("Публикация: " + published.Title)

	// Print the URL Image of the published ad
	fmt.Println("Фото: " + published.Image)

	// Print the URL of the published ad
	fmt.Println("Cсылка на объявление: https://www.olx.ua/" + published.HrefPublished)

	// Print the price of the published ad
	fmt.Println("Цена: " + published.Price)

	// Print the citi of the published ad
	fmt.Println("Город: " + published.City)

	// Print the current time in the format "15:04"
	fmt.Println("Время публикации:" + published.TimePublished)
}

// fetchAndParseHTML fetches and parses the HTML content of a given URL.
// It uses the provided URL to create a new HTTP request and retrieves the
// HTML content of the URL. It returns a goquery.Document and an error if
// any occurred.
func fetchAndParseHTML(url string) (*goquery.Document, error) {
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

// printAdTitles prints the titles and times of ads in the document.
// It checks if the document and selection are not nil, and if the title
// and time of each ad are not nil. If any of these checks fail, it exits
// the program with an error.
func returnPublished(doc *goquery.Document) Published {
	var titleText, href, timeText, price, urlImage, city string
	// Check if the document is nil
	if doc == nil {
		logAndExit(fmt.Errorf("document is nil"))
	}

	// Find the selection of ads
	selection := doc.Find("div#div-gpt-liting-after-promoted").Next()

	// Check if the selection is nil
	if selection == nil {
		logAndExit(fmt.Errorf("selection is nil"))
	}

	// Iterate over each ad and print its title and time
	selection.Each(func(i int, s *goquery.Selection) {
		// Find the title of the ad
		title := s.Find("h6")

		// Check if the title is nil
		if title == nil {
			logAndExit(fmt.Errorf("title is nil"))
		}

		// Get the text of the title
		titleText = title.Text()

		// Find the time of the ad
		timeAttr := s.Find(`p[data-testid="location-date"]`)

		// Check if the time is nil
		if timeAttr == nil {
			logAndExit(fmt.Errorf("time is nil"))
		}
		price = s.Find(`p[data-testid="ad-price"]`).Text()
		href, _ = s.Find(`a.css-z3gu2d`).Attr("href")

		// Get the text of the time
		timeSplit := strings.Split(timeAttr.Text(), " ")
		timeText = timeSplit[len(timeSplit)-1]

		urlImage, _ = s.Find(`div.css-gl6djm > img`).Attr("src")
		city = timeSplit[0]

	})
	return Published{
		Title:         titleText,
		Image:         urlImage,
		Price:         price,
		City:          city[:len(city)-1],
		HrefPublished: href,
		TimePublished: timeText,
	}
}

func logAndExit(err error) {
	fmt.Println(err)
	os.Exit(1)
}

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
