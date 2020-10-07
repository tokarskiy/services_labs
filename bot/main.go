package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"math/rand"

	"github.com/Syfaro/telegram-bot-api"
)

const tgBotToken = "<Telegram bot token>"
const booksServiceUrl = "http://books:80/api/books"

func main() {
	var bot *tgbotapi.BotAPI
	var err error

	if bot, err = tgbotapi.NewBotAPI(tgBotToken); err != nil {
		log.Panic(err)
	}

	ucfg := tgbotapi.NewUpdate(0)
	ucfg.Timeout = 60

	var updates tgbotapi.UpdatesChannel
	if updates, err = bot.GetUpdatesChan(ucfg); err != nil {
		log.Panic(err)
	}

	for {
		update := <-updates // ждем обновление чата
		if update.Message == nil {
			continue
		}

		if checkQuestion(update.Message.Text) {
			var book Book
			book, err = getRandomBook()
			if err != nil {
				continue
			}

			text := fmt.Sprintf("%s - %s", book.Author, book.Name)
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
			bot.Send(msg)
		}
	}

}

///
///	Check whether user asks about new book
///
func checkQuestion(message string) bool {
	return strings.ToLower(message) == "/newbook"
}

///
///	Function gets a random book from service
///
func getRandomBook() (Book, error) {
	resp, err := http.Get(booksServiceUrl)
	if err != nil {
		return Book{}, err
	}
	defer resp.Body.Close()

	bodyByte, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return Book{}, err
	}

	var data []Book
	if err = json.Unmarshal(bodyByte, &data); err != nil {
		return Book{}, err
	}

	dataLength := len(data)
	dataRandIdx := rand.Intn(dataLength)

	result := data[dataRandIdx]
	return result, nil
}

// Book - struct representing book
type Book struct {
	ID     int64  `json:"id"`
	Name   string `json:"bookName"`
	Author string `json:"authorName"`
}
