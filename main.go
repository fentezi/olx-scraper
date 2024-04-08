package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"time"

	"github.com/fentezi/olx-scraper/internal"
	"github.com/fentezi/olx-scraper/logger"
	"github.com/fentezi/olx-scraper/models"
	"github.com/fentezi/olx-scraper/utils"
)

const URL = "https://www.olx.ua/uk/nedvizhimost/kvartiry/dolgosrochnaya-arenda-kvartir/dnepr/?currency=UAH&search%5Border%5D=created_at:desc&view=list"

var log *slog.Logger

// main is the entry point of the program.
// It continuously fetches and prints the published ads
// that have a time greater than the current time.
func main() {
	var currentTime time.Time
	log = logger.Logger()

	log.Info(
		"application started",
	)
	// Enter an infinite loop to continuously check for new ads
	for {
		// Fetch and parse the HTML content of the URL
		doc, err := utils.FetchAndParseHTML(URL)
		if err != nil {
			// If failed to fetch and parse the HTML, log the error and exit
			log.Error(
				"fetch and parse HTML",
				logger.Err(err),
			)
		}
		// Get the published ad from the HTML document
		published := internal.GetPublished(doc, *log)
		// Check if the published ad should be printed
		if utils.ShouldPrintPublished(&published, currentTime) {
			// Print the published
			printPublished(published)
			err = utils.SaveImage(published.Image, *log)
			if err != nil {
				log.Warn(
					"save image",
					logger.Err(err),
				)
			}
			log.Info(
				"published ad printed",
			)
		}

		timeParse, _ := time.Parse("15:04", published.TimePublished)
		currentTime = timeParse
		// Wait for seconds before checking for new ads
		time.Sleep(time.Duration(rand.Intn(120)) * time.Second)
	}
}

// printPublished prints the published ad details to the console.
// It takes a Published struct as input and prints the title, URL, and current time.
func printPublished(published models.Published) {
	// Print the title of the published ad
	fmt.Println("Публикация: " + published.Title)

	// Print the URL Image of the published ad
	fmt.Println("Фото: " + published.Image)

	// Print the URL of the published ad
	fmt.Println("Cсылка на объявление: https://www.olx.ua" + published.HrefPublished)

	// Print the price of the published ad
	fmt.Println("Цена: " + published.Price)

	// Print the city of the published ad
	fmt.Println("Город: " + published.City)

	// Print the current time in the format "15:04"
	fmt.Println("Время публикации: " + published.TimePublished)
}