package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// ApiUrl https://api.telegram.org/bot<token>/METHOD_NAME
	ApiUrl = `https://api.telegram.org/bot%s/%s`
)

var (
	Client = &http.Client{
		Timeout: time.Second * 90,
	}

	dataBase = DataBase{}
)

type Bot struct {
	ApiUrl         string
	Offset         int
	GetUpdatesBody GetUpdatesBody
}

//TODO Настроить log. для сохранения логов в базу

func main() {
	dataBase.Connect()
	defer dataBase.Close()

	bot := Bot{
		ApiUrl: fmt.Sprintf(ApiUrl, ApiKey, "getUpdates"),
		Offset: 1,
		GetUpdatesBody: GetUpdatesBody{
			Offset:         1,
			Timeout:        60,
			AllowedUpdates: []string{"message"},
		},
	}

	for {
		bot.GetUpdates()
	}
}

// GetUpdates Парсер по своей сути однопоточный, потому что для получения обновлений нужно распарсить ответ и получить updateId.
//
// Вебхуки лучше, что для них нужен домен и внешний IP
func (bot *Bot) GetUpdates() {
	//Новый оффсет
	bot.GetUpdatesBody.Offset = bot.Offset + 1

	//Джейсоним
	getUpdatesBodyTemp, err := json.Marshal(bot.GetUpdatesBody)
	if err != nil {
		log.Printf("Error JSON Marshal func GetUpdates(offset int)\n%s\n", err)
		time.Sleep(time.Second * 30)
		return
	}

	//Получаю данные с сервера
	resp, err := Client.Post(bot.ApiUrl, "application/json", bytes.NewReader(getUpdatesBodyTemp))
	if err != nil {
		log.Printf("Error Client POST func GetUpdates(offset int)\n%s\n", err)
		time.Sleep(time.Second * 10)
		return
	}
	defer resp.Body.Close()

	//Читаю данные
	tempData, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("Error ReadAll resp.Body\n%s\n", err)
		time.Sleep(time.Second * 60)
		return
	}
	resp.Body.Close()

	//fmt.Printf("%s\n", tempData)

	//Парсим данные из JSON
	var data Update
	err = json.Unmarshal(tempData, &data)
	if err != nil {
		log.Printf("Error JSON parsing\n%s\n", err)
		time.Sleep(time.Second * 30)
		return
	}

	//Чекаем есть ли ответ
	if !data.Ok {
		log.Printf("GetUpdates Update not Ok\n%+v\n", data)
		time.Sleep(time.Second * 10)
		return
	}
	if len(data.Result) == 0 {
		time.Sleep(time.Second)
		return
	}

	//Парсим ответ и ищу наибольший UpdateId
	for _, result := range data.Result {
		go Parse(result)

		if result.UpdateId > bot.Offset {
			bot.Offset = result.UpdateId
		}
	}
}

// Parse
//
// Ищу /start (старт бота) и телефонный номер, и вызываю хэдлеры к ним
func Parse(result Result) {
	//fmt.Println("Parse called")

	//Чекаю а не бот ли это. У ботов нет телефонов.
	//Мб не делать эту проверку
	if result.Message.From.IsBot {
		return
	}

	//Это должен быть приватный чат.
	//Мб не делать эту проверку
	if !strings.EqualFold(result.Message.Chat.Type, "private") {
		return
	}

	//Челик прислал /start, значит надо запросить номер
	if strings.EqualFold(result.Message.Text, `/start`) {
		HandleStartMessage(result.Message)
		return
	}

	//Челик прислал номер, значит надо закинуть номер в базу
	if result.Message.Contact.PhoneNumber != "" {
		HandlePhoneNumber(result.Message)
		return
	}
}

func HandleStartMessage(message Message) {
	sendButton(message)
}

func HandlePhoneNumber(message Message) {
	//Чекаю что контакт принадлежит отправителю
	if message.Contact.UserId == message.From.Id {
		//fmt.Printf("%+v\n", message.Contact)
		//fmt.Println(message.Contact.PhoneNumber)
		//Добавляю в базу
		dataBase.Add(message)

		sendSuccessMessage(message)
	} else {
		//Контакт не принадлежит отправителю
		//fmt.Printf("%+v\n", message.Contact)
		//Если есть ID то надо добавить в базу
		if message.Contact.UserId > 0 {
			//fmt.Printf("%+v\n", message.Contact)
			//fmt.Println(message.Contact.PhoneNumber)
			dataBase.Add(message)
		}
		sendFailureMessage(message)
	}
}

func SendMessage(sendMessageBody SendMessageBody) {
	tempUrl := fmt.Sprintf(ApiUrl, ApiKey, "sendMessage")

	sendMessageBodyTemp, err := json.Marshal(sendMessageBody)
	if err != nil {
		log.Printf("Error JSON Marshal\n%s\n", err)
		return
	}

	_, err = Client.Post(tempUrl, "application/json", bytes.NewReader(sendMessageBodyTemp))
	if err != nil {
		log.Printf("Error Client POST\n%s\n", err)
		return
	}
}

func sendButton(message Message) {
	button := KeyboardButton{
		Text:           "Отправить номер телефона",
		RequestContact: true,
	}

	keyboard := ReplyKeyboardMarkup{
		Keyboard: [][]KeyboardButton{
			{button},
		},
		ResizeKeyboard: true,
	}

	sendMessageBody := SendMessageBody{
		ChatId:      message.Chat.Id,
		Text:        "Ля я ботиха!\nПрив чд кд,\nСкинь цифры",
		ReplyMarkup: keyboard,
	}

	SendMessage(sendMessageBody)
}

func sendSuccessMessage(message Message) {
	sendMessageBody := SendMessageBody{
		ChatId:                   message.Chat.Id,
		Text:                     "Номер получен",
		ReplyToMessageId:         message.MessageId,
		AllowSendingWithoutReply: true,
		ReplyMarkup: ReplyKeyboardRemove{
			RemoveKeyboard: true,
		},
	}

	SendMessage(sendMessageBody)
}

func sendFailureMessage(message Message) {
	sendMessageBody := SendMessageBody{
		ChatId:                   message.Chat.Id,
		Text:                     "ID контакта не соответствует вашему ID пользователя",
		ReplyToMessageId:         message.MessageId,
		AllowSendingWithoutReply: true,
	}

	SendMessage(sendMessageBody)
}
