# OLX Parser Telegram Bot

This is a Telegram bot that parses and sends notifications about the latest ads published on the OLX website. The bot scrapes the specified OLX category or search URL and sends an alert with the ad details, including an image, title, price, city, and a direct link to the ad.

## Features

- Parse and monitor new ads on OLX based on user-provided URL
- Send notifications with ad details (image, title, price, city, direct link)
- Inline keyboard button to open the ad URL
- Stop parsing command to pause monitoring

## Installation

1. Clone the repository:

```bash
git clone https://github.com/fentezi/olx-scraper
```

2. Install dependencies:

```bash
go get ./...
```

3. Set the Telegram bot token as an environment variable:

```bash
export TOKEN=your-bot-token
```

4. Build and run the bot:

```bash
go build ./olx-scraper
```