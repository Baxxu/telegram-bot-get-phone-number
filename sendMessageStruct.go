package main

type KeyboardButton struct {
	Text           string `json:"text"`
	RequestContact bool   `json:"request_contact,omitempty"`
}

type ReplyKeyboardMarkup struct {
	Keyboard       [][]KeyboardButton `json:"keyboard"`
	ResizeKeyboard bool               `json:"resize_keyboard,omitempty"`
	// OneTimeKeyboard Клавиатура скрывается, но не удаляется
	OneTimeKeyboard bool `json:"one_time_keyboard,omitempty"`
}

type ReplyKeyboardRemove struct {
	RemoveKeyboard bool `json:"remove_keyboard"`
}

type SendMessageBody struct {
	ChatId                   any    `json:"chat_id"`
	Text                     string `json:"text"`
	ReplyToMessageId         int    `json:"reply_to_message_id,omitempty"`
	AllowSendingWithoutReply bool   `json:"allow_sending_without_reply,omitempty"`
	ReplyMarkup              any    `json:"reply_markup,omitempty"`
}
