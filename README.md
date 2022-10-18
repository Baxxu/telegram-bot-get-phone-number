### telegram-bot-get-phone-number

# Telegram бот для получения номера телефона пользователя
(с согласия пользователя)

## Как использовать

Для работы бота нужен PostgreSql

- создать бота https://core.telegram.org/bots#6-botfather
- создать в коде глобальную константу с api key бота
- создать в коде глобальную константу для подключения к базе PostgreSql
```
const (
ApiKey string = "тут api key бота"

DataBaseUrl string = "postgres://username:password@localhost:5432/database_name"
)
```
- скомпилировать
- бот готов

### Адрес бота
https://t.me/GetPhoneNumber_123_Bot
