package main

import (
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/fentezi/olx-scraper/internal"
	"github.com/fentezi/olx-scraper/logger"
	"github.com/fentezi/olx-scraper/models"
	"github.com/fentezi/olx-scraper/utils"
	tele "gopkg.in/telebot.v3"
)

var (
    urlState     = make(map[int64]bool)
    urls         = make(map[int64]string)
    log          *slog.Logger
    stopChannels = make(map[int64]chan struct{})
)

func main() {
    log = logger.Logger()
    log.Info("application started")
    token := os.Getenv("TOKEN")
    if token == "" {
        log.Error("missing token")
        os.Exit(1)
    }

    pref := tele.Settings{
        Token:   token,
        Poller:  &tele.LongPoller{Timeout: 10 * time.Second},
        Verbose: false,
    }
    b, err := tele.NewBot(pref)
    if err != nil {
        os.Exit(1)
    }
    TelegramInit(b)
}

func parseLoop(id int64, b *tele.Bot) {
    var currentTime time.Time
    stopChan := stopChannels[id]
	log.Info(
		"stop channels",
		"len", len(stopChannels),
	)
    for {
        log := logger.Logger()
        log.Info("fetching ads", "id", id)

        doc, err := utils.FetchAndParseHTML(urls[id])
        if err != nil {
            log.Error("failed to fetch and parse HTML", logger.Err(err))
            continue
        }

        published := internal.GetPublished(doc, *log)
        if utils.ShouldPrintPublished(&published, currentTime) {
            SendMessagePhoto(b, id, &published)
            printPublished(published)
            err = utils.SaveImage(published.Image, *log)
            if err != nil {
                log.Warn("failed to save image", logger.Err(err))
            }
        }

        timeParse, _ := time.Parse("15:04", published.TimePublished)
        currentTime = timeParse

        select {
        case <-stopChan:
            log.Info(
                "stopping parsing",
                "id", id,
            )
			delete(stopChannels, id)
            return
        default:
            sleepDuration := time.Duration(rand.Intn(120)) * time.Second
            timer := time.NewTimer(sleepDuration)
            defer timer.Stop()

            select {
            case <-stopChan:
                log.Info(
                    "stopping parsing",
                    "id", id,
                )
				delete(stopChannels, id)
                return
            case <-timer.C:
            }
        }
    }
}

func SendMessagePhoto(b *tele.Bot, id int64, publish *models.Published) {
    splitPhoto := strings.Split(publish.Image, ";")
    photo := &tele.Photo{
        File: tele.FromURL(splitPhoto[0]),
        Caption: publish.Title + "\n" +
            "Цена: " + publish.Price + "\n" +
            "Город: " + publish.City + "\n" +
            "Время публикации: " + publish.TimePublished,
    }

    btnURL := tele.InlineButton{
        Unique: "myButton",
        Text:   "Объявление",
        URL:    "https://www.olx.ua" + publish.HrefPublished,
    }

    row := []tele.InlineButton{btnURL}

    b.Send(tele.ChatID(id), photo, &tele.SendOptions{
        ReplyMarkup: &tele.ReplyMarkup{
            InlineKeyboard: [][]tele.InlineButton{row},
        },
    })
}

func TelegramInit(b *tele.Bot) {
    b.Handle("/start", func(ctx tele.Context) error {
        userID := ctx.Sender().ID
        stopChannels[userID] = make(chan struct{})
        return ctx.Send("Hello, World!")
    })

    b.Handle("/addURL", func(ctx tele.Context) error {
        userID := ctx.Sender().ID
        urlState[userID] = true
        return ctx.Send("Enter the URL:")
    })

    b.Handle(tele.OnText, func(ctx tele.Context) error {
        userID := ctx.Sender().ID
        if urlState[userID] {
            text := ctx.Text()
            if isValidURL(text) {
                urls[userID] = text
                go parseLoop(userID, b)
                urlState[userID] = false
				stopChannels[userID] = make(chan struct{})
                return ctx.Send("URL saved: " + text)
            } else {
                return ctx.Send("Invalid URL. Please try again.")
            }
        }
        return ctx.Send(ctx.Text())
    })

    b.Handle("/stop", func(ctx tele.Context) error {
        userID := ctx.Sender().ID
        if stopChan, ok := stopChannels[userID]; ok {
            stopChan <- struct{}{}
            return ctx.Send("Parsing stopped")
        } else {
			return ctx.Send("Parsing not started yet")
		}
    })

    b.Start()
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

func isValidURL(str string) bool {
	pattern := `^(https?://)?(www\.)?olx\.(ua|pl|bg|ro|pt|com|co\.za|com\.br|com\.pk|lt|lv|hr|kz|uz|by|md|az)/.*$`
	match, _ := regexp.Match(pattern, []byte(str))
	if match {
		_, err := url.ParseRequestURI(str)
		return err == nil
	} else {
		return false
	}
}