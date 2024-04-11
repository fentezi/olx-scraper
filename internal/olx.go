package internal

import (
	"log/slog"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/fentezi/olx-scraper/models"
)

// printAdTitles prints the titles and times of ads in the document.
// It checks if the document and selection are not nil, and if the title
// and time of each ad are not nil. If any of these checks fail, it exits
// the program with an error.
func returnPublished(doc *goquery.Document, log *slog.Logger) models.Published {
	// Check if the document is nil
	if doc == nil {
		log.Error("document is nil")
		return models.Published{}
	}

	// Find the selection of ads
	selection := doc.Find("div#div-gpt-liting-after-promoted").Next()

	if selection == nil {
		log.Error("selection is nil")

	}

	// Find the title of the ad
	title := selection.Find("h6")

	// Check if the title is nil
	if title == nil {
		log.Warn("title is nil")
	}

	// Get the text of the title
	titleText := title.Text()

	// Find the time of the ad
	timeAttr := selection.Find(`p[data-testid="location-date"]`)

	// Check if the time is nil
	if timeAttr == nil {
		log.Warn("time is nil")
	}
	price := strings.Replace(selection.Find(`p[data-testid='ad-price'].css-tyui9s.er34gjf0`).Text(), ".css-1vxklie{color:#7F9799;font-size:12px;line-height:16px;font-weight:100;display:block;width:100%;text-align:right;}Договірна", "Договірна", -1)
	href, _ := selection.Find(`a.css-z3gu2d`).Attr("href")

	// Get the text of the time
	timeSplit := strings.Split(timeAttr.Text(), " ")
	timeText := timeSplit[len(timeSplit)-1]
	

	urlImage, _ := selection.Find(`div.css-gl6djm > img`).Attr("src")
	log.Info(
		"timeSplit",
		"text", price,
	)
	citySplit := strings.Split(timeAttr.Text(), " - ")
	city := citySplit[0]

	return models.Published{
		Title:         titleText,
		Image:         urlImage,
		Price:         price,
		City:          city,
		HrefPublished: href,
		TimePublished: timeText,
	}
}

// getPublished fetches and returns the published ad from the HTML document.
//
// It takes a goquery.Document pointer as input and returns a Published struct.
// It calls the returnPublished function to fetch and parse the HTML content.
func GetPublished(doc *goquery.Document, log slog.Logger) models.Published {
	// Call the returnPublished function to fetch and parse the HTML content.
	// It returns a Published struct.
	return returnPublished(doc, &log)
}