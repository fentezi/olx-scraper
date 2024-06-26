package main

import (
	"context"
	"fmt"
	"log/slog"
	"math/rand"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
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
            var phone string
            ctx, cancel := chromedp.NewContext(
            context.Background(),
        )
            linkSelector := `a[data-testid="contact-phone"]`
            _ = chromedp.Run(ctx, 
                chromedp.Navigate("https://www.olx.ua" + published.HrefPublished),
                chromedp.Sleep(5 * time.Second),
                chromedp.Click(`button[data-cy="ad-contact-phone"]`, chromedp.ByQuery),
                chromedp.TextContent(linkSelector, &phone),
            )
            cancel()
            SendMessagePhoto(b, id, &published, phone)
            log.Info(
                "sent message",
                "id", id,
                "title", published.Title,
                "price", published.Price,
                "city", published.City,
                "time", published.TimePublished,
                "phone", phone,
                "href", published.HrefPublished,
            )
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

func SendMessagePhoto(b *tele.Bot, id int64, publish *models.Published, phone string) {
    splitPhoto := strings.Split(publish.Image, ";")
    var photo *tele.Photo
    if phone == "" {
        photo = &tele.Photo{
            File: tele.FromURL(splitPhoto[0]),
            Caption: publish.Title + "\n" +
                "Цена: " + publish.Price + "\n" +
                "Город: " + publish.City + "\n" +
                "Время публикации: " + publish.TimePublished,
    }
    } else {
        photo = &tele.Photo{
        File: tele.FromURL(splitPhoto[0]),
        Caption: publish.Title + "\n" +
            "Цена: " + publish.Price + "\n" +
            "Город: " + publish.City + "\n" +
            "Номер телефона: " + phone + "\n" +
            "Время публикации: " + publish.TimePublished,
    }
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
        return ctx.Send(fmt.Sprintf("Привет, %s! Я бот для парсинга объявлений на OLX. Чтобы начать, нажмите команду /addurl", ctx.Sender().Username))
    })

    b.Handle("/addurl", func(ctx tele.Context) error {
        userID := ctx.Sender().ID
        urlState[userID] = true
		message := "Введите URL-адрес для парсинга." +
			"\n\nНа OLX выберите нужный город, категорию товаров, фильтры и параметры поиска. " +
			 "После того, как все критерии заданы, скопируйте URL-адрес из адресной строки браузера и отправьте его боту." +
			 "\n\nПример подходящей ссылки: https://www.olx.ua/uk/nedvizhimost/kvartiry/prodazha-kvartir/"

        return ctx.Send(message)
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
                return ctx.Send("URL успешно добавлен.\n\nКак только подходящие объявления появятся, бот оповестит вас.")
            } else {
                return ctx.Send("Неверный URL. Пожалуйста, попробуйте еще раз.")
            }
        }
        return ctx.Send(ctx.Text())
    })

    b.Handle("/stopparse", func(ctx tele.Context) error {
        userID := ctx.Sender().ID
        if stopChan, ok := stopChannels[userID]; ok {
            stopChan <- struct{}{}
            return ctx.Send("Парсер остановлен!")
        } else {
			return ctx.Send("Парсинг еще не начат!")
		}
    })

	b.Handle("/help", func(ctx tele.Context) error {
		helpMessage := "Доступные команды:\n/addurl - добавить ссылку на категорию или поиск OLX для отслеживания новых объявлений. Бот будет парсить указанную страницу и оповещать о свежих публикациях\n/stopparse - остановить парсинг\n/help - помощь"
		return ctx.Send(helpMessage)
	})

    b.Start()
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
