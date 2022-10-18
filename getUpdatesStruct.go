package main

type GetUpdatesBody struct {
	Offset         int      `json:"offset,omitempty"`
	Limit          int      `json:"limit,omitempty"`
	Timeout        int      `json:"timeout,omitempty"`
	AllowedUpdates []string `json:"allowed_updates,omitempty"`
}

type Update struct {
	Ok     bool     `json:"ok"`
	Result []Result `json:"result"`
}

type Result struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message,omitempty"`
}

type Message struct {
	MessageId int     `json:"message_id"`
	From      From    `json:"from,omitempty"`
	Chat      Chat    `json:"chat"`
	Date      int     `json:"date"`
	Text      string  `json:"text,omitempty"`
	Contact   Contact `json:"contact,omitempty"`
}

type From struct {
	Id           int    `json:"id"`
	IsBot        bool   `json:"is_bot"`
	FirstName    string `json:"first_name"`
	Username     string `json:"username"`
	LanguageCode string `json:"language_code"`
}

type Chat struct {
	Id        int    `json:"id"`
	FirstName string `json:"first_name"`
	Username  string `json:"username"`
	Type      string `json:"type"`
}

type Contact struct {
	PhoneNumber string `json:"phone_number"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name,omitempty"`
	UserId      int    `json:"user_id,omitempty"`
	Vcard       string `json:"vcard,omitempty"`
}
